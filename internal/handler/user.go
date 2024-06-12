package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// -------------------- COMMON FUNTION --------------------
func getIDFromParam(c *gin.Context) uint {
	rawUint64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(rawUint64)
}
