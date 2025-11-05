package controller

import (
    "net/http"
    "flux-panel/golang-backend/internal/app/model"
    "flux-panel/golang-backend/internal/app/response"
    dbpkg "flux-panel/golang-backend/internal/db"
    "github.com/gin-gonic/gin"
)

// POST /api/v1/alerts/recent {limit?}
func AlertsRecent(c *gin.Context) {
    var p struct{ Limit int `json:"limit"` }
    _ = c.ShouldBindJSON(&p)
    if p.Limit <= 0 || p.Limit > 200 { p.Limit = 50 }
    var list []model.Alert
    dbpkg.DB.Order("time_ms desc").Limit(p.Limit).Find(&list)
    c.JSON(http.StatusOK, response.Ok(list))
}

