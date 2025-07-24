package util

import (
	"fmt"
)

func GetOriginalVideoPath(path string) string {
	return fmt.Sprintf("%s/original.mp4", path)
}

func GetVideoPath(path string, size Size) string {
	return fmt.Sprintf("%s/%dx%d.webp", path, size.Width, size.Height)
}

func GetImagePath(path string, size Size) string {
	return fmt.Sprintf("%s/%dx%d.webp", path, size.Width, size.Height)
}

func GetOriginalImagePath(path string) string {
	return fmt.Sprintf("%s/original.webp", path)
}

func GetOriginalAudioPath(path string) string {
	return fmt.Sprintf("%s/original.mp3", path)
}
