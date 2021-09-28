package lib

type WebhookExecutionError struct {
	PreExecutionError bool
	ErrorObj          error
}

func NewWebhookExecutionError(preExec bool, err error) *WebhookExecutionError {
	return &WebhookExecutionError{
		PreExecutionError: preExec,
		ErrorObj:          err,
	}
}

func (whe WebhookExecutionError) Error() string {
	return whe.ErrorObj.Error()
}
