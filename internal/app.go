package internal

import (
	"fmt"
	"restapp/internal/environment"
	"restapp/internal/logic"
	"restapp/internal/logic/database"
	"restapp/internal/logic/logic_http"
	"restapp/internal/logic/logic_websocket"
	"restapp/websocket"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// Initialize gofiber application, including DB and view engine.
func NewApp() (*fiber.App, error) {
	environment.WaitForBuild()

	db, errDBLoad := database.InitDB()
	if errDBLoad != nil {
		log.Error(errDBLoad)
		return nil, errDBLoad
	}

	engine := NewAppHtmlEngine(db)
	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())

	appLogic := &logic.Logic{DB: db}
	UseHTTP := func(handler func(r logic_http.LogicHTTP) error) fiber.Handler {
		return func(c fiber.Ctx) error {
			logic := logic_http.LogicHTTP{
				Logic: appLogic,
				Ctx:   c,
			}
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return handler(logic)
		}
	}

	UseHTTPPage := func(
		templatePath string,
		bind *fiber.Map,
		redirectOn logic_http.RedirectCompute,
		layouts ...string,
	) fiber.Handler {
		bindx := fiber.Map{
			"Title": "?",
		}
		if bind != nil {
			for k, v := range *bind {
				bindx[k] = v
			}
		}
		return UseHTTP(func(r logic_http.LogicHTTP) error {
			return r.RenderPage(
				templatePath,
				&bindx,
				redirectOn,
				layouts...,
			)
		})
	}

	UseWebsocket := func(
		subscribe []string,
		handler func(ws *logic_websocket.LogicWebsocket) error,
	) fiber.Handler {
		return func(c fiber.Ctx) error {
			http := logic_http.LogicHTTP{
				Logic: appLogic,
				Ctx:   c,
			}
			if !websocket.IsWebSocketUpgrade(http.Ctx) {
				http.Ctx.Next()
			}

			ip := http.Ctx.IP()
			bind := http.MapPage(nil)

			websocket.New(func(conn *websocket.Conn) {
				// NOTE: Inside this method the 'c' variable is already updated
				// and some methods can not work. They should be used outside.
				ws := logic_websocket.New(
					appLogic,
					conn,
					app,
					ip,
					&bind,
				)

				user := ws.User()

				logic_websocket.WebsocketConnections.Connect(user.Id, &ws)
				defer logic_websocket.WebsocketConnections.Close(user.Id, &ws)
				for !ws.Closed {
					messageType, message, err := ws.Ctx.ReadMessage()
					if err != nil {
						break
					}

					// start := time.Now()
					ws.MessageType = messageType
					ws.Message = message
					err = handler(&ws)

					// colorErr := fiber.DefaultColors.Green
					// if err != nil {
					// 	colorErr = fiber.DefaultColors.Red
					// }

					// fmt.Printf(
					// 	"%s | %s%3s%s | %13s | %15s | %d | %s%s%s \n",
					// 	time.Now().Format("15:04:05"),
					// 	colorErr,
					// 	"ws",
					// 	fiber.DefaultColors.Reset,
					// 	time.Since(start),
					// 	ws.IP,
					// 	ws.MessageType,
					// 	fiber.DefaultColors.Yellow,
					// 	ws.Message,
					// 	fiber.DefaultColors.Reset,
					// )
					if err != nil {
						break
					}
				}
				ws.Ctx.Close()
			})(http.Ctx)

			return nil
		}
	}

	// static
	app.Get("/static/*", static.New("./web/static", static.Config{Browse: true}))
	app.Get("/partials*", static.New("./web/templates/partials", static.Config{Browse: true}))

	// pages
	app.Get("/", UseHTTPPage("index", &fiber.Map{"Title": "Discover"}, func(r logic_http.LogicHTTP, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/settings", UseHTTPPage("settings", &fiber.Map{"Title": "Settings"},
		func(r logic_http.LogicHTTP, bind *fiber.Map) string {
			if user := r.User(); user == nil {
				return "/"
			}
			return ""
		}, "partials/main"),
	)
	app.Get("/chat", UseHTTPPage("chat", &fiber.Map{"Title": "Home", "IsChatPage": true},
		func(r logic_http.LogicHTTP, bind *fiber.Map) string {
			return ""
		}),
	)
	app.Get("/chat/groups/:group_id", UseHTTPPage("chat", &fiber.Map{"Title": "Group", "IsChatPage": true},
		func(r logic_http.LogicHTTP, bind *fiber.Map) string {
			member, _, group := r.Member()

			if member == nil {
				return "/chat"
			}

			(*bind)["Title"] = group.Nick
			return ""
		}),
	)
	app.Get("/chat/groups/join/:group_name", UseHTTPPage("chat", &fiber.Map{"Title": "Join group", "IsChatPage": true},
		func(r logic_http.LogicHTTP, bind *fiber.Map) string {
			member, group, _ := r.Member()
			if group == nil {
				return "/chat"
			}

			if member != nil {
				return logic.PathRedirectGroup(group.Id)
			}

			(*bind)["Title"] = "Join " + group.Nick
			return ""
		}),
	)

	// get
	app.Get("/groups/:group_id/messages/page/:messages_page", UseHTTP(func(r logic_http.LogicHTTP) error { return r.MessagesPage() }))
	app.Get("/groups/:group_id/members/page/:members_page", UseHTTP(func(r logic_http.LogicHTTP) error { return r.MembersPage() }))

	// post
	app.Post("/account/create", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserSignUp() }))
	app.Post("/account/login", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserLogin() }))
	app.Post("/groups/create", UseHTTP(func(r logic_http.LogicHTTP) error { return r.GroupCreate() }))
	app.Post("/groups/:group_id/messages/create", UseHTTP(func(r logic_http.LogicHTTP) error { return r.MessageCreate() }))

	// put
	app.Put("/account/change/name", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserChangeName() }))
	app.Put("/account/change/email", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserChangeEmail() }))
	app.Put("/account/change/phone", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserChangePhone() }))
	app.Put("/account/change/password", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserChangePassword() }))
	app.Put("/account/logout", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserLogout() }))
	app.Put("/groups/:group_id/join", UseHTTP(func(r logic_http.LogicHTTP) error { return r.GroupJoin() }))
	app.Put("/groups/:group_id/change", UseHTTP(func(r logic_http.LogicHTTP) error { return r.GroupChange() }))

	// delete
	app.Delete("/groups/:group_id/leave", UseHTTP(func(r logic_http.LogicHTTP) error { return r.GroupLeave() }))
	app.Delete("/groups/:group_id", UseHTTP(func(r logic_http.LogicHTTP) error { return r.GroupDelete() }))
	app.Delete("/account/delete", UseHTTP(func(r logic_http.LogicHTTP) error { return r.UserDelete() }))

	// websoket
	app.Get("/groups/:group_id", UseWebsocket([]string{
		logic_websocket.SubForMessages,
	}, func(ws *logic_websocket.LogicWebsocket) error { return ws.SubscribeGroup() }))

	app.Use(UseHTTPPage("partials/x", &fiber.Map{
		"Title":         fmt.Sprintf("%d", fiber.StatusNotFound),
		"StatusCode":    fiber.StatusNotFound,
		"StatusMessage": fiber.ErrNotFound.Message,
		"CenterContent": true,
	}, func(r logic_http.LogicHTTP, bind *fiber.Map) string { return "" }, "partials/main"))

	return app, nil
}
