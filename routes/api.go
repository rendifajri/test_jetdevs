package routes

import (
	"net/http"

	"test_jetdevs/controllers"

	"github.com/labstack/echo/v4"
)

func Build(e *echo.Echo) {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	echo.NotFoundHandler = func(c echo.Context) error {
		resp["status"] = "not_found"
		resp["message"] = "Route Not Found"
		return c.JSON(http.StatusNotFound, resp)
	}
	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		resp["status"] = "not_allowed"
		resp["message"] = "Route Not Allowed"
		return c.JSON(http.StatusMethodNotAllowed, resp)
	}
	e.GET("/article", controllers.ArticleIndex)
	e.GET("/article/content/:id", controllers.ArticleContent)
	e.POST("/article", controllers.ArticlePost)

	e.GET("/article/comment/:id", controllers.CommentIndex)
	e.POST("/article/comment/:id", controllers.CommentPost)
}
