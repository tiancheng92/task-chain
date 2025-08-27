package task_chain

import (
	"context"
	"time"

	"github.com/looplab/fsm"
)

type TaskInterface interface {
	Run() error
	AddParameterToNextTask(key string, value any)
	getRetryTimes() int
	getStartTime() *time.Time
	getEndTime() *time.Time
	getStatus() string
	getName() string
	getNameForMsg() string
	getParameters() map[string]any
	getParameter(key string) any
	getFailedReason() string
	getNextTaskParameter() map[string]any
	isIgnoreFailed() bool
	isMustExecute() bool
	addParameter(key string, value any)
	setStatusToAbandon()
	setStatusToRunning()
	setStatusToSuccess()
	setStatusToFailed(err error)
	getChan() <-chan string
	close()
}

func runTask(t TaskInterface) {
	t.setStatusToRunning()
	var err error
	for range t.getRetryTimes() + 1 {
		err = t.Run()
		if err == nil {
			break
		}
	}
	if err != nil {
		t.setStatusToFailed(err)
	} else {
		t.setStatusToSuccess()
	}
}

type Task struct {
	ctx               context.Context
	status            *fsm.FSM
	c                 chan string
	startTime         *time.Time
	endTime           *time.Time
	parameter         map[string]any
	nextTaskParameter map[string]any
	name              string
	nameForMsg        string
	statusStr         string
	failedReason      string
	ignoreFailed      bool
	mustExecute       bool
	retryTimes        int
}

func NewTask(opts ...taskOption) *Task {
	s := getTaskSetting(opts)

	t := &Task{
		ctx:          context.Background(),
		c:            make(chan string),
		name:         s.taskName,
		nameForMsg:   s.taskNameForMsg,
		ignoreFailed: s.ignoreFailed,
		mustExecute:  s.mustExecute,
		parameter:    s.parameter,
		retryTimes:   s.retryTimes,
	}

	t.status = fsm.NewFSM("waiting", fsm.Events{
		{Name: "abandon", Src: []string{"waiting"}, Dst: "abandon"},
		{Name: "running", Src: []string{"waiting"}, Dst: "running"},
		{Name: "success", Src: []string{"running"}, Dst: "success"},
		{Name: "failed", Src: []string{"running"}, Dst: "failed"},
	}, fsm.Callbacks{
		"enter_abandon": func(_ context.Context, _ *fsm.Event) {
			now := time.Now()
			t.startTime = &now
			t.endTime = &now
		},
		"enter_running": func(_ context.Context, _ *fsm.Event) {
			now := time.Now()
			t.startTime = &now
		},
		"enter_success": func(_ context.Context, _ *fsm.Event) {
			now := time.Now()
			t.endTime = &now
		},
		"enter_failed": func(_ context.Context, _ *fsm.Event) {
			now := time.Now()
			t.endTime = &now
		},
		"enter_state": func(_ context.Context, e *fsm.Event) {
			t.statusStr = e.Dst
			t.c <- e.Dst
		},
	})

	t.statusStr = t.status.Current()

	return t
}

func (t *Task) AddParameterToNextTask(key string, value any) {
	if t.nextTaskParameter == nil {
		t.nextTaskParameter = make(map[string]any)
	}

	t.nextTaskParameter[key] = value
}

func (t *Task) addParameter(key string, value any) {
	if t.parameter == nil {
		t.parameter = make(map[string]any)
	}

	t.parameter[key] = value
}

func (t *Task) getChan() <-chan string {
	return t.c
}

func (t *Task) setStatusToAbandon() {
	err := t.status.Event(t.ctx, "abandon")
	if err != nil {
		panic(err)
	}
}

func (t *Task) setStatusToRunning() {
	err := t.status.Event(t.ctx, "running")
	if err != nil {
		panic(err)
	}
}

func (t *Task) setStatusToSuccess() {
	err := t.status.Event(t.ctx, "success")
	if err != nil {
		panic(err)
	}
}

func (t *Task) setStatusToFailed(err error) {
	t.failedReason = err.Error()
	if e := t.status.Event(t.ctx, "failed"); e != nil {
		panic(e)
	}
}

func (t *Task) getStatus() string {
	return t.status.Current()
}

func (t *Task) getFailedReason() string {
	return t.failedReason
}

func (t *Task) getName() string {
	return t.name
}

func (t *Task) getNameForMsg() string {
	return t.nameForMsg
}

func (t *Task) getParameter(key string) any {
	if value, ok := t.parameter[key]; !ok {
		return nil
	} else {
		return value
	}
}

func (t *Task) getStartTime() *time.Time {
	return t.startTime
}

func (t *Task) getEndTime() *time.Time {
	return t.endTime
}

func (t *Task) getParameters() map[string]any {
	return t.parameter
}

func (t *Task) getNextTaskParameter() map[string]any {
	return t.nextTaskParameter
}

func (t *Task) getRetryTimes() int {
	return t.retryTimes
}

func (t *Task) isIgnoreFailed() bool {
	return t.ignoreFailed
}

func (t *Task) isMustExecute() bool {
	return t.mustExecute
}

func (t *Task) close() {
	close(t.c)
}
