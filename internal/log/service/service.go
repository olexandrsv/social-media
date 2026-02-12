package service

import (
	"bufio"
	"io"
	"log"
	"os"
	"social-media/internal/common"
	logClient "social-media/internal/common/app/log"
	"sync"

	"github.com/pkg/errors"
)

type Service interface {
	Error(string)
	Info(string)
	GetLogs() (string, error)
}

type logService struct {
	writer *bufio.Writer
	file   *os.File
	log    *log.Logger
	mux    sync.Mutex
}

func NewService() Service {
	file, err := os.OpenFile("./internal/common/app/log/log.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	l := log.New(writer, "", log.LstdFlags|log.Lshortfile)
	return &logService{
		file:   file,
		log:    l,
		writer: writer,
	}
}

func (srv *logService) Error(msg string) {
	srv.println(msg)
}

func (srv *logService) Info(msg string) {
	srv.println(msg)
}

func (srv *logService) println(msg string) {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	srv.log.Println(msg)
}

func (srv *logService) GetLogs() (string, error) {
	srv.mux.Lock()
	if err := srv.writer.Flush(); err != nil {
		srv.mux.Unlock()
		logClient.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	file, err := os.Open("./internal/common/app/log/log.txt")
	if err != nil {
		srv.mux.Unlock()
		logClient.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	content, err := io.ReadAll(file)
	if err != nil {
		srv.mux.Unlock()
		logClient.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}

	srv.mux.Unlock()
	return string(content), nil
}
