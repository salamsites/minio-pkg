package sminiochat

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/client/audio"
	"github.com/salamsites/minio-pkg/image"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"github.com/salamsites/minio-pkg/video"
	"net/http"
	"path/filepath"
)

type Chat struct {
	client  *minio.Client
	tempDir string
}

func NewChatClient(options sminio.Options) (sminio.ChatClient, error) {
	client, err := minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}
	return &Chat{client: client, tempDir: options.TempDir}, nil
}

func (s *Chat) UploadFile(ctx context.Context, roomId int64, request *http.Request, key string) (util.Media, util.Err) {
	fmt.Printf("\n")
	fmt.Println("UploadFile chat")
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

	mediaMimeTypes := append(mimetype.Images, mimetype.Videos...)
	acceptedMimeTypes := append(mediaMimeTypes, mimetype.Files...)
	acceptedMimeTypes = append(acceptedMimeTypes, mimetype.Audios...)

	for i := range files {
		file, fErr := files[i].Open()
		if fErr != nil {
			isError = true
			errSave = util.UnsupportedErr(files[i].Filename)
			break
		}
		mimeType, detectErr := mimetype.Detect(file, acceptedMimeTypes)
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
	folderName := GetPathByTime(roomId)
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

		prefix, mime := mimetype.GetPrefixExt(mimeTypes[i])
		fmt.Println(prefix)
		fmt.Println(mimeTypes[i])

		path := GetPath(folderName, int64(i))
		switch prefix {
		case mimetype.PrefixImage:
			saveError := image.Save(ctx, s.client, mimeTypes[i], file, Size, path, util.ChatBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the image"}
				break
			}
			img := util.FeedResultTypeImage{
				Path: path,
				Type: mimetype.PrefixImage,
				Mime: ".webp",
			}
			content = append(content, img)
		case mimetype.PrefixVideo:
			duration, saveError := video.Save(ctx, s.client, s.tempDir, mimeTypes[i], file, files[i].Filename, Size, path, util.ChatBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the video"}
				break
			}

			ex := filepath.Ext(files[i].Filename)

			img := util.FeedResultTypeVideo{
				Path:     path,
				Type:     mimetype.PrefixVideo,
				Duration: duration,
				Mime:     ex,
			}
			content = append(content, img)
			fmt.Println(fmt.Sprintf("%d sec", duration))
			break
		case mimetype.PrefixAudio:
			duration, saveError := audio.Save(ctx, s.client, s.tempDir, mimeTypes[i], file, files[i].Filename, Size, path, util.ChatBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the audio"}
				break
			}

			ex := filepath.Ext(files[i].Filename)

			img := util.FeedResultTypeAudio{
				Path:     path,
				Type:     mimetype.PrefixAudio,
				Duration: duration,
				Mime:     ex,
			}
			content = append(content, img)
			fmt.Println(fmt.Sprintf("%d sec", duration))
			break

		default:
			saveError := SaveFile(ctx, s.client, mimeTypes[i], file, mime, path, util.ChatBucket)
			if saveError != nil {
				isError = true
				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occurred while saving the file"}
				break
			}
			ex := filepath.Ext(files[i].Filename)

			img := util.FeedResultTypeFile{
				Path:     path,
				Type:     mimetype.PrefixFile,
				FileSize: files[i].Size,
				Mime:     ex,
			}
			content = append(content, img)
			break
		}

		err := file.Close()
		if err != nil {
			fmt.Println("error closing file", err)
		}
	}
	fmt.Println("content---->", content)

	if isError {
		return result, errSave
	}

	result.Sizes = Size
	result.Content = content
	return result, util.Err{}
}

//func (s *chat) UploadFileOld(ctx context.Context, id int64, request *http.Request, key string) (util.Media, util.Err) {
//	fmt.Printf("\n")
//	fmt.Println("UploadFile")
//	result := util.Media{}
//	err := util.Validate(request, key, 50<<20, 10)
//	if err.StatusCode > 0 {
//		return result, err
//	}
//
//	files := request.MultipartForm.File[key]
//	if files == nil {
//		return result, util.Err{}
//	}
//
//	var errSave util.Err
//	isError := false
//
//	var mimeTypes []string
//
//	for i := range files {
//		file, fErr := files[i].Open()
//		if fErr != nil {
//			isError = true
//			errSave = util.UnsupportedErr(files[i].Filename)
//			break
//		}
//
//		mime, err := gabrielvasile.DetectReader(file)
//		if err != nil {
//			panic(err)
//		}
//
//		mimeType := mime.String()
//
//		fmt.Println("MIME type:", mimeType)
//
//		//mimeType, detectErr := mimetype.Detect(file)
//		//if detectErr != nil || mimeType == "" {
//		//	isError = true
//		//	errSave = util.UnsupportedErr(files[i].Filename)
//		//	break
//		//}
//
//		mimeTypes = append(mimeTypes, mimeType)
//		err = file.Close()
//		if err != nil {
//			fmt.Println("error closing file", err)
//		}
//	}
//
//	if isError {
//		return result, errSave
//	}
//
//	var content []interface{}
//
//	for i := range files {
//		if isError {
//			break
//		}
//		file, fErr := files[i].Open()
//		if fErr != nil {
//			isError = true
//			errSave = util.UnsupportedErr(files[i].Filename)
//			break
//		}
//
//		// Генерируем путь сохранения
//		path := GetPath(id, int64(i))
//
//		switch {
//		case mimetype.IsImage(mimeTypes[i]):
//			saveError := image.Save(ctx, s.client, mimeTypes[i], file, Size, path, util.FeedBucket)
//			if saveError != nil {
//				isError = true
//				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the image"}
//				break
//			}
//			img := util.FeedResultTypeImage{
//				Id:   int64(i),
//				Type: mimetype.PrefixImage,
//			}
//			content = append(content, img)
//		case mimetype.IsVideo(mimeTypes[i]):
//			saveError := video.Save(ctx, s.client, mimeTypes[i], file, files[i].Filename, Size, path, util.FeedBucket)
//			if saveError != nil {
//				isError = true
//				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the video"}
//				break
//			}
//			vid := util.FeedResultTypeVideo{
//				Id:   int64(i),
//				Type: mimetype.PrefixVideo,
//			}
//			content = append(content, vid)
//		default:
//			saveError := generic.Save(ctx, s.client, mimeTypes[i], file, files[i].Filename, path, util.FeedBucket)
//			if saveError != nil {
//				isError = true
//				errSave = util.Err{StatusCode: http.StatusBadRequest, Message: "error occured while saving the file"}
//				break
//			}
//			doc := util.FeedResultTypeFile{
//				Id:   int64(i),
//				Type: mimeTypes[i],
//			}
//			content = append(content, doc)
//		}
//
//		err := file.Close()
//		if err != nil {
//			fmt.Println("error closing file", err)
//		}
//	}
//
//	if isError {
//		return result, errSave
//	}
//
//	result.Sizes = Size
//	result.Content = content
//	return result, util.Err{}
//}
