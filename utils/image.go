package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadImage(c *fiber.Ctx, formKey string) (string, error) {
	file, err := c.FormFile(formKey)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// İçeriğe göre hash üretelim
	hash := sha1.New()
	io.Copy(hash, src)

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%x_%d%s", hash.Sum(nil), time.Now().UnixNano(), ext)
	path := filepath.Join("uploads", filename)

	// Kaydetmek için baştan dosya oku
	src.Seek(0, io.SeekStart)
	if err := c.SaveFile(file, path); err != nil {
		return "", err
	}

	return path, nil
}

func DeleteImage(path string) error {
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
