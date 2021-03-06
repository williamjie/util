package util

import (
	"time"

	"github.com/satori/go.uuid"
)

const (
	HIGH uint = iota
	NORMAL
	LOW
)

type Task struct {
	ID       string      `json:"id"`
	Priority uint        `json:"priority"`
	Time     time.Time   `json:"time"`
	Retries  uint        `json:"retries"`
	Content  interface{} `json:"content"`
}

func NewTask(id string, content interface{}) Task {
	return Task{
		ID:       id,
		Time:     time.Now().Local(),
		Retries:  0,
		Content:  content,
		Priority: NORMAL,
	}
}

type TaskService struct {
	*RedisService
}

func NewTaskService(redisService *RedisService) *TaskService {
	CheckNil(redisService, "redis service should not be nil")
	return &TaskService{
		redisService,
	}
}

func (self *TaskService) Enq(content interface{}) error {
	CheckNil(content, "content should not be nil")
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	v, err := ToJsonString(NewTask(id.String(), content))
	if err != nil {
		return err
	}
	return self.Rpush(v)
}

func (self *TaskService) Deq() (*Task, error) {
	reply, err := self.Lpop()
	if err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, nil
	}
	var task Task
	err = ToInstance(reply, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
