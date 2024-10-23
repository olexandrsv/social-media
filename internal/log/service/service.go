package service

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Service interface {
	Error(string)
}

type logService struct {
	log *log.Logger
	mux sync.Mutex
}

func NewService() Service {
	file, err := os.OpenFile("./internal/common/app/log/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	l := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return &logService{
		log: l,
	}
}

func (srv *logService) Error(msg string) {
	srv.mux.Lock()
	fmt.Println(msg)
	defer srv.mux.Unlock()
	srv.log.Println(msg)
}
