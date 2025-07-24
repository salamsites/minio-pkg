package music

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
)

func Save(ctx context.Context, client *minio.Client, localDir, path, bucketName string) error {
	remoteDir := path + "/hls"
	m3u8FileName, err := findFileWithExtension(localDir, ".m3u8")
	if err != nil {
		return err
	}

	// Если M3U8 файл найден, загрузить его
	if m3u8FileName != "" {
		m3u8Path := filepath.Join(localDir, m3u8FileName)
		if err := uploadFile(ctx, client, bucketName, remoteDir, m3u8Path, "application/vnd.apple.mpegurl"); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("M3U8 файл не найден в директории %s", localDir)
	}

	// Затем сохраняем все TS файлы
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("filepath.Walk---->", err)
			return err
		}

		if filepath.Ext(path) == ".ts" {
			if err := uploadFile(ctx, client, bucketName, remoteDir, path, "video/mp2t"); err != nil {
				fmt.Println("filepath.Ext(path) == .ts", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func uploadFile(ctx context.Context, client *minio.Client, bucketName, remoteDir, filePath, contentType string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	objectName := filepath.Join(remoteDir, filepath.Base(filePath))

	// Получаем размер файла
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, bucketName, objectName, file, stat.Size(), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	fmt.Printf("Файл %s успешно загружен в %s/%s\n", filePath, bucketName, objectName)
	return nil
}

func findFileWithExtension(dir, ext string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ext {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("файл с расширением %s не найден в директории %s", ext, dir)
}
