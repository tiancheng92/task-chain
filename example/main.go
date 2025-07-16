package main

import (
	"errors"
	"github.com/tiancheng92/task-chain"
	"github.com/tiancheng92/task-chain/internal/log"
	"strconv"
	"strings"
	"time"
)

func init() {
	task_chain.Init(
		task_chain.SetDBByDsn("xxxxxxxxxx:xxxxxxxxxx@tcp(xxxxxxxxxx:3306)/xxxxxxxxxx?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true"),
		task_chain.SetLarkAppInfo("xxxxxxxxxx", "xxxxxxxxxx", "https://xxxxxxxxxx.xxxxxxxxxx.com/xxxxxxxxxx/xxxxxxxxxx?task_id=%d"),
		task_chain.SetLLMInfo("https://xxxxxxxxxx.xxxxxxxxxx.com", "xxxxxxxxxx", "xxxxxxxxxx/xxxxxxxxxx"),
		task_chain.SetCustomPlugins(new(PrintNodeStatusPlugin)),
	)

}

type PrintNodeStatusPlugin struct {
	*task_chain.DefaultCustomPlugin
}

func (p *PrintNodeStatusPlugin) WhenNodeStatusChange(chainID uint64, nodeID uint64, status string) error {
	log.Infof("chain id: %d node id: %d status: %s", chainID, nodeID, status)
	return nil
}

type TaskOne struct {
	*task_chain.Task
}

func NewTaskOne(i int) *TaskOne {
	return &TaskOne{
		task_chain.NewTask(
			task_chain.SetTaskName("task-"+strconv.Itoa(i)),
			task_chain.SetTaskParameter(map[string]any{
				"index":     i,
				"index_str": strconv.Itoa(i),
			}),
			task_chain.SetTaskNameForMsg("测试任务 - "+strconv.Itoa(i)),
			task_chain.SetRetryTimes(3),
			task_chain.IsIgnoreFailed(i == 14 || i == 45),
			task_chain.IsMustExecute(i == 35),
		),
	}
}

func (to *TaskOne) Run() error {
	time.Sleep(2 * time.Second)
	indexStr, err := task_chain.GetParameter[string](to, "index_str")
	if err != nil {
		return err
	}

	if indexStr == "2" {
		to.AddParameterToNextTask("parent_index", 2)
	}

	if strings.Contains(indexStr, "4") {
		return errors.New("dry run error")
	}
	return nil
}

func main() {
	tc := task_chain.NewTaskChain(
		task_chain.SetTaskChainInitiator("xxxxxxxxxx"),
		task_chain.SetTaskChainInitiatorForMsg("6666666666"),
		task_chain.SetTaskChainName("task_chain_1"),
		task_chain.SetTaskChainNameForMsg("任务链-1"),
		task_chain.SetTaskChainInfoForMsg(map[string]string{
			"测试发信字段1": "1",
			"测试发信字段2": "2",
			"测试发信字段3": "3",
			"测试发信字段4": "4",
			"测试发信字段5": "5",
			"测试发信字段6": "6",
		}),
		task_chain.SetTaskChainSendLarkUserID(false, "xxxxxxxxxx"),
	)

	ln := tc.AddTask(NewTaskOne(1)).
		AddTask(NewTaskOne(2)).
		AddTask(NewTaskOne(3))
	for i := range 4 {
		ln.AddTask(NewTaskOne(4 + i*10)).
			AddTask(NewTaskOne(5 + i*10)).
			AddTask(NewTaskOne(6 + i*10))
	}
	tc.AddTask(NewTaskOne(91)).
		AddTask(NewTaskOne(92)).
		AddTask(NewTaskOne(93))

	err := tc.Run()
	if err != nil {
		log.Errorf("%+v", err)
	}
}
