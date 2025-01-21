package api

import (
	"net/http"

	api "github.com/alexeybudnikov/go_final_project/internal/api/auth"
	"github.com/alexeybudnikov/go_final_project/internal/service"
	"github.com/go-chi/chi"
)

func NewRouter(service service.TaskService) *chi.Mux {
	handler := NewHandler(service)
	r := chi.NewRouter()

	webDir := "./web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handler.ApiNextDate)
	r.Get("/api/task", api.ValidateJWT(handler.ApiTaskGet))
	r.Get("/api/tasks", api.ValidateJWT(handler.ApiTasksGetAll))
	r.Post("/api/task/done", api.ValidateJWT(handler.ApiTaskDone))
	r.Post("/api/task", api.ValidateJWT(handler.ApiTaskCreate))
	r.Put("/api/task", api.ValidateJWT(handler.ApiTaskUpdate))
	r.Delete("/api/task", api.ValidateJWT(handler.ApiDeleteTask))
	r.Post("/api/signin", handler.LoginHanler)

	return r
}
