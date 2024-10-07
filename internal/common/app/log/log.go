package log

import (
	"context"
	"fmt"
	"social-media/api/pb/log"
	"social-media/internal/common/app/config"
	"testing"

	"google.golang.org/grpc"
)

type Logger interface {
	Error(error)
}

var logger Logger

func Init() {
	conn, err := grpc.Dial(":"+config.App.LogService.Port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	logger = logClient{
		log.NewLogClient(conn),
	}
}

func InitMock(t testing.TB) {
	logger = mockLogger{
		t: t,
	}
}

func Error(err error) {
	logger.Error(err)
}

type logClient struct {
	log.LogClient
}

func (c logClient) Error(err error) {
	_, err = c.LogClient.Error(context.Background(), &log.LogRequest{Msg: err.Error()})
	if err != nil {
		fmt.Println(err)
	}
}

type mockLogger struct {
	t testing.TB
}

func (l mockLogger) Error(err error) {
	l.t.Log(err)
}
