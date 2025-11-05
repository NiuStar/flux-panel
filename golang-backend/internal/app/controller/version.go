package controller

import (
    "net/http"
    "os"

    appver "flux-panel/golang-backend/internal/app/version"
    "flux-panel/golang-backend/internal/app/response"
    "github.com/gin-gonic/gin"
)

// GET /api/v1/version
func Version(c *gin.Context) {
    // server.version from main package
    serverVer := appver.Get()
    // agent version (expected agent binary baseline)
    agentVer := os.Getenv("AGENT_VERSION")
    if agentVer == "" { agentVer = "go-agent-1.0.4" }
    c.JSON(http.StatusOK, response.Ok(map[string]string{
        "server":   serverVer,
        "agent":    agentVer,
    }))
}
