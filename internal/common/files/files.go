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

func Remained(resources, deletedResources []string) (remainedResources []string, finalErr error) {
	resourcesSate := make(map[string]resourceLifecycle)

	for i := 0; i < len(resources); i++ {
		resourcesSate[resources[i]] = Created
	}

	for i := 0; i < len(deletedResources); i++ {
		deletedResource := deletedResources[i]
		_, ok := resourcesSate[deletedResource]
		if !ok {
			finalErr = common.ErrInvalidData
			return
		}
		resourcesSate[deletedResource] = Deleted
	}

	for resource, state := range resourcesSate {
		if state == Created {
			remainedResources = append(remainedResources, resource)
		}
	}
	return
}

func Process(folder string, files []*multipart.FileHeader) ([]string, error) {
	filesNames := []string{}
	for _, file := range files {
		filename, err := processFile(folder, file)
		if err != nil {
			return nil, err
		}
		filesNames = append(filesNames, filename)
	}
	return filesNames, nil
}

func Get(folder, fileID string) (*os.File, error){
	path := "./upload/"+folder+"/"+fileID
	file, err := os.Open(path)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return file, nil
}

func Delete(folder string, files []string) error {
	for _, file := range files {
		if err := removeFile(folder, file); err != nil {
			return err
		}
	}
	return nil
}

func processFile(folder string, file *multipart.FileHeader) (string, error) {
	extension := filepath.Ext(file.Filename)
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	filename := id.String() + extension
	path := "./upload/"+ folder+ "/" + filename
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

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

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

func removeFile(folder, filename string) error {
	mx.Lock()
	defer mx.Unlock()

	path := "./upload/"+ folder+ "/" + filename
	if err := os.Remove(path); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}
