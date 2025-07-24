package sminiochat

import (
	"fmt"
	"time"
)

func GetPathByTime(id int64) string {
	folderByTime := currentTime()
	return fmt.Sprintf("%d/%s", id, folderByTime)
}

func GetPath(folderName string, chatId int64) string {
	return fmt.Sprintf("%s/%d", folderName, chatId)
}

func currentTime() string {
	currentTime := time.Now().UnixMilli()

	// Преобразуем миллисекунды в строку
	folderByTime := fmt.Sprintf("%d", currentTime)
	return folderByTime
}
