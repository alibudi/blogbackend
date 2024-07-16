package main

import (
	"log"
	"os"

	"github.com/alibudi/blogbackend/database"
	"github.com/alibudi/blogbackend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	database.Connect()
	err:=godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	port:=os.Getenv("PORT")
	app:=fiber.New()
	routes.Setup(app)
	app.Listen(":"+port)
}