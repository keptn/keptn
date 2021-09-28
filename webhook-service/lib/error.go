package lib

type WebhookExecutionError struct {
	PreExecutionError bool
	ErrorObj          error
	ExecutedRequests  int
}

type WebhookExecutionErrorOpt func(executionError *WebhookExecutionError)

func WithNrOfExecutedRequests(nrRequests int) WebhookExecutionErrorOpt {
	return func(executionError *WebhookExecutionError) {
		executionError.ExecutedRequests = nrRequests
	}
}

func NewWebhookExecutionError(preExec bool, err error, opts ...WebhookExecutionErrorOpt) *WebhookExecutionError {
	whe := &WebhookExecutionError{
		PreExecutionError: preExec,
		ErrorObj:          err,
		ExecutedRequests:  0,
	}

	for _, opt := range opts {
		opt(whe)
	}

	return whe
}

func (whe WebhookExecutionError) Error() string {
	return whe.ErrorObj.Error()
}
