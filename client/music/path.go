package music

import "fmt"

func GetMusicPath(id int64) string {
	return fmt.Sprintf("%d", id)
}
