package main

import (
	"log"

	"github.com/flux-panel/go-backend/handlers"
	"github.com/flux-panel/go-backend/models"
	"github.com/flux-panel/go-backend/services"
	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=True"
	db, err := models.InitDB(dsn)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	nodeService := services.NewNodeService(db)
	nodeHandler := handlers.NewNodeHandler(nodeService)

	tunnelService := services.NewTunnelService(db)
	tunnelHandler := handlers.NewTunnelHandler(tunnelService)

	r := gin.Default()
	userHandler.RegisterRoutes(r)
	nodeHandler.RegisterRoutes(r)
	tunnelHandler.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
