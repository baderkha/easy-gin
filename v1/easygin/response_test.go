package easygin

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestResponse(t *testing.T) {
	dummy := struct {
		Name string
		Date time.Time
	}{
		Name: "some_name",
		Date: time.Now(),
	}
	tests := []struct {
		Input       *Response
		Expected    Response
		Description string
	}{
		{
			Input: Res("something"),
			Expected: Response{
				Data: "something",
			},
			Description: "primitive values should be the same",
		},
		{
			Input: Res("something").Status(http.StatusInternalServerError),
			Expected: Response{
				Data:           "something",
				HTTPStatusCode: http.StatusInternalServerError,
			},
			Description: "status Code Override should work",
		},
		{
			Input: Res(dummy).Status(http.StatusInternalServerError),
			Expected: Response{
				Data:           dummy,
				HTTPStatusCode: http.StatusInternalServerError,
			},
			Description: "object values should be the same",
		},
		{
			Input: Res(&dummy).Status(http.StatusInternalServerError),
			Expected: Response{
				Data:           &dummy,
				HTTPStatusCode: http.StatusInternalServerError,
			},
			Description: "pointer values should be the same",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected.Data, test.Input.Data)
		assert.Equal(t, test.Expected.HTTPStatusCode, test.Input.HTTPStatusCode)
	}
}
