package feed

import "fmt"

func GetFeedPath(id, feedId int64) string {
	return fmt.Sprintf("%d/%d", id, feedId)
}
