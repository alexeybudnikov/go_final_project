package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ResolveHost() string {
	port := 7540
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	host := fmt.Sprintf(":%d", port)
	fmt.Printf("Приложение запущено. Порт приложения: %d", port)
	return host
}
