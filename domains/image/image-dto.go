package image

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

type SaveImageDto struct {
	File    *multipart.FileHeader `form:"file" binding:"required"`
	AltText string                `form:"alt_text"`
}

func (d SaveImageDto) IsValid() string {
	if d.File.Size > 5*1024*1024 {
		return "file size exceeds 5MB limit"
	}

	f, err := d.File.Open()
	if err != nil {
		return "unsupported file"
	}
	defer f.Close()
	_, format, err := image.Decode(f)
	if err != nil {
		return "unsupported image format"
	}

	if format != "jpeg" && format != "png" {
		return "only JPEG and PNG formats are supported"
	}
	return ""
}

type ImageResponse struct {
	ID      uint   `json:"id"`
	URL     string `json:"url"`
	AltText string `json:"alt_text"`
	UserID  uint   `json:"user_id"`
}
