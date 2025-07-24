package audio

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// originalFilePath, compressedFilePath, some_error
func createTempFile(tempDir string, inputFile multipart.File, fileName string) (string, error) {
	fmt.Println("file_name audio--->", fileName)
	originalFile, err := os.CreateTemp(tempDir, fmt.Sprintf("*%s", fileName))
	if err != nil {
		fmt.Println(err, "--->os.CreateTemp originalFile")
		return "", err
	}
	defer originalFile.Close()

	if _, err := io.Copy(originalFile, inputFile); err != nil {
		fmt.Println(err, "--->os.Copy")
		return "", err
	}

	return originalFile.Name(), nil
}

func removeTempFile(originalFilePath string) error {
	if err := os.Remove(originalFilePath); err != nil {
		return err
	}
	return nil
}
