package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func listLinksHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		rows, _ := db.Query("SELECT code,target,clicks,created_at FROM links WHERE user_id=?", userID)
		defer rows.Close()
		// Build and return JSON array of links…
	}
}
