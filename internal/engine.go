package internal

import (
	"restapp/internal/environment"
	"restapp/internal/logic"
	"restapp/internal/logic/database"
	"restapp/internal/logic/model_database"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

// Initialize the view engine.
func NewAppHtmlEngine(db *database.Database) *html.Engine {
	engine := html.New("./web/templates", ".html")

	if environment.Environment == environment.EnvironmentDevelopment {
		engine.Reload(true)
	}

	engine.AddFuncMap(map[string]interface{}{
		"paginateGroups": paginate[model_database.Group],
		"hideEmail": func(text string) string {
			splits := strings.Split(text, "@")
			if len(splits) != 2 {
				return strings.Repeat("*", len(text))
			}
			// a in a@b.c
			before := splits[0]
			// @b.c in a@b.c
			after := "@" + splits[1]

			if len(before) > 5 {
				before = before[:3] + strings.Repeat("*", len(before[3:]))
			} else {
				before = strings.Repeat("*", len(before))
			}
			return before + after
		},
		"jsonTime": func(t time.Time) string {
			return t.Format("2006-01-02T15:04:05.000Z")
		},
		"hidePhone": func(text string) string {
			if len(text) > 5 {
				return text[:4] + strings.Repeat("*", len(text)-4)
			}
			return strings.Repeat("*", len(text))
		},
		"hide": func(text string) string {
			return strings.Repeat("*", len(text))
		},
		"newMap": func(args ...any) fiber.Map {
			result := fiber.Map{}
			for i := 0; i < len(args)-1; i = i + 2 {
				k := args[i].(string)
				v := args[i+1]
				result[k] = v
			}
			return result
		},
		"newArr": func(args ...any) []any {
			return args
		},

		"groupLink": func(group model_database.Group) string {
			return "localhost:3000" + logic.PathRedirectGroupJoin(group.Name)
		},
		"userRightsOf":    db.MemberRights,
		"userMemberOf":    db.MemberById,
		"userMemberships": db.UserGroupList,
		"groupMembers":    db.MemberList,
	})

	return engine
}

func paginate[T any](slice []T, n int) [][]T {
	result := [][]T{}
	for v := range slices.Chunk(slice, n) {
		result = append(result, v)
	}

	return result
}
