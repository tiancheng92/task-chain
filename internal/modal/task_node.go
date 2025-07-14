package modal

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type TaskNode struct {
	ID                         uint64                      `json:"id" gorm:"primary_key;type:bigint unsigned;not null"`
	ParentID                   uint64                      `json:"parent_id" gorm:"type:bigint unsigned;not null;comment:父节点ID"`
	NextNodeIDs                datatypes.JSONSlice[uint64] `json:"next_node_ids" gorm:"type:json;comment:子节点ID列表"`
	StartTime                  *time.Time                  `json:"start_time" gorm:"comment:开始时间"`
	EndTime                    *time.Time                  `json:"end_time" gorm:"comment:结束时间"`
	Name                       string                      `json:"name" gorm:"type:varchar(64);not null;comment:任务名称"`
	NameForMsg                 string                      `json:"name_for_msg" gorm:"comment:发信用任务名称"`
	Parameter                  datatypes.JSONMap           `json:"parameter" gorm:"type:json;comment:参数"`
	Status                     string                      `json:"status" gorm:"type:enum('waiting', 'running', 'failed', 'success', 'abandon');not null;default:'waiting';comment:状态"`
	FailedReason               string                      `json:"failed_reason" gorm:"type:longtext;not null;comment:失败原因"`
	FailedReasonAfterAIAnalyze string                      `json:"failed_reason_after_ai_analyze" gorm:"type:longtext;not null;comment:AI分析的错误原因"`
	IgnoreFailed               bool                        `json:"ignore_failed" gorm:"comment:是否忽略报错"`
	ChainID                    uint64                      `json:"chain_id" gorm:"type:bigint unsigned;index;not null;comment:任务链ID"`
}

func GetTaskNode(db *gorm.DB, id uint64) (*TaskNode, error) {
	var ent TaskNode
	err := db.Model(new(TaskNode)).Where("id = ?", id).First(&ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &ent, nil
}

func GetTaskNodesByChainID(db *gorm.DB, chainID uint64) ([]*TaskNode, error) {
	var ent []*TaskNode
	err := db.Model(new(TaskNode)).Where("chain_id = ?", chainID).Find(&ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ent, nil
}

func CreateStartTaskNode(db *gorm.DB, chainID uint64) (*TaskNode, error) {
	ent := &TaskNode{
		ChainID:      chainID,
		ParentID:     0,
		Name:         "start",
		NameForMsg:   "开始",
		Status:       "success",
		IgnoreFailed: false,
	}
	err := db.Model(new(TaskNode)).Create(ent).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ent, nil
}

func UpdateTaskNode(db *gorm.DB, id uint64, taskNode *TaskNode) (*TaskNode, error) {
	err := db.Model(new(TaskNode)).Where("id = ?", id).Updates(taskNode).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return taskNode, nil
}

func AppendTaskNode(db *gorm.DB, id, nextNodeID uint64) (*TaskNode, error) {
	ent, err := GetTaskNode(db, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ent.NextNodeIDs = append(ent.NextNodeIDs, nextNodeID)
	return UpdateTaskNode(db, id, ent)
}

func CreateTaskNode(db *gorm.DB, chainID uint64, parentID uint64, name string, nameForMsg string, parameter map[string]any, ignoreFailed bool) (*TaskNode, error) {
	ent := &TaskNode{
		ChainID:      chainID,
		ParentID:     parentID,
		Name:         name,
		NameForMsg:   nameForMsg,
		Status:       "waiting",
		Parameter:    parameter,
		IgnoreFailed: ignoreFailed,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Model(new(TaskNode)).Create(ent).Error
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = AppendTaskNode(tx, parentID, ent.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ent, nil
}
