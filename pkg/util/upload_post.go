package util

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func UploadPostImage(c *fiber.Ctx, userID string, uploadDir string) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", nil
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	sanitizedFileName := sanitizeFileName(file.Filename)
	savePath := uploadDir + sanitizedFileName
	if err := c.SaveFile(file, savePath); err != nil {
		return "", err
	}

	return "https://raion-assessment.elginbrian.com" + savePath, nil
}