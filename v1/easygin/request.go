package easygin

// IRequest : request interface that is expected by the easy gin handler
type IRequest interface {
	// extra validations happen here
	Validate() error
	// format the error response
	ValidationErrorFormat(err error) any
}
