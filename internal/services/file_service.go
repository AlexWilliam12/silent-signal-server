package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Calculate the file hash and save the file on directory 'uploads'.
// Return three values (hash, filename, error)
func SaveFile(fileHeader *multipart.FileHeader) (string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	path := filepath.Join("uploads", filename)

	out, err := os.Create(path)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", "", err
	}

	hash, err := calculateHash(out)
	if err != nil {
		return "", "", err
	}

	return hash, filename, nil
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

func calculateHash(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}
