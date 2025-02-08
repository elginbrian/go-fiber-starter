package main

import (
	"fiber-starter/config"
	"fiber-starter/internal/di"
	"fiber-starter/internal/routes"

	_ "fiber-starter/docs"
)

// @title FIBER STARTER API
// @version 1.0
// @description This is a RESTful API for a simple social media application. It allows users to manage their posts, including creating, updating, and deleting posts, and provides authentication using JWT. The API is built using the Fiber framework and interacts with a PostgreSQL database.
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:8084
// @BasePath /api/v1/
func main() {
	serverPort := config.GetServerPort()
	jwtSecret := config.GetJWTSecret()
	refreshSecret := config.GetRefreshSecret()

	db := config.InitDatabase()
	defer db.Close()

	container := di.NewContainer(db, jwtSecret, refreshSecret)
	app := config.SetupFiber()
	routes.SetupRoutes(app, *container, jwtSecret)

	config.StartServer(app, serverPort)
}