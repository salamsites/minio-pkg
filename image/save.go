package image

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/salamsites/minio-pkg/mimetype"
	"github.com/salamsites/minio-pkg/util"
	"image"
	"io"
	"mime/multipart"
)

func Save(ctx context.Context, client *minio.Client, mimeType string, file multipart.File, sizes []util.Size, path, bucketName string) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	var resultSizes []util.Size
	var resizeErr error
	img, err := decode(file, mimeType)
	if err != nil {
		fmt.Println("decode  err")
		return err
	}

	for _, size := range sizes {

		var resizedImg *image.NRGBA
		if size.Height == 0 {
			resizedImg = imaging.Resize(img, size.Width, size.Height, imaging.Lanczos)
		} else {
			resizedImg = imaging.Fill(img, size.Width, size.Height, imaging.Center, imaging.Lanczos)
		}

		var buf bytes.Buffer
		if err := webp.Encode(&buf, resizedImg, &webp.Options{Quality: float32(size.Quality)}); err != nil {
			fmt.Println("11. error:", err)
			resizeErr = err
			break
		}

		objectName := util.GetImagePath(path, size)
		fmt.Println("objectname----->", objectName)
		_, err = client.PutObject(ctx, bucketName, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: mimetype.WEBP})
		if err != nil {
			resizeErr = err
			fmt.Println("12. error:", err)
			break
		}
		resultSizes = append(resultSizes, size)
	}

	if resizeErr != nil {
		for _, size := range resultSizes {
			objectName := util.GetImagePath(path, size)
			err := client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
			if err != nil {
				// TODO
				fmt.Println(err)
			}
		}

		return resizeErr
	}

	// TODO error berse yokarda save edenlerni pozmaly
	//save original quality
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: true}); err != nil {
		return err
	}
	objectName := util.GetOriginalImagePath(path)
	_, err = client.PutObject(ctx, bucketName, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: mimetype.WEBP})
	if err != nil {
		return err
	}

	return nil
}
