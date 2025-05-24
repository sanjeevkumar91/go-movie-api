package movies

import (
	"github.com/gin-gonic/gin"
)

type MoviesController struct{}

func NewMoviesController() *MoviesController {
	return &MoviesController{}
}

func (h MoviesController) SendMessage(c *gin.Context) {
	c.String(200, "Hello, World!")
}
