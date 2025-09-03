package modal

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tiancheng92/datatypes"
	"gorm.io/gorm"
)

type TaskChain struct {
	ID             uint64                      `json:"id" gorm:"primary_key;type:bigint unsigned;not null"`
	StartTime      *time.Time                  `json:"start_time" gorm:"comment:开始时间"`
	EndTime        *time.Time                  `json:"end_time" gorm:"comment:结束时间"`
	Name           string                      `json:"name" gorm:"comment:任务名"`
	Username       string                      `json:"username" gorm:"comment:执行用户"`
	Status         string                      `json:"status" gorm:"type:enum('running', 'failed', 'success');not null;default:'running';comment:状态"`
	NameForMsg     string                      `json:"name_for_msg" gorm:"comment:发信用任务名称"`
	InfoForMsg     datatypes.JSONMap           `json:"info_for_msg" gorm:"type:json;comment:发信用信息"`
	MsgIDs         datatypes.JSONSlice[string] `json:"msg_ids" gorm:"type:json;comment:飞书信息ID"`
	UsernameForMsg string                      `json:"username_for_msg" gorm:"comment:执行用户(发信用)"`
	Nodes          []*TaskNode                 `json:"nodes" gorm:"-"`
}

func GetTaskChainByID(db *gorm.DB, id uint64, withNode bool) (*TaskChain, error) {
	var ent TaskChain
	err := db.Model(new(TaskChain)).Where("id = ?", id).First(&ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if withNode {
		ent.Nodes, err = GetTaskNodesByChainID(db, ent.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &ent, nil
}

func CreateTaskChain(db *gorm.DB, username, name, nameForMsg string, infoForMsg map[string]any, usernameForMsg string) (*TaskChain, error) {
	ent := &TaskChain{
		Name:           name,
		Username:       username,
		NameForMsg:     nameForMsg,
		InfoForMsg:     infoForMsg,
		UsernameForMsg: usernameForMsg,
	}
	err := db.Model(new(TaskChain)).Create(ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ent, nil
}

func UpdateTaskChainMsgID(db *gorm.DB, id uint64, msgIDs []string) (*TaskChain, error) {
	ent, err := GetTaskChainByID(db, id, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ent.MsgIDs = msgIDs

	err = db.Model(new(TaskChain)).Where("id = ?", id).Updates(&ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ent, nil
}

func UpdateTaskChain(db *gorm.DB, id uint64) (*TaskChain, error) {
	ent, err := GetTaskChainByID(db, id, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var (
		startTime, endTime *time.Time
		status             string
	)

	for i := range ent.Nodes {
		if ent.Nodes[i].StartTime != nil {
			if startTime == nil {
				startTime = ent.Nodes[i].StartTime
			} else {
				if ent.Nodes[i].StartTime.Before(*startTime) {
					startTime = ent.Nodes[i].StartTime
				}
			}

		}
	}

	for i := range ent.Nodes {
		if ent.Nodes[i].Status == "running" {
			status = "running"
			break
		}
		if ent.Nodes[i].Status == "failed" {
			status = "failed"
			continue
		}
	}

	if status == "" {
		status = "success"
	}

	if status != "running" {
		for i := range ent.Nodes {
			if ent.Nodes[i].EndTime != nil {
				if endTime == nil {
					endTime = ent.Nodes[i].EndTime
				} else {
					if ent.Nodes[i].EndTime.After(*endTime) {
						endTime = ent.Nodes[i].EndTime
					}
				}

			}
		}
	}

	ent.Status = status
	ent.StartTime = startTime
	ent.EndTime = endTime

	err = db.Model(new(TaskChain)).Where("id = ?", id).Updates(&ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ent, nil
}
