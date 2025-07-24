package test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

const url = "http://localhost" + PORT

func TestMain(m *testing.M) {
	log.Println("TestMain started")
	code := m.Run()
	os.Exit(code)
}

func TestSetMimeTypeAll(t *testing.T) {
	// Open the file
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files := []string{
		//pwd + "/test_files/2.jpg",
		pwd + "/test_files/123.mp3",
		//pwd + "/test_files/ALGEBRA_10.pdf",
		//pwd + "/test_files/4955.pdf",
		//pwd + "/test_files/1.HEIC",
		//pwd + "/test_files/873.tiff",
		//pwd + "/test_files/maxresdefault.bmp",
		//pwd + "/test_files/owl.svg",
		//pwd + "/test_files/patrol_icon.ico",
		//pwd + "/test_files/69u.gif",
		//pwd + "/test_files/IMG_4458.heic",
		//pwd + "/test_files/20240530_163214.jpg",
		//pwd + "/test_files/1.video.mov",
		//pwd + "/test_files/upx.mp4",
		//pwd + "/test_files/insta/news-video/1_5-minut_video/Dinamica Zbor.mp4",
	}

	//files := []string{
	//	pwd + "/test_files/",
	//}

	//	fmt.Println("files: ", files)

	// Create a new buffer to store the request body
	body := &bytes.Buffer{}

	// Create a new multipart writer
	writer := multipart.NewWriter(body)

	// Iterate over the files
	for _, filename := range files {
		// Open the file
		file, err := os.Open(filename)
		assert.NoError(t, err)

		defer file.Close()

		// Create a new form file field
		part, err := writer.CreateFormFile(KEY, filename)
		assert.NoError(t, err)

		// Copy the file content to the form file field
		_, err = io.Copy(part, file)
		assert.NoError(t, err)

	}

	// Close the multipart writer
	writer.Close()

	//fmt.Println("body-->", body)

	req, err := http.NewRequest("POST", url+UploadImageURL, body)
	assert.NoError(t, err)

	// Set the Content-Type header to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	// Print the response status code
	fmt.Println("Response status:", resp.Status)
	fmt.Printf("Response message: %v\n", resp.Body)
}
