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
	Info(string)
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

func Info(msg string){
	logger.Info(msg)
}

type logClient struct {
	log.LogClient
}

func (c logClient) Error(err error) {
	msg := fmt.Sprintf("%+v\n", err)
	_, err = c.LogClient.Error(context.Background(), &log.LogRequest{Msg: msg})
	if err != nil {
		fmt.Println(err)
	}
}

func (c logClient) Info(msg string) {
	_, err := c.LogClient.Error(context.Background(), &log.LogRequest{Msg: msg})
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

func (l mockLogger) Info(msg string) {
	l.t.Log(msg)
}
