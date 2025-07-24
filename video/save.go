package video

import (
	"context"
	"fmt"
	"github.com/mowshon/moviego"
	"github.com/salamsites/minio-pkg/image"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"io"
	"mime/multipart"
	"os"
)

func Save(ctx context.Context, client *minio.Client, tempDir string, mimeType string, file multipart.File, fileName string, sizes []util.Size, path, bucketName string) (int64, error) {
	fmt.Println("UploadFeed Save Video")

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	// create temp
	//tempDir := "/Users/meylis/Documents/salam-messenger/backend/pkg/sminio/test/temp"
	originalFilePath, compressedFilePath, err := createTempFile(tempDir, file, fileName)
	if err != nil {
		return 0, err
	}

	defer removeTempFile(originalFilePath, compressedFilePath)

	// compress
	err = compressToMP4(originalFilePath, compressedFilePath, 23)
	if err != nil {
		return 0, err
	}

	duration, err := audioVideoDuration(originalFilePath)
	if err != nil {
		return 0, err
	}

	fmt.Println("len(sizes): ", len(sizes))
	// screen
	if len(sizes) > 0 {
		tempImagePath := fmt.Sprintf("%s.%s", compressedFilePath, "jpeg")
		fmt.Println(tempImagePath)
		fmt.Println(tempImagePath)
		first, err := moviego.Load(compressedFilePath)
		if err != nil {
			fmt.Println("Error loading first: ", err)
			return 0, err
		}

		screenPath, err := first.Screenshot(0.01, tempImagePath)
		if err != nil {
			fmt.Println("Error converting to path: ", err)
			return 0, err
		}

		f, err := util.FileToMultipartFile(screenPath)
		if err != nil {
			fmt.Println("Error converting to multipart file: ", err)
		}

		errSave := image.Save(ctx, client, mimetype.JPEG, f, sizes, path, bucketName)
		if errSave != nil {
			fmt.Println("Error saving file: ", errSave)
			return 0, errSave
		}

		if err := os.Remove(screenPath); err != nil {
			//return err
		}

		//sizesPath, err := resizeImageToWebpJpeg(outputFilePath+".jpg", sizes)
		//media.SizesPath = sizesPath
	}

	videoBuff, err := bufferFromFile(compressedFilePath)

	objectName := util.GetOriginalVideoPath(path)
	_, err = client.PutObject(ctx, bucketName, objectName, &videoBuff, int64(videoBuff.Len()), minio.PutObjectOptions{ContentType: mimetype.MP4})
	if err != nil {
		return 0, err
	}

	fmt.Println("Video uploaded successfully")
	//
	//duration, err := AudioVideoDuration(compressedFilePath)

	return duration, nil
}
