package image

import (
	"bytes"
	"fmt"
	_ "github.com/Kodeworks/golang-image-ico"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"io"
	"mime/multipart"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// convertHeicToJpg converts a HEIC image from a multipart.File to an image.Image
func convertHeicToJpg(file multipart.File) (image.Image, error) {
	// Read the file into a buffer
	fmt.Println("converting heic to jpg")
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	// Create an io.ReaderAt for the buffer
	readerAt := bytes.NewReader(buf.Bytes())

	// Decode the HEIC image
	img, err := goheif.Decode(readerAt)
	if err != nil {
		return nil, err
	}

	// Extract EXIF data
	exifData, err := goheif.ExtractExif(readerAt)
	if err != nil {
		return nil, err
	}

	// Adjust the image orientation based on EXIF data
	img, err = adjustOrientation(img, exifData)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// adjustOrientation adjusts the image orientation based on EXIF data
func adjustOrientation(img image.Image, exifData []byte) (image.Image, error) {
	if len(exifData) == 0 {
		return img, nil
	}

	exifReader := bytes.NewReader(exifData)
	exifInfo, err := exif.Decode(exifReader)
	if err != nil {
		return nil, err
	}

	orientation, err := exifInfo.Get(exif.Orientation)
	if err != nil {
		return img, nil
	}

	orientValue, err := orientation.Int(0)
	if err != nil {
		return img, nil
	}

	switch orientValue {
	case 3:
		return imaging.Rotate180(img), nil
	case 6:
		return imaging.Rotate270(img), nil
	case 8:
		return imaging.Rotate90(img), nil
	default:
		return img, nil
	}
}

// convertFormatToJpg converts a GIF, BMP, TIFF  image from a multipart.File to an image.Image
func convertFormatToJpg(file multipart.File, format string) (image.Image, error) {
	// Чтение файла в буфер
	fmt.Println("barde")
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		fmt.Println("error copy")
		return nil, err
	}

	// Создание io.Reader из буфера
	reader := bytes.NewReader(buf.Bytes())

	// Декодирование изображения в зависимости от формата
	var img image.Image
	fmt.Println("format--->", format)
	switch format {
	case "image/gif":
		img, err = gif.Decode(reader)
	case "image/bmp":
		img, err = bmp.Decode(reader)
	case "image/tiff":
		img, err = tiff.Decode(reader)
	case "image/svg+xml":
		img, err = decodeSVG(reader)
	case "image/x-icon":
		img, err = decodeICO(reader)
	}
	if err != nil {
		fmt.Println("decode error format", format)
		fmt.Println("error----->", err)
		return nil, err
	}

	return img, nil
}

func decodeSVG(reader io.Reader) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(reader)
	if err != nil {
		return nil, err
	}

	width := int(icon.ViewBox.W)
	height := int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	icon.SetTarget(0, 0, float64(width), float64(height))

	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1)

	return img, nil
}

// decodeICO декодирует ICO изображение
func decodeICO(reader io.Reader) (image.Image, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}
