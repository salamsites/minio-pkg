package user

import (
	"context"
	"fmt"
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

func NewUserClient(options sminio.Options) (sminio.UserClient, error) {
	client, err := minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}
	return &User{client: client}, nil
}

// RemoveUser id - user id
func (s *User) RemoveUser(ctx context.Context, id int64) error {
	// TODO remove avatar
	// TODO remove feed
	return nil
}

// UploadAvatar
// Err{StatusCode: http.StatusRequestEntityTooLarge, Message: "files is to large"}
func (s *User) UploadAvatar(ctx context.Context, id int64, request *http.Request, key string) ([]util.Size, util.Err) {
	err := util.Validate(request, key, 50<<20, 1)
	if err.StatusCode > 0 {
		return nil, err
	}

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

	path := GetAvatarPath(id)
	fmt.Println("path avatar--->", path)
	saveError := image.Save(ctx, s.client, mimeType, file, Size, path, util.AvatarBucket)
	if saveError != nil {
		return nil, util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the image"}
	}

	return Size, util.Err{}
}

func (s *User) RemoveAvatar(ctx context.Context, id int64) error {
	objectCh := s.client.ListObjects(ctx, util.AvatarBucket, minio.ListObjectsOptions{
		Prefix:    GetAvatarDir(id),
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		err := s.client.RemoveObject(ctx, util.AvatarBucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
