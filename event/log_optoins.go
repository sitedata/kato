package event

// GetLoggerOption
func GetLoggerOption(status string) map[string]string {
	return map[string]string{"step": "appruntime", "status": status}
}

//GetCallbackLoggerOption
func GetCallbackLoggerOption() map[string]string {
	return map[string]string{"step": "callback", "status": "failure"}
}

//GetTimeoutLoggerOption
func GetTimeoutLoggerOption() map[string]string {
	return map[string]string{"step": "callback", "status": "timeout"}
}

//GetLastLoggerOption
func GetLastLoggerOption() map[string]string {
	return map[string]string{"step": "last", "status": "success"}
}
