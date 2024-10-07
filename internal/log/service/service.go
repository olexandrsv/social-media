package service

import "fmt"

type Service interface {
	Error(string)
}

type logService struct{}

func NewService() Service {
	return &logService{}
}

func (srv *logService) Error(msg string) {
	fmt.Println(msg)
}
