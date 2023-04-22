package main

import (
	"errors"

	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	UserID string // must have the same name as your dto field
}

type UserInput struct {
	UserID string `json:"user_id" uri:"user_id"`
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
	return easygin.Res(u.UserID) // will return back 123
}

func AuthMiddleware(ctx *gin.Context) {
	ctx.Set("user_login_info", LoginInfo{
		UserID: "123",
	})
	ctx.Next()
}

func main() {
	en := gin.Default()
	en.GET("_login_info", AuthMiddleware, easygin.To(HandleUsers, easygin.BindContext("user_login_info")))
	en.Run(":80")
}
