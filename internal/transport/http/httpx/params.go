package httpx

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParamInt(c *gin.Context, name string) (int, bool) {
	v := c.Param(name)
	i, err := strconv.Atoi(v)
	return i, err == nil
}
