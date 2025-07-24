package music

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/image"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"net/http"
	"os"
)

type Music struct {
	client *minio.Client
}

func NewMusicClient(options sminio.Options) (sminio.MusicClient, error) {
	client, err := minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}
	return &Music{client: client}, nil
}

func (s *Music) UploadMusicPhotoReq(ctx context.Context, id int64, request *http.Request, key string) ([]util.Size, util.Err) {
	fmt.Println("")
	err := util.Validate(request, key, 50<<20, 1)
	if err.StatusCode > 0 {
		return nil, err
	}

	fmt.Println("id--->", id)
	f := request.MultipartForm.File[key][0]

	file, fErr := f.Open()
	if fErr != nil {
		return nil, util.UnsupportedErr(f.Filename)
	}
	defer file.Close()

	mimeType, detectErr := mimetype.Detect(file, mimetype.Images)
	if detectErr != nil || mimeType == "" {
		return nil, util.UnsupportedErr(f.Filename)
	}

	path := GetMusicPath(id)
	saveError := image.Save(ctx, s.client, mimeType, file, Size, path, util.MusicPhotoBucket)
	if saveError != nil {
		return nil, util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the image"}
	}

	return Size, util.Err{}
}

func (s *Music) UploadMusicPhoto(ctx context.Context, id int64, imagePath string) ([]util.Size, util.Err) {
	// Открываем файл с диска по пути imagePath
	file, fErr := os.Open(imagePath)
	if fErr != nil {
		return nil, util.UnsupportedErr(imagePath)
	}
	defer file.Close()

	// Определяем MIME-тип файла
	mimeType, detectErr := mimetype.Detect(file, mimetype.Images)
	if detectErr != nil || mimeType == "" {
		return nil, util.UnsupportedErr(imagePath)
	}

	// Генерируем путь для сохранения файла
	path := GetMusicPath(id)

	// Сохраняем файл, используя клиент MinIO
	saveError := image.Save(ctx, s.client, mimeType, file, Size, path, util.MusicPhotoBucket)
	if saveError != nil {
		return nil, util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the image"}
	}

	// Возвращаем размер сохраненного изображения
	return Size, util.Err{}
}

func (s *Music) UploadMusic(ctx context.Context, id int64, localDir string) (string, util.Err) {
	fmt.Println("")
	fmt.Printf("upload music in------")
	path := GetMusicPath(id)
	//localDir := "/home/user/Videos/Receive"
	saveError := Save(ctx, s.client, localDir, path, util.MusicPhotoBucket)
	if saveError != nil {
		return "", util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the music"}
	}

	return "success", util.Err{}
}
