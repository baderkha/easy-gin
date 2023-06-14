package easygin

var (
	requestErrorWraper = func(err error) any {
		return &StandardResponse{
			Data:    err.Error(),
			Message: "Error",
		}
	}
	requestErrorWrapperNoOp = func(err error) any {
		return err.Error()
	}
)

func NoOpRequestErrResWrapper() {
	requestErrorWraper = requestErrorWrapperNoOp
}

func SetRequestErrResWrapper(f func(err error) any) {
	requestErrorWraper = f
}

// IRequest : request interface that is expected by the easy gin handler
type IRequest interface {
	// extra validations happen here
	Validate() error
}
