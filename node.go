package task_chain

import (
	"strings"
	"sync"

	"github.com/Yostardev/gf"
	"github.com/tiancheng92/task-chain/internal/lark"
	"github.com/tiancheng92/task-chain/internal/llm"
	"github.com/tiancheng92/task-chain/internal/log"
	"github.com/tiancheng92/task-chain/internal/modal"
	"gorm.io/gorm"
)

type startNode struct {
	chainID       uint64
	id            uint64
	nextTaskNodes []*node
}

func (n *startNode) AddTask(task TaskInterface) *node {
	ln := &node{
		task: task,
	}

	n.nextTaskNodes = append(n.nextTaskNodes, ln)
	return ln
}

func (n *startNode) prepareDB(tx *gorm.DB, chainID uint64) error {
	m, err := modal.CreateStartTaskNode(tx, chainID)
	if err != nil {
		return err
	}

	n.chainID = chainID
	n.id = m.ID

	for i := range n.nextTaskNodes {
		if err = n.nextTaskNodes[i].prepareDB(tx, chainID, m.ID); err != nil {
			return err
		}
	}
	return nil
}

func (n *startNode) startWatch(wg *sync.WaitGroup) {
	for i := range n.nextTaskNodes {
		n.nextTaskNodes[i].startWatch(wg)
	}
}

func (n *startNode) run() {
	wg := new(sync.WaitGroup)
	for i := range n.nextTaskNodes {
		wg.Go(func() {
			n.nextTaskNodes[i].run(false)
		})
	}
	wg.Wait()
	return
}

func (n *startNode) close() {
	wg := new(sync.WaitGroup)
	for i := range n.nextTaskNodes {
		wg.Go(func() {
			n.nextTaskNodes[i].close()
		})
	}
	wg.Wait()
	return
}

type node struct {
	chainID       uint64
	id            uint64
	task          TaskInterface
	nextTaskNodes []*node
}

func (n *node) AddTask(task TaskInterface) *node {
	ln := &node{
		task: task,
	}
	n.nextTaskNodes = append(n.nextTaskNodes, ln)

	return ln
}

func (n *node) startWatch(wg *sync.WaitGroup) {
	n.watch(wg)
	for i := range n.nextTaskNodes {
		n.nextTaskNodes[i].startWatch(wg)
	}
}

func (n *node) watch(wg *sync.WaitGroup) {
	wg.Go(func() {
		c := n.task.getChan()
		for {
			select {
			case status := <-c:
				if gf.ArrayContains(plugins, "use_db") {
					ent, err := modal.GetTaskNode(modal.GetDB(), n.id)
					if err != nil {
						log.Errorf("%+v", err)
						return
					}

					ent.Status = status
					ent.StartTime = n.task.getStartTime()
					ent.EndTime = n.task.getEndTime()
					ent.FailedReason = n.task.getFailedReason()
					ent.Parameter = n.task.getParameters()

					if _, err = modal.UpdateTaskNode(modal.GetDB(), n.id, ent); err != nil {
						log.Errorf("%+v", err)
						return
					}

					if _, err = modal.UpdateTaskChain(modal.GetDB(), n.chainID); err != nil {
						log.Errorf("%+v", err)
						return
					}

					if gf.ArrayContains(plugins, "use_llm") && ent.FailedReason != "" {
						ent.FailedReasonAfterAIAnalyze, err = llm.AnalyzeError(ent.FailedReason)
						if err != nil {
							log.Errorf("%+v", err)
						} else {
							if _, err = modal.UpdateTaskNode(modal.GetDB(), n.id, ent); err != nil {
								log.Errorf("%+v", err)
							}
						}
					}
				}
				if gf.ArrayContains(plugins, "send_lark_msg") {
					if err := lark.UpdateMsg(n.chainID); err != nil {
						log.Errorf("%+v", err)
						return
					}
				}

				for i := range customPluginList {
					if err := customPluginList[i].WhenNodeStatusChange(n.chainID, n.id, status); err != nil {
						log.Errorf("%+v", err)
						return
					}
				}

				if strings.Contains("failed/success/abandon", status) {
					return
				}
			}
		}
	})
}

func (n *node) prepareDB(tx *gorm.DB, chainID, parentID uint64) error {
	m, err := modal.CreateTaskNode(tx, chainID, parentID, n.task.getName(), n.task.getNameForMsg(), n.task.getParameters(), n.task.isIgnoreFailed(), n.task.isMustExecute())
	if err != nil {
		return err
	}

	n.chainID = chainID
	n.id = m.ID

	for i := range n.nextTaskNodes {
		if err = n.nextTaskNodes[i].prepareDB(tx, chainID, m.ID); err != nil {
			return err
		}
	}
	return nil
}

func (n *node) run(isAbandon bool) {
	if isAbandon {
		n.task.setStatusToAbandon()
	} else {
		runTask(n.task)
	}

	for k, v := range n.task.getNextTaskParameter() {
		n.setNextTaskParameter(k, v)
	}

	wg := new(sync.WaitGroup)
	for i := range n.nextTaskNodes {
		wg.Go(func() {
			n.nextTaskNodes[i].run(!(n.task.getStatus() == "success" || n.task.isIgnoreFailed() || n.nextTaskNodes[i].task.isMustExecute()))
		})
	}
	wg.Wait()
	return
}

func (n *node) setNextTaskParameter(k string, v any) {
	for i := range n.nextTaskNodes {
		n.nextTaskNodes[i].task.addParameter(k, v)
		n.nextTaskNodes[i].setNextTaskParameter(k, v)
	}
}

func (n *node) close() {
	n.task.close()
	wg := new(sync.WaitGroup)
	for i := range n.nextTaskNodes {
		wg.Go(func() {
			n.nextTaskNodes[i].close()
		})
	}
	wg.Wait()
}
