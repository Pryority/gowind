package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	templateFolder = "./web/templates"
	resourceFolder = "./web/static"
	localPort      = ":8080"
)

func main() {
	r := gin.Default()

	// Configure HTML rendering
	r.LoadHTMLGlob(templateFolder + "/*.tpl.html")

	// Serve static files
	r.Static("/static", resourceFolder)

	// Define routes
	r.GET("/", loadMain)
	r.GET("/:route_name", loadTemplate)

	log.Printf("Server started at http://localhost%s", localPort)
	r.Run(localPort)
}

func loadMain(c *gin.Context) {
	templateName := "index.tpl.html"

	// Render the template
	c.HTML(http.StatusOK, templateName, nil)
}

func loadTemplate(c *gin.Context) {
	routeName := c.Param("route_name")

	// Render the template
	c.HTML(http.StatusOK, "layout.tpl.html", gin.H{
		"RouteName": routeName,
	})
}
