package main

import (
	"fmt"
	"net/http"

	"github.com/alexeybudnikov/go_final_project/internal/api"
	"github.com/alexeybudnikov/go_final_project/internal/database"
	"github.com/alexeybudnikov/go_final_project/internal/repository"
	"github.com/alexeybudnikov/go_final_project/internal/service"
	"github.com/alexeybudnikov/go_final_project/utils"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := database.InitDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

	repo := repository.NewTaskRepository(db)
	service := service.NewTaskService(repo)

	r := api.NewRouter(service)

	err = http.ListenAndServe(utils.ResolveHost(), r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
