package easygin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// MethodToStatusCode : default status codes returned for each HTTP method
	MethodToStatusCode = map[string]int{
		http.MethodGet:    http.StatusOK,
		http.MethodPost:   http.StatusCreated,
		http.MethodPut:    http.StatusOK,
		http.MethodPatch:  http.StatusOK,
		http.MethodDelete: http.StatusOK,
	}
)

// To : Converts an easygin handler to a gin handler
func To[T IRequest](inputFunc func(reqDTO T) *Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dtoCastFrom T
		err := ctx.Bind(&dtoCastFrom)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtoCastFrom.ValidationErrorFormat(err))
			return
		}
		err = ctx.BindQuery(&dtoCastFrom)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtoCastFrom.ValidationErrorFormat(err))
			return
		}
		err = ctx.BindUri(&dtoCastFrom)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtoCastFrom.ValidationErrorFormat(err))
			return
		}
		err = dtoCastFrom.Validate()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtoCastFrom.ValidationErrorFormat(err))
			return
		}
		res := inputFunc(dtoCastFrom)
		if res == nil {
			ctx.JSON(http.StatusInternalServerError, dtoCastFrom.ValidationErrorFormat(errors.New("No response in bdoy")))
			return
		}
		responseCode := res.hTTPStatusCode
		if responseCode == 0 {
			responseCode = MethodToStatusCode[ctx.Request.Method]
		}
		ctx.JSON(responseCode, res.data)
	}
}
