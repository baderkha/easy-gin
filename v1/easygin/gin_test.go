package easygin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var (
	dummyErr = errors.New("expected user id")
)

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// MockRequestInput : mock any request
type MockRequestInput struct {
	Ctx                   *gin.Context
	HTTPMethod            string
	PathParams            []gin.Param
	QueryParams           url.Values
	Body                  any
	KeyValueMiddlewareSet map[string]any
}

func MockInboundRequest(m *MockRequestInput) *gin.Context {
	m.Ctx.Request.Method = m.HTTPMethod
	if m.Body != nil {
		m.Ctx.Request.Header.Set("Content-Type", "application/json")
		jsonbytes, err := json.Marshal(m.Body)
		if err != nil {
			panic(err)
		}
		// the request body must be an io.ReadCloser
		// the bytes buffer though doesn't implement io.Closer,
		// so you wrap it in a no-op closer
		m.Ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	} else {
		m.Ctx.Request.Body = http.NoBody
	}
	m.Ctx.Params = m.PathParams
	m.Ctx.Request.URL.RawQuery = m.QueryParams.Encode()
	if m.KeyValueMiddlewareSet != nil && len(m.KeyValueMiddlewareSet) > 0 {
		for k, v := range m.KeyValueMiddlewareSet {
			m.Ctx.Set(k, v)
		}
	}
	return m.Ctx
}

type dummyJSONOnly struct {
	Name  string `json:"name" form:"name" uri:"name"`
	Other int    `json:"other" form:"other" uri:"other" binding:"ne=10"`
}

func (d dummyJSONOnly) Validate() error {
	if d.Name == "" {
		return dummyErr
	}
	return nil
}

func (d dummyJSONOnly) ValidationErrorFormat(err error) any {
	return err.Error()
}

func TestTo(t *testing.T) {

	{
		handler := func(d dummyJSONOnly) *Response {
			return Res(d.Name).Status(305)
		}
		// test case 1 JSON
		{

			// no mode specify // failure EOF json body, should return back error string
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodPost,
					Body:        false, // invalid value
					QueryParams: make(url.Values),
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.MatchRegex(t, w.Body.String(), "json: cannot unmarshal")
			}
			// no mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodPost,
					Body: map[string]any{
						"name": "baderkha",
					}, // invalid value
					QueryParams: make(url.Values),
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				//assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}

			// json mode specify // failure EOF json body, should return back error string
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodPost,
					Body:        false, // invalid value
					QueryParams: make(url.Values),
				})
				ginFunc := To(handler, BindJSON)
				ginFunc(gCtx)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.MatchRegex(t, w.Body.String(), "json: cannot unmarshal")
			}
			// json mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodPost,
					Body: map[string]any{
						"name": "baderkha",
					}, // invalid value
					QueryParams: make(url.Values),
				})
				ginFunc := To(handler, BindJSON)
				ginFunc(gCtx)
				//assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}
		}
		// test case 2 query param binding
		{

			// no mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodGet,
					QueryParams: map[string][]string{
						"name": {
							"baderkha",
						},
					},
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}

			// no mode specify , failure gin validation , not custom // should error
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodGet,
					QueryParams: map[string][]string{
						"name": {
							"baderkha",
						},
						"other": {
							"10",
						},
					},
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.MatchRegex(t, w.Body.String(), "Error")
			}

			// query mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodGet,
					QueryParams: map[string][]string{
						"name": {
							"baderkha",
						},
					},
				})
				ginFunc := To(handler, BindQuery)
				ginFunc(gCtx)
				assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}
		}
		// test case 3 uri param binding
		{
			// no mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodGet,
					QueryParams: make(url.Values),
					PathParams: []gin.Param{
						{
							Key:   "name",
							Value: "baderkha",
						},
					},
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}

			// no mode specify , failure gin validation , not custom // should error
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:        ctx,
					HTTPMethod: http.MethodGet,
					PathParams: []gin.Param{
						{
							Key:   "name",
							Value: "baderkha",
						},
						{
							Key:   "other",
							Value: "10",
						},
					},
				})
				ginFunc := To(handler)
				ginFunc(gCtx)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.MatchRegex(t, w.Body.String(), "Error")
			}

			// uri mode specify // ok , should output
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodGet,
					QueryParams: make(url.Values),
					PathParams: []gin.Param{
						{
							Key:   "name",
							Value: "baderkha",
						},
					},
				})
				ginFunc := To(handler, BindURI)
				ginFunc(gCtx)
				assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}
		}
		// test case 4 set / get struct binding
		{

			// regular struct , this should work and we should have the value passed to our dto handler
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodGet,
					QueryParams: make(url.Values),
					KeyValueMiddlewareSet: map[string]any{
						"some_mwr": struct {
							Name string
						}{
							Name: "baderkha",
						},
					},
				})
				ginFunc := To(handler, BindContext("some_mwr"))
				ginFunc(gCtx)
				assert.Equal(t, 305, w.Code)
				assert.Equal(t, `"baderkha"`, w.Body.String())
			}

			// regular struct key does not exist, this should work but the value will be empty
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodGet,
					QueryParams: make(url.Values),
					KeyValueMiddlewareSet: map[string]any{
						"some_mwr": struct {
							Name string
						}{
							Name: "baderkha",
						},
					},
				})
				ginFunc := To(handler, BindContext("some_mwr2"))
				ginFunc(gCtx)
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}

			// map should not bind and should throw an error
			{
				w := httptest.NewRecorder()
				ctx := GetTestGinContext(w)
				gCtx := MockInboundRequest(&MockRequestInput{
					Ctx:         ctx,
					HTTPMethod:  http.MethodGet,
					QueryParams: make(url.Values),
					KeyValueMiddlewareSet: map[string]any{
						"some_mwr": nil,
					},
				})
				ginFunc := To(handler, BindContext("some_mwr"))
				ginFunc(gCtx)
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.NotEqual(t, `"baderkha"`, w.Body.String())
			}

		}

		// test case 5 failed manual validation
		{
			w := httptest.NewRecorder()
			ctx := GetTestGinContext(w)
			gCtx := MockInboundRequest(&MockRequestInput{
				Ctx:        ctx,
				HTTPMethod: http.MethodGet,
				QueryParams: map[string][]string{
					"name": {
						"",
					},
				},
			})
			ginFunc := To(handler)
			ginFunc(gCtx)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, fmt.Sprintf(`"%s"`, dummyErr.Error()), w.Body.String())
		}

	}
}
