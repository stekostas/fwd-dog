package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AboutPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.gohtml", gin.H{})
}

func CreditsPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "credits.gohtml", gin.H{})
}
