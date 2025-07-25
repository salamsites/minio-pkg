package feed

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/image"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"github.com/salamsites/minio-pkg/video"
	"net/http"
)

type Feed struct {
	client  *minio.Client
	tempDir string
}

func NewFeedClient(options sminio.Options) (sminio.FeedClient, error) {
	client, err := minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}
	return &Feed{client: client, tempDir: options.TempDir}, nil
}

func (s *Feed) UploadFeed(ctx context.Context, userid, feedId int64, request *http.Request, key string) (util.Media, util.Err) {
	fmt.Printf("\n")
	fmt.Println("UploadFeed")
	result := util.Media{}
	err := util.Validate(request, key, 50<<20, 10)
	if err.StatusCode > 0 {
		return result, err
	}

	files := request.MultipartForm.File[key]
	if files == nil {
		return result, util.Err{}
	}

	var errSave util.Err
	isError := false

	var mimeTypes []string

	feedMimeTypes := append(mimetype.Images, mimetype.Videos...)

	for i := range files {
		file, fErr := files[i].Open()
		if fErr != nil {
			isError = true
			errSave = util.UnsupportedErr(files[i].Filename)
			break
		}

		mimeType, detectErr := mimetype.Detect(file, feedMimeTypes)
		if detectErr != nil || mimeType == "" {
			isError = true
			errSave = util.UnsupportedErr(files[i].Filename)
			break
		}
		fmt.Println("mimetype->", mimeType)
		mimeTypes = append(mimeTypes, mimeType)
		err := file.Close()
		if err != nil {
			//TODO
			fmt.Println("error closing file", err)
		}
	}

	if isError {
		return result, errSave
	}

	var content []interface{}

	for i := range files {
		if isError {
			break
		}
		file, fErr := files[i].Open()
		if fErr != nil {
			isError = true
			errSave = util.UnsupportedErr(files[i].Filename)
			break
		}

		prefix, _ := mimetype.GetPrefixExt(mimeTypes[i])
		switch prefix {
		case mimetype.PrefixImage:
			path := GetFeedPath(userid, feedId, mimetype.PrefixImage, int64(i))
			saveError := image.Save(ctx, s.client, mimeTypes[i], file, Size, path, util.FeedBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the image"}
				break
			}
			img := util.FeedResultTypeImage{
				Path: path,
				Type: mimetype.PrefixImage,
				Mime: mimeTypes[i],
			}
			content = append(content, img)
		case mimetype.PrefixVideo:
			path := GetFeedPath(userid, feedId, mimetype.PrefixVideo, int64(i))
			duration, saveError := video.Save(ctx, s.client, s.tempDir, mimeTypes[i], file, files[i].Filename, Size, path, util.FeedBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the video"}
				break
			}

			video := util.FeedResultTypeVideo{
				Path:     path,
				Type:     mimetype.PrefixVideo,
				Duration: duration,
				Mime:     mimeTypes[i],
			}
			content = append(content, video)
			break
		}

		err := file.Close()
		if err != nil {
			//TODO
			fmt.Println("error closing file", err)
		}
	}

	if isError {
		// TODO
	}

	result.Sizes = Size
	result.Content = content
	return result, util.Err{}
}
