package user

import "fmt"

func GetAvatarPath(id int64) string {
	return fmt.Sprintf("%d", id)
}

func GetAvatarDir(id int64) string {
	return fmt.Sprintf("%d", id)
}
