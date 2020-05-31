package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RecoveryHandler handles the recovery process when a fatal error occurs to render the `500.gohtml` template instead of
// simply rendering an empty page.
func RecoveryHandler(c *gin.Context, _ interface{}) {
	c.HTML(http.StatusInternalServerError, "500.gohtml", gin.H{})
}
