package main

import (
	"fmt"
	"os"
	"thinkmate/database"
	"thinkmate/router"
	"thinkmate/services"
)

func main() {
	database.StartDB()
	services.StartOpenAIClient()
	var PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))

	r := router.SetupRouter()

	r.Run(PORT)
}
