package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	path := filepath.Join("uploads", filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}

	return filename, nil
}

func DeleteFile(filename string) error {
	if err := os.Remove("uploads/" + filename); err != nil {
		return err
	}
	return nil
}

func BuildFileURL(filename string) string {
	var builder strings.Builder
	builder.WriteString("http://")
	builder.WriteString(os.Getenv("SERVER_HOST"))
	builder.WriteString(":")
	builder.WriteString(os.Getenv("SERVER_PORT"))
	builder.WriteString("/uploads/")
	builder.WriteString(filename)
	return builder.String()
}

func EncodedImgThumbnail(file *multipart.File) (string, error) {
	img, _, err := image.Decode(*file)
	if err != nil {
		return "", err
	}

	resizedImg := imaging.Resize(img, 100, 0, imaging.Lanczos)
	blurredImg := imaging.Blur(resizedImg, 2.0)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, blurredImg, nil); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

func EncodedVideoThumbnail(file *multipart.File, videoType string) (string, error) {
	videoFile, err := os.CreateTemp("temp", "video-*."+videoType)
	if err != nil {
		return "", err
	}
	defer os.Remove(videoFile.Name())

	_, err = io.Copy(videoFile, *file)
	if err != nil {
		return "", err
	}

	videoPath := videoFile.Name()

	thumbnailFile, err := os.CreateTemp("temp", "thumbnail-*.jpg")
	if err != nil {
		return "", err
	}
	defer os.Remove(thumbnailFile.Name())
	thumbnailPath := thumbnailFile.Name()

	cmd := exec.Command("ffmpeg", "-ss", "00:00:00", "-i", videoPath, "-s", "320x240", "-vframes", "1", thumbnailPath)
	if err = cmd.Run(); err != nil {
		return "", err
	}

	img, _, err := image.Decode(thumbnailFile)
	if err != nil {
		return "", err
	}

	resizedImg := imaging.Resize(img, 100, 0, imaging.Lanczos)
	blurredImg := imaging.Blur(resizedImg, 2.0)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, blurredImg, nil); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}
