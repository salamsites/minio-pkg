package sminiochat

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
)

func SaveFile(ctx context.Context, client *minio.Client, mimeType string, file multipart.File, mime, path, bucketName string) error {
	// Перемещаем указатель в начало файла
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	//tempDir := "/home/user/sallamm/sminio/test/temp"
	//tempFile, err := os.CreateTemp(tempDir, "upload-*.tmp")
	//if err != nil {
	//	return fmt.Errorf("создание временного файла: %w", err)
	//}
	//
	//defer os.Remove(tempFile.Name())
	//defer tempFile.Close()
	//
	//// Копируем содержимое файла во временный файл
	//fileSize, err := io.Copy(tempFile, file)
	//if err != nil {
	//	return fmt.Errorf("копирование файла: %w", err)
	//}
	//
	//// Перемещаем указатель в начало оригинального файла
	//_, err = file.Seek(0, io.SeekStart)
	//if err != nil {
	//	return fmt.Errorf("перемещение указателя в начало файла: %w", err)
	//}

	fileSize, err := GetFileSize(file)
	if err != nil {
		return fmt.Errorf("get file size failed: %v", err)
	}

	// Создаем объект в MinIO с оригинальным именем файла
	objectName := fmt.Sprintf("%s/original.%s", path, mime)
	_, err = client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		return err
	}

	return nil
}

func GetFileSize(file multipart.File) (int64, error) {
	// Save the current position
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Seek to the end to get the size
	endPos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Seek back to the original position
	_, err = file.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endPos, nil
}
