package main

import (
    "log"
    "os"

    app "flux-panel/golang-backend/internal/app"
    appver "flux-panel/golang-backend/internal/app/version"
    "flux-panel/golang-backend/internal/app/util"
    dbpkg "flux-panel/golang-backend/internal/db"
    "flux-panel/golang-backend/internal/app/scheduler"

	"github.com/gin-gonic/gin"
)

func main() {
    // load .env if present
    util.LoadEnv()
    if err := dbpkg.Init(); err != nil {
        log.Fatalf("db init error: %v", err)
    }
    // start schedulers
    scheduler.Start()

	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	app.RegisterRoutes(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "6365"
    }
    log.Printf("flux-panel server version %s", appver.Get())
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
