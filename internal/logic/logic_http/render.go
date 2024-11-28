package logic_http

import (
	"restapp/internal/logic"

	"github.com/gofiber/fiber/v3"
)

// Should return redirect path or empty string.
type RedirectCompute func(r LogicHTTP, bind *fiber.Map) string

// Render a page using a template.
// Special
func (r LogicHTTP) RenderPage(templatePath string, bind *fiber.Map, redirect RedirectCompute, layouts ...string) error {
	bindx := r.MapPage(bind)
	if redirect != nil {
		if path := redirect(r, bind); path != "" {
			return r.Ctx.Redirect().To(path)
		}
	}
	if title, ok := (*bind)["Title"].(string); ok {
		bindx["Title"] = "Restapp - " + title
	}
	return r.Ctx.Render(templatePath, bindx, layouts...)
}

func (r LogicHTTP) MapPage(bind *fiber.Map) fiber.Map {
	bindx := fiber.Map{}
	if user, tokenErr := r.User(); user != nil {
		bindx["User"] = user
	} else if tokenErr != nil {
		bindx["TokenError"] = true
		bindx["Message"] = "Authorization error"
		bindx["Id"] = "local-token-error"
	}

	if group := r.Group(); group != nil {
		bindx["Group"] = group
	}

	bindx = logic.MapMerge(&bindx, bind)
	return bindx
}

// Renders the danger message html element.
func (r LogicHTTP) RenderDanger(message, id string) error {
	return r.Ctx.Render("partials/danger", fiber.Map{
		"Id":      id,
		"Message": message,
	})
}

// Renders the warning message html element.
func (r LogicHTTP) RenderWarning(message, id string) error {
	return r.Ctx.Render("partials/warning", fiber.Map{
		"Id":      id,
		"Message": message,
	})
}

// Renders the success message html element.
func (r LogicHTTP) RenderSuccess(message, id string) error {
	return r.Ctx.Render("partials/success", fiber.Map{
		"Id":      id,
		"Message": message,
	})
}