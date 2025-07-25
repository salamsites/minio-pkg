package feed

import "fmt"

func GetFeedPath(id, feedId int64, mimeType string, generateId int64) string {
	return fmt.Sprintf("%d/%d/%s/%d", id, feedId, mimeType, generateId)
}
