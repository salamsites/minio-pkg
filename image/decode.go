package image

import (
	"fmt"
	"github.com/salamsites/minio-pkg/mimetype"
	"image"
	"io"
	"mime/multipart"
)

func decode(file multipart.File, mimeType string) (image.Image, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("seek", err)
	}

	fmt.Println("Decoding image")
	fmt.Println("mime---", mimeType)
	switch mimeType {
	case mimetype.HEIF, mimetype.HEIC:
		return convertHeicToJpg(file)
	case mimetype.GIF, mimetype.BMP, mimetype.TIFF, mimetype.SVG, mimetype.ICO:
		return convertFormatToJpg(file, mimeType)

	default:
		img, _, err := image.Decode(file)
		return img, err
	}
}
