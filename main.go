package main

import (
	"thinkmate/database"
	"thinkmate/router"
	"thinkmate/services"
)

func main() {
	database.StartDB()
	services.StartOpenAIClient()
	var PORT = ":9090"

	r := router.SetupRouter()

	r.Run(PORT)
}
