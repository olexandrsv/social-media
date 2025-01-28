package message

import (
	"mime/multipart"
	"social-media/internal/common/files"
)

type Message struct {
	id          string
	userID      int
	text        string
	imagesPaths []string
	filesPaths  []string
}

func New(id string, userID int, text string, imagesPaths, filesPaths []string) *Message {
	return &Message{
		id:          id,
		userID:      userID,
		text:        text,
		imagesPaths: imagesPaths,
		filesPaths:  filesPaths,
	}
}

func (m *Message) ID() string {
	return m.id
}

func (m *Message) UserID() int {
	return m.userID
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) ImagesPaths() []string {
	return m.imagesPaths
}

func (m *Message) FilesPaths() []string {
	return m.filesPaths
}

func (m *Message) Update(text string, newImages, newFiles []*multipart.FileHeader, deletedImages, deletedFiles []string) error {
	m.text = text
	addedImages, err := files.Process(newImages)
	if err != nil {
		return err
	}
	addedFiles, err := files.Process(newFiles)
	if err != nil {
		return err
	}

	remainedImages, err := files.Remained(m.imagesPaths, deletedImages)
	if err != nil {
		return err
	}
	remainedFiles, err := files.Remained(m.filesPaths, deletedFiles)
	if err != nil {
		return err
	}

	if err := files.Delete(deletedImages); err != nil {
		return err
	}
	if err := files.Delete(deletedFiles); err != nil {
		return err
	}

	m.imagesPaths = append(remainedImages, addedImages...)
	m.filesPaths = append(remainedFiles, addedFiles...)

	return nil
}
