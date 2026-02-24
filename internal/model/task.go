package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/pb"
	"github.com/google/uuid"
)

type TaskStatus int

const (
	Pending TaskStatus = iota
	Scheduled
	Completed
	Failed
)

type Task struct {
	ID uuid.UUID `db:"id"`

	URL     string     `db:"url"`
	Method  string     `db:"method"`
	Headers JSONHeader `db:"headers"`
	Body    []byte     `db:"body"`

	Status       TaskStatus `db:"status"`
	RunAt        time.Time  `db:"run_at"`
	RetryCount   int        `db:"retry_count"`
	Created      time.Time  `db:"created"`
	Updated      time.Time  `db:"updated"`
	ErrorMessage string     `db:"error_message"`
}

type Tasks []Task

type JSONHeader http.Header

func (task Task) ToProto() *pb.ProtoTask {
	headers := make(map[string]*pb.HeaderValues)
	for k, v := range task.Headers {
		headers[k] = &pb.HeaderValues{Values: v}
	}

	return &pb.ProtoTask{
		Id:      task.ID.String(),
		Url:     task.URL,
		Method:  task.Method,
		Headers: headers,
		Body:    task.Body,
	}
}

func FromProto(pbTask *pb.ProtoTask) (Task, error) {
	if pbTask == nil {
		return Task{}, errors.New("pbTask is nil")
	}

	id, err := uuid.Parse(pbTask.Id)
	if err != nil {
		return Task{}, err
	}

	headers := make(JSONHeader, len(pbTask.Headers))
	for k, v := range pbTask.Headers {
		headers[k] = v.Values
	}

	return Task{
		ID:      id,
		URL:     pbTask.Url,
		Method:  pbTask.Method,
		Headers: headers,
		Body:    pbTask.Body,
	}, nil
}

func (header JSONHeader) Value() (driver.Value, error) {
	return json.Marshal(header)
}

func (header *JSONHeader) Scan(value any) error {
	if b, ok := value.([]byte); !ok {
		return errors.New("type assertion to []byte failed")
	} else {
		return json.Unmarshal(b, &header)
	}
}
