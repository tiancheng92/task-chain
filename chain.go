package task_chain

import (
	"github.com/Yostardev/gf"
	"github.com/pkg/errors"
	"github.com/tiancheng92/task-chain/internal/lark"
	"github.com/tiancheng92/task-chain/internal/log"
	"github.com/tiancheng92/task-chain/internal/modal"
	"gorm.io/gorm"
	"sync"
)

var plugins []string

type chain struct {
	*startNode
	id                  uint64
	username            string
	taskChainName       string
	taskChainNameForMsg string
	infoForMsg          map[string]any
	larkUserID          []string
	larkGroupID         []string
	larkMsgID           []string
}

func NewTaskChain(opts ...chainOption) *chain {
	s := getChainSetting(opts)
	return &chain{
		username:            s.username,
		taskChainName:       s.taskChainName,
		taskChainNameForMsg: s.taskChainNameForMsg,
		infoForMsg:          s.infoForMsg,
		larkUserID:          s.larkUserID,
		larkGroupID:         s.larkGroupID,
	}
}

func (c *chain) AddTask(taskInterface TaskInterface) *node {
	if c.startNode == nil {
		c.startNode = new(startNode)
	}
	return c.startNode.AddTask(taskInterface)
}

func (c *chain) prepareDB() error {
	return modal.GetDB().Transaction(func(tx *gorm.DB) error {
		tc, err := modal.CreateTaskChain(tx, c.username, c.taskChainName, c.taskChainNameForMsg, c.infoForMsg)
		if err != nil {
			return err
		}
		c.id = tc.ID

		return c.startNode.prepareDB(tx, c.id)
	})
}

func (c *chain) check() error {
	if gf.ArrayContains(plugins, "send_lark_msg") && !gf.ArrayContains(plugins, "use_db") {
		return errors.New("if want to send lark msg, please set db first")
	}

	if gf.ArrayContains(plugins, "use_llm") && !gf.ArrayContains(plugins, "use_db") {
		return errors.New("if want to use llm, please set db first")
	}

	if !gf.ArrayContains(plugins, "send_lark_msg") && (len(c.larkGroupID)+len(c.larkUserID) > 0) {
		return errors.New("if want to send lark msg, please set lark app info in init function")
	}

	return nil
}

func (c *chain) Run() error {
	if c.startNode != nil {
		if err := c.check(); err != nil {
			return err
		}

		// 写入数据库，并初始化任务链中的数据（若写入失败，则中断任务链并回滚数据库中数据）
		if gf.ArrayContains(plugins, "use_db") {
			if err := c.prepareDB(); err != nil {
				return errors.WithStack(err)
			}
		}

		// 发送飞书通知（若发送失败，不阻塞任务执行）
		if gf.ArrayContains(plugins, "send_lark_msg") && (len(c.larkGroupID)+len(c.larkUserID) > 0) {
			if err := lark.SendMsg(c.chainID, c.larkGroupID, c.larkUserID); err != nil {
				log.Errorf("%+v", err)
			}
		}

		wg := new(sync.WaitGroup)
		// 开始监听每个任务节点的状态
		c.startNode.startWatch(wg)
		// 开始执行任务
		c.run()
		//等待监听结束
		wg.Wait()
		// 资源回收
		c.close()
		if gf.ArrayContains(plugins, "send_lark_msg") {
			delete(lark.DataMap, c.chainID)
			delete(lark.TimeTickMap, c.chainID)
		}
	}
	return nil
}
