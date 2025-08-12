package feed

import "fmt"

func GetFeedPath(userId int64, feedType string, feedId int64, mimeType string, generateId int64) string {
	return fmt.Sprintf("%d/%s/%d/%s/%d", userId, feedType, feedId, mimeType, generateId)
}
