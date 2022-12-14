package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test_jetdevs/database"
	"test_jetdevs/structs"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

func ArticleIndex(c echo.Context) error {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	db := database.ShareConnection
	if db == nil {
		// log.Println(err.Error())
		resp["status"] = "error"
		resp["message"] = "database connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	page := 1
	if c.QueryParam("page") != "" {
		page, _ = strconv.Atoi(c.QueryParam("page"))
	}
	offset := (page - 1) * 20
	query := "SELECT id, nickname, title, content, created_on FROM article LIMIT 20 OFFSET " + strconv.Itoa(offset)
	var data []structs.Article
	err := db.Raw(query).Scan(&data)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp["status"] = "success"
	resp["message"] = "Success"
	resp["data"] = data
	return c.JSON(http.StatusOK, resp)
}

func ArticleContent(c echo.Context) error {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	db := database.ShareConnection
	if db == nil {
		// log.Println(err.Error())
		resp["status"] = "error"
		resp["message"] = "database connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	id := c.Param("id")
	query := "SELECT content FROM article WHERE id=" + id
	var data structs.Article
	err := db.Raw(query).Scan(&data)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp["status"] = "success"
	resp["message"] = "Success"
	resp["data"] = data
	return c.JSON(http.StatusOK, resp)
}

func ArticlePost(c echo.Context) error {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	var requestBody = map[string]interface{}{}
	log.Println("controllers.ArticlePost() received param : ")
	json.NewDecoder(c.Request().Body).Decode(&requestBody)
	bb, _ := json.MarshalIndent(requestBody, "", "  ")
	log.Println(string(bb))
	validator := govalidator.New(govalidator.Options{
		Rules: govalidator.MapData{
			"nickname": []string{"required", "max:20"},
			"title":    []string{"required", "max:20"},
			"content":  []string{"required"},
		},
		Data: &requestBody,
	}).ValidateStruct()
	if len(validator) > 0 {
		log.Println("controllers.ArticlePost() invalid validation : ")
		vb, _ := json.MarshalIndent(validator, "", "  ")
		log.Println(string(vb))
		resp["status"] = "error"
		resp["message"] = "Invalid Validation"
		resp["errors"] = validator
		return c.JSON(http.StatusBadRequest, resp)
	}
	db := database.ShareConnection
	if db == nil {
		// log.Println(err.Error())
		resp["status"] = "error"
		resp["message"] = "database connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	query := "INSERT INTO article (nickname, title, content) VALUES ('" + fmt.Sprintf("%v", requestBody["nickname"]) + "', '" + fmt.Sprintf("%v", requestBody["title"]) + "', '" + fmt.Sprintf("%v", requestBody["content"]) + "')"
	err := db.Exec(query)
	if err.Error != nil {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp["status"] = "success"
	resp["message"] = "Success"
	resp["data"] = nil
	return c.JSON(http.StatusOK, resp)
}
