package image

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/image"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"net/http"
)

type User struct {
	client *minio.Client
}

func NewUserClient(options sminio.Options) (sminio.ImageClient, error) {
	client, err := minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}
	return &User{client: client}, nil
}

// RemoveUser id - image id

// UploadAvatar
// Err{StatusCode: http.StatusRequestEntityTooLarge, Message: "files is to large"}
func (s *User) UploadImage(ctx context.Context, request *http.Request, key, path string, Size []util.Size, buckedName string) util.Err {
	err := util.Validate(request, key, 50<<20, 1)
	if err.StatusCode > 0 {
		return err
	}

	f := request.MultipartForm.File[key][0]
	file, fErr := f.Open()
	if fErr != nil {
		return util.UnsupportedErr(f.Filename)
	}
	defer file.Close()

	mimeType, detectErr := mimetype.Detect(file, mimetype.Images)
	if detectErr != nil || mimeType == "" {
		return util.UnsupportedErr(f.Filename)
	}

	//path := GetAvatarPath(id)
	//fmt.Println("path avatar--->", path)
	saveError := image.Save(ctx, s.client, mimeType, file, Size, path, buckedName)
	if saveError != nil {
		return util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the image"}
	}

	return util.Err{}
}

func (s *User) RemoveImage(ctx context.Context, path, buckedName string) error {
	objectCh := s.client.ListObjects(ctx, buckedName, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		err := s.client.RemoveObject(ctx, buckedName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
