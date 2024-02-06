package common

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type Validator struct {
	err error
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) AppendMsg(msg string) {
	err := errors.New(msg)
	if v.err != nil {
		v.err = errors.Join(v.err, err)
	} else {
		v.err = err
	}
}

func (v *Validator) NotEmpty(name, value string) {
	if len(value) == 0 {
		msg := concatenate("Field ", name, " can't be empty")
		v.AppendMsg(msg)
	}
}

func (v *Validator) NotLess(name, value string, len int) {
	if utf8.RuneCountInString(value) < len {
		msg := concatenate("Field ", name, " can't be less than ", strconv.Itoa(len), " symbols")
		v.AppendMsg(msg)
	}
}

func (v *Validator) Err() error {
	if v.err != nil {
		return NewError(http.StatusBadRequest, v.err.Error())
	}
	return nil
}

func concatenate(list ...string) string {
	var b bytes.Buffer
	for _, s := range list {
		b.WriteString(s)
	}
	return b.String()
}
