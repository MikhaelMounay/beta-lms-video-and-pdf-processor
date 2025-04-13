package shared

import (
	"path/filepath"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/utils"
)

type IFile interface {
	GetPath() string
	GetName() string
	GetGrade() string
	GetId() string
	GetDir() string
	SetGrade(grade string)
	SetPath(path string)
	SetId(id string)
	GetNewPath() string
	ProcessFile() error
	GetUploadPaths() []string
}

type File struct {
	Path        string
	Name        string
	Grade       string
	Id          string
	Dir         string
	UploadPaths []string
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetGrade() string {
	return f.Grade
}

func (f *File) GetId() string {
	return f.Id
}

func (f *File) GetDir() string {
	return f.Dir
}

func (f *File) SetGrade(grade string) {
	f.Grade = grade
}

func (f *File) SetPath(path string) {
	f.Path = path
	f.Dir = filepath.Dir(path)
	f.Name = filepath.Base(path)
}

func (f *File) SetId(id string) {
	if id != "" {
		f.Id = id
	} else {
		f.Id = utils.GenerateRandomString(10)
	}
}

func (f *File) GetNewPath() string {
	return f.Dir + "/" + f.Id
}

func (f *File) ProcessFile() error {
	return nil
}

func (f *File) GetUploadPaths() []string {
	return f.UploadPaths
}
