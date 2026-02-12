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
	host := config.App.LogService.Host
	port := config.App.LogService.GrpcPort
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
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
	fmt.Printf("%+v\n", err)
	logger.Error(err)
}

func Info(msg string) {
	fmt.Println(msg)
	logger.Info(msg)
}

func Infof(pattern string, args ...any) {
	Info(fmt.Sprintf(pattern, args...))
}

func Logf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	Info(msg)
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
