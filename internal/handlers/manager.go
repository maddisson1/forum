package handlers

import (
	"forum/app"
	"forum/internal/service"
)

type handlers struct {
	service service.ServiceI
	app     *app.Application
}

func New(s service.ServiceI, app *app.Application) *handlers {
	return &handlers{
		s,
		app,
	}
}
