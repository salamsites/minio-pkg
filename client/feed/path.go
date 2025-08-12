package feed

import (
	"fmt"
	"time"
)

func GetFeedPath(userId int64, feedType string, feedId int64, mimeType string) string {
	folderByTime := currentTime()
	return fmt.Sprintf("%d/%s/%d/%s/%s", userId, feedType, feedId, mimeType, folderByTime)
}

func currentTime() string {
	currentTime := time.Now().UnixMilli()

	// Преобразуем миллисекунды в строку
	folderByTime := fmt.Sprintf("%d", currentTime)
	return folderByTime
}
