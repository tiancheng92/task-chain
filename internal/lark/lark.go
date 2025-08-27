package lark

import (
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tiancheng92/task-chain/internal/log"
	"github.com/tiancheng92/task-chain/internal/modal"
)

var (
	appID, appSecret, superLinkUrlFmt string
	TimeTickMap                       = make(map[uint64]<-chan time.Time)
	DataMap                           = make(map[uint64]string)
)

func SetAppInfo(id, secret, urlFmt string) {
	appID = id
	appSecret = secret
	superLinkUrlFmt = urlFmt
}

func getLarkMsgContent(taskChainID uint64) (string, error) {
	tc, err := modal.GetTaskChainByID(modal.GetDB(), taskChainID, true)
	if err != nil {
		return "", errors.WithStack(err)
	}

	msg := newMessageContentForTask().
		setTaskName(tc.NameForMsg).
		setTaskID(tc.ID).
		setStatus(tc.Status).
		setBasicInfo(tc.StartTime, tc.EndTime, tc.UsernameForMsg)

	var keys []string

	for k := range tc.InfoForMsg {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for i := range keys {
		msg.addParameter(keys[i], tc.InfoForMsg[keys[i]].(string))
	}

	if len(tc.Nodes) <= 30 {
		deepMap := make(map[uint64]int)

		nodeCount := len(tc.Nodes)

		maxWidth := 0

		for i := range tc.Nodes {
			if len(tc.Nodes[i].NextNodeIDs) > maxWidth {
				maxWidth = len(tc.Nodes[i].NextNodeIDs)
			}
		}

		for i := range tc.Nodes {
			n := tc.Nodes[i]

			if n.Name == "start" {
				deepMap[n.ID] = -1
				continue
			}

			deep := deepMap[n.ParentID] + 1

			deepMap[n.ID] = deep

			if nodeCount == 2 || maxWidth == 1 {
				deep = 0
			}

			msg.addSubtasks(deep, n.NameForMsg, n.Status, n.StartTime, n.EndTime, n.IgnoreFailed)
		}
	} else {
		msg.tooManySubtasks()
	}

	return msg.string(), nil
}

func SendMsg(chainID uint64, larkGroupIDs, larkUserIDs []string) error {
	var (
		msgIDs []string
		wg     sync.WaitGroup
	)

	content, err := getLarkMsgContent(chainID)
	if err != nil {
		return errors.WithStack(err)
	}

	for i := range larkUserIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msgID, err := sendToUser(larkUserIDs[i], content)
			if err != nil {
				log.Errorf("%+v", err)
				return
			}
			msgIDs = append(msgIDs, msgID)
		}()
	}

	for i := range larkGroupIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msgID, err := sendToGroup(larkGroupIDs[i], content)
			if err != nil {
				log.Errorf("%+v", err)
				return
			}
			msgIDs = append(msgIDs, msgID)
		}()
	}

	wg.Wait()

	_, err = modal.UpdateTaskChainMsgID(modal.GetDB(), chainID, msgIDs)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func UpdateMsg(taskChainID uint64) error {
	if _, ok := TimeTickMap[taskChainID]; !ok {
		TimeTickMap[taskChainID] = time.Tick(200 * time.Millisecond)
	}

	<-TimeTickMap[taskChainID]

	tc, err := modal.GetTaskChainByID(modal.GetDB(), taskChainID, true)
	if err != nil {
		return errors.WithStack(err)
	}

	content, err := getLarkMsgContent(taskChainID)
	if err != nil {
		return errors.WithStack(err)
	}

	if d, ok := DataMap[taskChainID]; ok && d == content {
		return nil
	}

	DataMap[taskChainID] = content
	for i := range tc.MsgIDs {
		err = updateMsg(tc.MsgIDs[i], content)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
