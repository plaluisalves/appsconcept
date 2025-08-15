package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FizzBuzzParams struct {
	Int1  int    `form:"int1" json:"int1" binding:"required,gt=0"`
	Int2  int    `form:"int2" json:"int2" binding:"required,gt=0"`
	Limit int    `form:"limit" json:"limit" binding:"required,gt=0,lte=10000"`
	Str1  string `form:"str1" json:"str1" binding:"required,max=50"`
	Str2  string `form:"str2" json:"str2" binding:"required,max=50"`
}

// HandlerFizzBuzz handles the FizzBuzz endpoint.
func (a *App) HandlerFizzBuzz(c *gin.Context) {

	// log.Printf("[%s] handler hit\n", FizzBuzzEndpoint)

	var params FizzBuzzParams

	// Bind and validate query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameters: " + err.Error(),
		})
		log.Printf("[%s] %v\n", FizzBuzzEndpoint, err)
		return
	}

	// Increase metric
	a.metrics.IncHits(params.Int1, params.Int2, params.Limit, params.Str1, params.Str2)

	// Build result using strings.Builder
	var sb strings.Builder
	results := make([]string, 0, params.Limit)
	for i := 1; i <= params.Limit; i++ {
		sb.Reset()
		if i%params.Int1 == 0 {
			sb.WriteString(params.Str1)
		}
		if i%params.Int2 == 0 {
			sb.WriteString(params.Str2)
		}
		if sb.Len() == 0 {
			sb.WriteString(strconv.Itoa(i))
		}
		results = append(results, sb.String())
	}

	c.JSON(http.StatusOK, results)

	// log.Printf("[%s] handler complete\n", FizzBuzzEndpoint)
}
