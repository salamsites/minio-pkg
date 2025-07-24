package video

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// originalFilePath, compressedFilePath, some_error
func createTempFile(tempDir string, inputFile multipart.File, fileName string) (string, string, error) {
	fmt.Println("file_name video--->", fileName)
	originalFile, err := os.CreateTemp(tempDir, fmt.Sprintf("*%s", fileName))
	if err != nil {
		fmt.Println(err, "--->os.CreateTemp originalFile")
		return "", "", err
	}
	defer originalFile.Close()

	if _, err := io.Copy(originalFile, inputFile); err != nil {
		fmt.Println(err, "--->os.Copy")
		return "", "", err
	}

	compressedFile, err := os.CreateTemp(tempDir, "output-*.mp4")
	if err != nil {
		fmt.Println(err, "--->os.CreateTemp compressedFile")
		return "", "", err
	}
	defer compressedFile.Close()

	return originalFile.Name(), compressedFile.Name(), nil
}

func removeTempFile(originalFilePath, compressedFilePath string) error {
	if err := os.Remove(originalFilePath); err != nil {
		return err
	}
	if err := os.Remove(compressedFilePath); err != nil {
		return err
	}
	return nil
}
