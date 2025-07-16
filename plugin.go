package task_chain

var customPluginList []CustomPlugin

type CustomPlugin interface {
	BeforeRunning(chainID uint64) error
	WhenNodeStatusChange(chainID, nodeID uint64, status string) error
	AfterRunning(chainID uint64) error
}

type DefaultCustomPlugin struct{}

func (d *DefaultCustomPlugin) BeforeRunning(_ uint64) error {
	return nil
}

func (d *DefaultCustomPlugin) WhenNodeStatusChange(_, _ uint64, _ string) error {
	return nil
}

func (d *DefaultCustomPlugin) AfterRunning(_ uint64) error {
	return nil
}
