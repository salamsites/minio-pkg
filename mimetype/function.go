package mimetype

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"mime/multipart"
	"strings"
)

// GetPrefixExt Prefix, extension
func GetPrefixExt(mimeType string) (string, string) {
	fmt.Println("GetPrefixExt: ", mimeType)
	a := strings.Split(mimeType, "/")
	return a[0], a[1]
}

func Detect(file multipart.File, mimeTypes []string) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	mtype := mimetype.Detect(buffer)

	fmt.Println("mtype: ", mtype)

	for _, mimeType := range mimeTypes {
		//fmt.Println("mime: ", mime)
		if mtype.Is(mimeType) {
			return mimeType, nil
		}
	}

	return "", err
}
