package easygin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

const (
	// BindJSON : if you want to specify specific json bind for your easy gin handler
	BindJSON = "json_bind"
	// BindQuery : if you want to specify specific query bind for your easy gin handler
	BindQuery = "query_bind"
	// BindURI : if you want to specify specific uri bind for your easy gin handler
	BindURI = "uri_bind"

	bindContext = "context_bind"
	ctxSep      = ":::"
)

var (
	// MethodToStatusCode : default status codes returned for each HTTP method
	MethodToStatusCode = map[string]int{
		http.MethodGet:     http.StatusOK,
		http.MethodPost:    http.StatusCreated,
		http.MethodPut:     http.StatusOK,
		http.MethodPatch:   http.StatusOK,
		http.MethodDelete:  http.StatusOK,
		http.MethodHead:    http.StatusNoContent,
		http.MethodOptions: http.StatusOK,
	}
)

// BindContext : if you want to bind from context , make sure it's a struct that has the same fields in the one you're copying from
func BindContext(ctxKey string) string {
	return bindContext + ctxSep + ctxKey
}

func extractContextBinds(bindFroms []string) (ctxKeys []string) {
	for _, v := range bindFroms {
		if strings.Contains(v, bindContext) {
			ctxKeys = append(ctxKeys, strings.Split(v, ctxSep)[1])
		}
	}
	return ctxKeys
}

func withoutContexts(bindFroms []string) (nonCtx []string) {
	for _, v := range bindFroms {
		if !strings.Contains(v, bindContext) {
			nonCtx = append(nonCtx, v)
		}
	}
	return nonCtx
}

// To : Converts an easygin handler to a gin handler
func To[T IRequest](inputFunc func(reqDTO T) *Response, bindFrom ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			dtoCastFrom T
			err         error

			contextKeysToBind = extractContextBinds(bindFrom)
			canBindAll        = len(withoutContexts(bindFrom)) == 0
		)

		if ctx.Request.Body != http.NoBody && canBindAll || SliceContains(bindFrom, BindJSON) {
			err = ctx.BindJSON(&dtoCastFrom)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, requestErrorWraper(err))
				return
			}
		}

		if canBindAll || SliceContains(bindFrom, BindQuery) {
			err = ctx.BindQuery(&dtoCastFrom)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, requestErrorWraper(err))
				return
			}
		}

		if canBindAll || SliceContains(bindFrom, BindURI) {
			err = ctx.BindUri(&dtoCastFrom)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, requestErrorWraper(err))
				return
			}
		}

		for _, cKey := range contextKeysToBind {
			val, exists := ctx.Get(cKey)
			if !exists {
				fmt.Println("warning!! context key == " + cKey + " does not exist")
				continue
			}
			err := copier.Copy(&dtoCastFrom, val)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, fmt.Errorf("context key %s made a bad cast ! %w", cKey, err).Error())
				return
			}
		}

		err = dtoCastFrom.Validate()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, requestErrorWraper(err))
			return
		}
		res := inputFunc(dtoCastFrom)
		if res == nil {
			ctx.JSON(http.StatusInternalServerError, requestErrorWraper(errors.New("No response in bdoy")))
			return
		}
		responseCode := res.HTTPStatusCode
		if responseCode == 0 {
			responseCode = MethodToStatusCode[ctx.Request.Method]
			res.HTTPStatusCode = responseCode
		}
		ctx.JSON(responseCode, res.JSONDATA())
	}
}
