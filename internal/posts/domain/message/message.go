package message

import (
	"mime/multipart"
	"social-media/internal/common/files"
	"social-media/internal/posts/domain/user"
)

type MessageI interface {
	ID() string
	UserID() int
	UserName() string
	UserSurname() string
	Text() string
	ImagesPaths() []string
	FilesPaths() []string
	AddUserFullName(name, surname string)
	Update(text string, newImages, newFiles []*multipart.FileHeader, filesFolder string, deletedImages, deletedFiles []string) error
}

type Message struct {
	id          string
	user        *user.User
	text        string
	imagesPaths []string
	filesPaths  []string
}

func New(id string, userID int, text string, imagesPaths, filesPaths []string) *Message {
	return &Message{
		id:          id,
		user:        user.New(userID, "", ""),
		text:        text,
		imagesPaths: imagesPaths,
		filesPaths:  filesPaths,
	}
}

func (m *Message) ID() string {
	return m.id
}

func (m *Message) UserID() int {
	return m.user.ID()
}

func (m *Message) UserName() string {
	return m.user.Name()
}

func (m *Message) UserSurname() string {
	return m.user.Surname()
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

func (m *Message) SetImagesPaths(paths []string) {
	m.imagesPaths = paths
}

func (m *Message) AddUserFullName(name, surname string) {
	m.user.AddFullName(name, surname)
}

func (m *Message) Update(text string, newImages, newFiles []*multipart.FileHeader, filesFolder string, deletedImages, deletedFiles []string) error {
	m.text = text
	addedImages, err := files.Process(filesFolder, newImages)
	if err != nil {
		return err
	}
	addedFiles, err := files.Process(filesFolder, newFiles)
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

	if err := files.Delete(filesFolder, deletedImages); err != nil {
		return err
	}
	if err := files.Delete(filesFolder, deletedFiles); err != nil {
		return err
	}

	m.imagesPaths = append(remainedImages, addedImages...)
	m.filesPaths = append(remainedFiles, addedFiles...)

	return nil
}
