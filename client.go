package sminio

import (
	"context"
	"github.com/salamsites/minio-pkg/util"
	"net/http"
)

//type Some struct {
//	ChatClient  files.Chat
//	MusicClient music.Music
//	FeedClient  feed.Feed
//	UserClient  image.User
//}

type ImageClient interface {
	UploadImage(ctx context.Context, request *http.Request, key, path string, Size []util.Size, buckedName string) util.Err
	RemoveImage(ctx context.Context, path, buckedName string) error
}

type MusicClient interface {
	UploadMusicPhoto(ctx context.Context, id int64, path string) ([]util.Size, util.Err)
	UploadMusic(ctx context.Context, id int64, path string) (string, util.Err)
}

type FeedClient interface {
	UploadFeed(ctx context.Context, userid, feedId int64, request *http.Request, key string) (util.Media, util.Err)
	//RemoveFeed(ctx context.Context) error
	//
	//UploadStories(ctx context.Context) error
	//RemoveStories(ctx context.Context) error
}

type FileClient interface {
	UploadFile(ctx context.Context, request *http.Request, key string, path []string, Size []util.Size, buckedName string) (util.Media, util.Err)
}
