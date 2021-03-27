package internal

type ErrorHandler struct {
	logger *Logger
}

func NewErrorHandler(logger *Logger) *ErrorHandler {
	return &ErrorHandler{logger}
}

func (eh *ErrorHandler) IsError(err error) {
	if err != nil {
		eh.logger.Error(err)
	}
}

func (eh *ErrorHandler) IsFatalError(err error) {
	if err != nil {
		eh.logger.Error(err)
	}
}
