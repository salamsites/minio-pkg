package util

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
)

func Validate(request *http.Request, key string, maxMemory int64, fileSize int) Err {
	err := request.ParseMultipartForm(maxMemory)
	if err != nil {
		return Err{StatusCode: http.StatusRequestEntityTooLarge, Message: "files is to large"}
	}

	files := request.MultipartForm.File[key]
	if files == nil {
		return Err{StatusCode: http.StatusNotFound, Message: "file not found"}
	}

	if len(files) > fileSize {
		return Err{StatusCode: http.StatusRequestEntityTooLarge, Message: fmt.Sprintf("files max count %d", fileSize)}
	}

	return Err{}
}

func UnsupportedErr(filename string) Err {
	return Err{StatusCode: http.StatusBadRequest, Message: "unsupported file type: " + filename}
}

func FileToMultipartFile(filepath string) (multipart.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
