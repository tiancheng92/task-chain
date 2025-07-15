package task_chain

type taskSetting struct {
	taskName       string
	taskNameForMsg string
	parameter      map[string]any
	ignoreFailed   bool
	mustExecute    bool
	retryTimes     int
}

func defaultTaskSetting() *taskSetting {
	return &taskSetting{
		parameter: make(map[string]any),
	}
}

type taskOption func(*taskSetting)

func SetTaskName(taskName string) taskOption {
	return func(setting *taskSetting) {
		setting.taskName = taskName
	}
}
func SetTaskNameForMsg(taskNameForMsg string) taskOption {
	return func(setting *taskSetting) {
		setting.taskNameForMsg = taskNameForMsg
	}
}

func SetTaskParameter(parameter map[string]any) taskOption {
	return func(setting *taskSetting) {
		setting.parameter = parameter
	}
}

func IsIgnoreFailed(t bool) taskOption {
	return func(setting *taskSetting) {
		setting.ignoreFailed = t
	}
}

func IsMustExecute(t bool) taskOption {
	return func(setting *taskSetting) {
		setting.mustExecute = t
	}
}

func SetRetryTimes(retryTimes int) taskOption {
	return func(setting *taskSetting) {
		setting.retryTimes = retryTimes
	}
}

func getTaskSetting(opts []taskOption) *taskSetting {
	s := defaultTaskSetting()

	for i := range opts {
		opts[i](s)
	}
	return s
}
