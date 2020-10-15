package service

import (
	"errors"
	"io"
)

type UploadRequest struct {
	FileName   string        `json:"file_name"`
	FolderName string        `json:"folder_name"`
	File       io.ReadSeeker `json:"file"`
}
type UploadResponse struct {
	FileLink string `json:"file_link"`
}

func (ur *UploadRequest) validate() error {
	if ur.FileName == "" {
		return errors.New("empty filename")
	}
	if ur.FolderName == "" {
		return errors.New("empty foldername")
	}
	if ur.File == nil {
		return errors.New("empty file")
	}
	return nil
}
