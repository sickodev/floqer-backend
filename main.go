package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sickodev/floqer-backend/store"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main(){
  if err := godotenv.Load(); err != nil{
    log.Fatalf("Error loading env file")
  }

  app := fiber.New(fiber.Config{
    Prefork: true,
    ServerHeader: "Floqer-Backend",
  })

  app.Use(logger.New())
  app.Use(cors.New(cors.Config{
    AllowMethods: "GET, POST",
    ExposeHeaders: "Content-Type",
  }))

  api := app.Group("/api")

  // VERSION 1 SETTINGS

  v1 := api.Group("/v1")

  v1.Get("/store",store.GetData)

  log.Fatal(app.Listen(os.Getenv("PORT")))
}
