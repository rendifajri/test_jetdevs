package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test_jetdevs/database"
	"test_jetdevs/helpers"
	"test_jetdevs/structs"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func UserIndex(c echo.Context) error {
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
	query := "SELECT id, name FROM user"
	var data []structs.User
	err := db.Raw(query).Scan(&data)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	for i := range data {
		query = "SELECT point FROM user_point WHERE user_id=" + data[i].Id
		var point structs.UserPoint
		err = db.Raw(query).Scan(&point)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println(err.Error)
			resp["status"] = "error"
			resp["message"] = err.Error
			return c.JSON(http.StatusInternalServerError, resp)
		}
		data[i].Point = point

		query = "SELECT point, point_final FROM user_point_history WHERE user_id=" + data[i].Id
		var pointHistory []structs.UserPointHistory
		err = db.Raw(query).Scan(&pointHistory)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println(err.Error)
			resp["status"] = "error"
			resp["message"] = err.Error
			return c.JSON(http.StatusInternalServerError, resp)
		}
		data[i].PointHistory = pointHistory
	}
	resp["status"] = "success"
	resp["message"] = "Success"
	resp["data"] = data
	return c.JSON(http.StatusOK, resp)
}

func UserDetail(c echo.Context) error {
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
	query := "SELECT id, name FROM user WHERE id=" + id
	var user structs.User
	err := db.Raw(query).Scan(&user)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}

	query = "SELECT point FROM user_point WHERE user_id=" + id
	var point structs.UserPoint
	err = db.Raw(query).Scan(&point)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	user.Point = point

	query = "SELECT point, point_final FROM user_point_history WHERE user_id=" + id
	var pointHistory []structs.UserPointHistory
	err = db.Raw(query).Scan(&pointHistory)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println(err.Error)
		resp["status"] = "error"
		resp["message"] = err.Error
		return c.JSON(http.StatusInternalServerError, resp)
	}
	user.PointHistory = pointHistory

	resp["status"] = "success"
	resp["message"] = "Success"
	resp["data"] = user
	return c.JSON(http.StatusOK, resp)
}

func UserPoint(c echo.Context) error {
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
	point, err := strconv.Atoi(c.FormValue("point"))
	if err != nil || point <= 0 {
		resp["status"] = "error"
		resp["message"] = "Point is invalid"
		return c.JSON(http.StatusBadRequest, resp)
	}
	pointFinal := point
	rc, err := helpers.RedisConnect()
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	ticket := time.Now().Format("2006-01-02-150405")
	resp["status"] = "processing"
	resp["message"] = "Transaction is on progress"
	var respData = map[string]interface{}{
		"ticket": ticket,
	}
	resp["data"] = respData
	mrb, _ := json.MarshalIndent(resp, "", "  ")
	log.Println("UserPoint() redis ticket : " + ticket)
	log.Println("UserPoint() redis mrb : " + string(mrb))
	messageDuration := 5 * time.Minute
	err = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis store problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}

	go func(ticket string) {
		rc, _ := helpers.RedisConnect()
		messageDuration := 5 * time.Minute

		query := "SELECT point FROM user_point WHERE user_id=" + id
		var pointDb structs.UserPoint
		errs := db.Raw(query).Scan(&pointDb)
		if errs.Error == gorm.ErrRecordNotFound {
			_, err = db.DB().Exec("INSERT INTO user_point (user_id, point) VALUES (?, ?)", id, point)
			if err != nil {
				log.Println(err.Error())
				resp["status"] = "error"
				resp["message"] = err.Error
				mrb, _ = json.MarshalIndent(resp, "", "  ")
				log.Println("UserPoint() redis mrb : " + string(mrb))
				_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
			}
		} else {
			pointFinal = pointDb.Point + point
			_, err = db.DB().Exec("UPDATE user_point SET point=? WHERE user_id=?", pointFinal, id)
			if err != nil {
				log.Println(err.Error())
				resp["status"] = "error"
				resp["message"] = err.Error
				mrb, _ = json.MarshalIndent(resp, "", "  ")
				log.Println("UserPoint() redis mrb : " + string(mrb))
				_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
			}
		}

		_, err = db.DB().Exec("INSERT INTO user_point_history (user_id, point, point_final) VALUES (?, ?, ?)", id, point, pointFinal)
		if err != nil {
			log.Println(err.Error())
			resp["status"] = "error"
			resp["message"] = err.Error
			mrb, _ = json.MarshalIndent(resp, "", "  ")
			log.Println("UserPoint() redis mrb : " + string(mrb))
			_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
		}
		resp["status"] = "success"
		resp["message"] = "Balance updated"
		mrb, _ = json.MarshalIndent(resp, "", "  ")
		log.Println("UserPoint() redis mrb : " + string(mrb))
		_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
	}(ticket)

	return c.JSON(http.StatusOK, resp)
}

func UserPointMinus(c echo.Context) error {
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
	point, err := strconv.Atoi(c.FormValue("point"))
	if err != nil || point <= 0 {
		resp["status"] = "error"
		resp["message"] = "Point is invalid"
		return c.JSON(http.StatusBadRequest, resp)
	}
	pointFinal := point
	query := "SELECT point FROM user_point WHERE user_id=" + id
	var pointDb structs.UserPoint
	errs := db.Raw(query).Scan(&pointDb)
	if errs.Error == gorm.ErrRecordNotFound {
		resp["status"] = "error"
		resp["message"] = "Not enough point"
		return c.JSON(http.StatusInternalServerError, resp)
	} else {
		if pointDb.Point < pointFinal {
			resp["status"] = "error"
			resp["message"] = "Not enough point"
			return c.JSON(http.StatusInternalServerError, resp)
		}
	}
	rc, err := helpers.RedisConnect()
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	ticket := time.Now().Format("2006-01-02-150405")
	resp["status"] = "processing"
	resp["message"] = "Transaction is on progress"
	var respData = map[string]interface{}{
		"ticket": ticket,
	}
	resp["data"] = respData
	mrb, _ := json.MarshalIndent(resp, "", "  ")
	log.Println("UserPoint() redis ticket : " + ticket)
	log.Println("UserPoint() redis mrb : " + string(mrb))
	messageDuration := 5 * time.Minute
	err = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis store problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}

	go func(ticket string) {
		rc, _ := helpers.RedisConnect()
		messageDuration := 5 * time.Minute

		query := "SELECT point FROM user_point WHERE user_id=" + id
		var pointDb structs.UserPoint
		errs := db.Raw(query).Scan(&pointDb)
		if errs.Error == gorm.ErrRecordNotFound {
			resp["status"] = "error"
			resp["message"] = "Not enough point"
			mrb, _ = json.MarshalIndent(resp, "", "  ")
			log.Println("UserPoint() redis mrb : " + string(mrb))
			_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
			return
		} else {
			if pointDb.Point < pointFinal {
				resp["status"] = "error"
				resp["message"] = "Not enough point"
				mrb, _ = json.MarshalIndent(resp, "", "  ")
				log.Println("UserPoint() redis mrb : " + string(mrb))
				_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
				return
			}
			pointFinal = pointDb.Point - point
			_, err = db.DB().Exec("UPDATE user_point SET point=? WHERE user_id=?", pointFinal, id)
			if err != nil {
				log.Println(err.Error())
				resp["status"] = "error"
				resp["message"] = err.Error
				mrb, _ = json.MarshalIndent(resp, "", "  ")
				log.Println("UserPoint() redis mrb : " + string(mrb))
				_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
			}
		}

		_, err = db.DB().Exec("INSERT INTO user_point_history (user_id, point, point_final) VALUES (?, ?, ?)", id, point, pointFinal)
		if err != nil {
			log.Println(err.Error())
			resp["status"] = "error"
			resp["message"] = err.Error
			mrb, _ = json.MarshalIndent(resp, "", "  ")
			log.Println("UserPoint() redis mrb : " + string(mrb))
			_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
		}
		resp["status"] = "success"
		resp["message"] = "Balance updated"
		mrb, _ = json.MarshalIndent(resp, "", "  ")
		log.Println("UserPoint() redis mrb : " + string(mrb))
		_ = rc.RedisSet(context.Background(), ticket, string(mrb), messageDuration)
	}(ticket)

	return c.JSON(http.StatusOK, resp)
}
func UserProcessStatus(c echo.Context) error {
	var resp = map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}
	ticket := c.FormValue("ticket")
	rc, err := helpers.RedisConnect()
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis connection problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	val, err := rc.RedisGet(context.Background(), ticket)
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "redis get problem"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if val == "nil" {
		resp["status"] = "error"
		resp["message"] = "Data is nil"
		return c.JSON(http.StatusBadRequest, resp)
	} else if val == "" {
		resp["status"] = "error"
		resp["message"] = "Key is not exists"
		return c.JSON(http.StatusBadRequest, resp)
	}
	var respData = map[string]interface{}{}
	err = json.Unmarshal([]byte(val), &respData)
	if err != nil {
		resp["status"] = "error"
		resp["message"] = "Failed to unmarshal"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp = respData
	return c.JSON(http.StatusOK, resp)
}
