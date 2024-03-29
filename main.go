package main

import (
	"os"
	"thinkmate/database"
	"thinkmate/router"
	"thinkmate/services"
)

func main() {
	database.StartDB()
	services.StartOpenAIClient()
	var PORT = os.Getenv("PORT")

	r := router.SetupRouter()

	r.Run(PORT)
}
