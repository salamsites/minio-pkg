package sminio

import (
	"context"
	"github.com/salamsites/minio-pkg/util"
	"net/http"
)

//type Some struct {
//	ChatClient  chat.Chat
//	MusicClient music.Music
//	FeedClient  feed.Feed
//	UserClient  user.User
//}

type UserClient interface {
	RemoveUser(ctx context.Context, id int64) error
	UploadAvatar(ctx context.Context, id int64, request *http.Request, key string) ([]util.Size, util.Err)
	RemoveAvatar(ctx context.Context, id int64) error
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

type ChatClient interface {
	UploadFile(ctx context.Context, id int64, request *http.Request, key string) (util.Media, util.Err)
}
