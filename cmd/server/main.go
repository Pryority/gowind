// cmd/server/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/pryority/gowind/cmd/server/handlers"
	"github.com/pryority/gowind/cmd/server/initializers"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables:", err.Error())
	}
	initializers.ConnectDB(&config)
	handlers.SetDB(initializers.DB)

}

func main() {
	engine := html.New("./cmd/web/views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	SetupRoutes(app)

	log.Printf("🌩 Server started at http://localhost%s", LocalPort)

	// Start the server
	if err := app.Listen(LocalPort); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
