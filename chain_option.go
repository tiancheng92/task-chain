package task_chain

type chainSetting struct {
	taskChainName       string
	taskChainNameForMsg string
	username            string
	infoForMsg          map[string]any
	larkUserID          []string
	larkGroupID         []string
}

func defaultChainSetting() *chainSetting {
	return &chainSetting{
		infoForMsg: make(map[string]any),
	}
}

type chainOption func(*chainSetting)

func SetTaskChainName(taskName string) chainOption {
	return func(setting *chainSetting) {
		setting.taskChainName = taskName
	}
}

func SetTaskChainInitiator(username string) chainOption {
	return func(setting *chainSetting) {
		setting.username = username
	}
}

func SetTaskChainNameForMsg(name string) chainOption {
	return func(setting *chainSetting) {
		setting.taskChainNameForMsg = name
	}
}

func SetTaskChainInfoForMsg(info map[string]string) chainOption {
	return func(setting *chainSetting) {
		for k, v := range info {
			setting.infoForMsg[k] = v
		}
	}
}

func SetTaskChainSendLarkUserID(userIDs ...string) chainOption {
	return func(setting *chainSetting) {
		setting.larkUserID = userIDs
	}
}

func SetTaskChainSendLarkGroupID(groupIDs ...string) chainOption {
	return func(setting *chainSetting) {
		setting.larkGroupID = groupIDs
	}
}

func getChainSetting(opts []chainOption) *chainSetting {
	s := defaultChainSetting()

	for i := range opts {
		opts[i](s)
	}
	return s
}
