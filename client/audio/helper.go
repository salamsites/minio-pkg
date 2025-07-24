package audio

import (
	"bytes"
	"fmt"
	"github.com/vansante/go-ffprobe"
	"io"
	"os"
	"os/exec"
	"time"
)

func compressToMP4(originalFilePath, compressedFilePath string, crf int) error {

	cmd := exec.Command("ffmpeg", "-y", "-i", originalFilePath, "-c:v", "libx264", "-crf", fmt.Sprintf("%d", crf), "-preset", "medium", "-c:a", "aac", "-b:a", "128k", "-movflags", "+faststart", compressedFilePath)
	//cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func bufferFromFile(filePath string) (bytes.Buffer, error) {
	var result bytes.Buffer
	outputFile, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("failed to open output file: %v", err)
	}
	defer outputFile.Close()

	if _, err := io.Copy(&result, outputFile); err != nil {
		return result, fmt.Errorf("failed to copy output file: %v", err)
	}

	return result, nil
}

func audioVideoDuration(filePath string) (int64, error) {
	data, err := ffprobe.GetProbeData(filePath, 1*time.Minute)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	duration := data.Format.Duration().Milliseconds()
	return duration, nil
}
