package common

import (
	"os"
	"social-media/internal/common/app/log"

	"github.com/pkg/errors"
)

func ReadSecret(name string) (string, error) {
	content, err := os.ReadFile("/run/secrets/" + name)
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", ErrInternal
	}

	return string(content), err
}
