package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"test_jetdevs/database"
	"test_jetdevs/structs"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

func CommentIndex(c echo.Context) error {
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
	query := "SELECT id, nickname, content, created_on FROM comment WHERE article_id=" + id
	var data []structs.Comment
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

func CommentPost(c echo.Context) error {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	var requestBody = map[string]interface{}{}
	log.Println("controllers.CommentPost() received param : ")
	json.NewDecoder(c.Request().Body).Decode(&requestBody)
	bb, _ := json.MarshalIndent(requestBody, "", "  ")
	log.Println(string(bb))
	validator := govalidator.New(govalidator.Options{
		Rules: govalidator.MapData{
			"nickname": []string{"required", "max:20"},
			"content":  []string{"required"},
		},
		Data: &requestBody,
	}).ValidateStruct()
	if len(validator) > 0 {
		log.Println("controllers.CommentPost() invalid validation : ")
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
	id := c.Param("id")
	comment_id := "null"
	if requestBody["comment_id"] != nil && fmt.Sprintf("%v", requestBody["comment_id"]) != "" {
		comment_id = fmt.Sprintf("%v", requestBody["comment_id"])
	}
	query := "INSERT INTO comment (article_id, comment_id, nickname, content) VALUES (" + id + ", " + comment_id + ", '" + fmt.Sprintf("%v", requestBody["nickname"]) + "', '" + fmt.Sprintf("%v", requestBody["content"]) + "')"
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
