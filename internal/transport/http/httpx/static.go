package httpx

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func SetStaticCache(c *gin.Context, seconds int) {
	c.Header("Cache-Control", "public, max-age="+strconv.Itoa(seconds))
}
