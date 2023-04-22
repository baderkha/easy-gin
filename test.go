package main

import (
	"errors"
	"net/http"

	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/gin-gonic/gin"
)

type UserInput struct {
	UserID            string `json:"user_id" uri:"user_id"`
	SomethingInTheWay string `json:"something_in_the_way" form:"q"`
}

func (u UserInput) Validate() error {
	if u.UserID == "" {
		return errors.New("user id not set")
	}
	return nil
}

func (u UserInput) ValidationErrorFormat(err error) any {
	return map[string]any{
		"err": err.Error(),
	}
}

func HandleUsers(u UserInput) *easygin.Response {
	return easygin.Res(u).Status(http.StatusNotFound)
}

func main() {
	en := gin.Default()
	en.GET("/:user_id", easygin.To(HandleUsers, easygin.BindJSON, easygin.BindQuery))
	en.Run(":80")
}
