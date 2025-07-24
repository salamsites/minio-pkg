package audio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"io"
	"mime/multipart"
)

func Save(ctx context.Context, client *minio.Client, tempDir string, mimeType string, file multipart.File, fileName string, sizes []util.Size, path, bucketName string) (int64, error) {
	fmt.Println("Upload Save Audio")

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	// create temp
	//tempDir := "/Users/meylis/Documents/salam-messenger/backend/pkg/sminio/test/temp"
	originalFilePath, err := createTempFile(tempDir, file, fileName)
	if err != nil {
		return 0, err
	}

	defer removeTempFile(originalFilePath)

	duration, err := audioVideoDuration(originalFilePath)
	if err != nil {
		return 0, err
	}

	fmt.Println(duration)

	// screen

	audioBuff, err := bufferFromFile(originalFilePath)

	objectName := util.GetOriginalAudioPath(path)

	_, err = client.PutObject(ctx, bucketName, objectName, &audioBuff, int64(audioBuff.Len()), minio.PutObjectOptions{ContentType: mimetype.MP3})
	if err != nil {
		return 0, err
	}

	fmt.Println("Audio uploaded successfully")

	return duration, nil
}
