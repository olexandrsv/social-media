package files

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var mx sync.Mutex

type resourceLifecycle int

const (
	Created resourceLifecycle = iota
	Deleted
)

func Remained(resources, deletedResources []string) ([]string, error) {
	resourcesSate := make(map[string]resourceLifecycle)

	for i := 0; i < len(resources); i++ {
		resourcesSate[resources[i]] = Created
	}

	for i := 0; i < len(deletedResources); i++ {
		deletedResource := deletedResources[i]
		_, ok := resourcesSate[deletedResource]
		if !ok {
			log.Error(errors.Errorf("attempt to update file from another post: %+v", deletedResource))
			return nil, common.ErrInvalidData
		}
		resourcesSate[deletedResource] = Deleted
	}

	remainedResources := make([]string, 0)
	for resource, state := range resourcesSate {
		if state == Created {
			remainedResources = append(remainedResources, resource)
		}
	}
	return remainedResources, nil
}

func Process(files []*multipart.FileHeader) ([]string, error) {
	filesNames := []string{}
	for _, file := range files {
		filename, err := processFile(file)
		if err != nil {
			return nil, err
		}
		filesNames = append(filesNames, filename)
	}
	return filesNames, nil
}

func Delete(files []string) error {
	for _, file := range files {
		if err := removeFile(file); err != nil {
			return err
		}
	}
	return nil
}

func processFile(file *multipart.FileHeader) (string, error) {
	extension := filepath.Ext(file.Filename)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	filename := id.String() + extension
	path := "./upload/" + filename
	if err := save(file, path); err != nil {
		return "", err
	}
	return filename, nil
}

func save(fileHeader *multipart.FileHeader, path string) error {
	mx.Lock()
	defer mx.Unlock()
	
	file, err := fileHeader.Open()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}
	defer file.Close()
	newFile, err := os.Create(path)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func removeFile(filename string) error {
	mx.Lock()
	defer mx.Unlock()

	path := "./upload/" + filename
	if err := os.Remove(path); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}
