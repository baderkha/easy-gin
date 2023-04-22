package easygin

// Response : response type for easy gin , it's a wrapper around the data with status code override
type Response struct {
	data           interface{}
	hTTPStatusCode int
}

// Status : override wrapper with a status code instead of the default mapped ones
//
// Refer to : MethodToStatusCode to see the default ones
func (r *Response) Status(code int) *Response {
	r.hTTPStatusCode = code
	return r
}

// Res : response
func Res(data any) *Response {
	return &Response{
		data: data,
	}
}
