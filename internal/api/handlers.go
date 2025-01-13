package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	api "github.com/alexeybudnikov/go_final_project/internal/api/auth"
	"github.com/alexeybudnikov/go_final_project/internal/models"
	"github.com/alexeybudnikov/go_final_project/internal/service"
	"github.com/alexeybudnikov/go_final_project/utils"
)

type Handler struct {
	service service.TaskService
}

func NewHandler(service service.TaskService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ApiNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return
	}
	res, err := utils.NextDate(nowTime, date, repeat)

	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *Handler) ApiTaskCreate(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	taskId, err := h.service.CreateTask(task)
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]string{"id": fmt.Sprintf("%d", taskId)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApiTasksGetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	taskResponses := []models.TaskResponse{}
	for _, task := range tasks {
		taskResponses = append(taskResponses, models.TaskResponse{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	response := models.TasksResponse{Tasks: taskResponses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApiTaskGet(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	idParsed, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": "не указан идентификатор"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	task, err := h.service.GetTaskByID(int64(idParsed))
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	taskResponse := models.TaskResponse{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	response := taskResponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApiTaskUpdate(w http.ResponseWriter, r *http.Request) {
	var task models.TaskResponse

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(task.ID)
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)

		return
	}
	t := models.Task{
		ID:      int64(taskID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	err = h.service.UpdateTaskByID(t)
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApiTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	err = h.service.DoneTask(int64(parsedId))
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApiDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	err = h.service.DeleteTask(int64(parsedId))
	if err != nil {
		response := map[string]string{"error": fmt.Sprintf("%s", err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LoginHanler(w http.ResponseWriter, r *http.Request) {
	secretPassword := os.Getenv("TODO_PASSWORD")

	var passwd models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&passwd)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if passwd.Password != secretPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := api.GenerateJWT(secretPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) StartPage(w http.ResponseWriter, r *http.Request) {
	tokenString, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login.html", http.StatusOK)
		return
	}

	if tokenString.Value == "" {
		http.Redirect(w, r, "/login.html", http.StatusOK)
		return
	}

	http.Redirect(w, r, "/index.html", http.StatusOK)
}
