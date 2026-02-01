package processing

import (
	"bytes"
	"context"
	"errors"
	"vixel/domains/image"

	internalImg "image"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"gorm.io/gorm"
)

type ProcessingService struct {
	db            *gorm.DB
	uploadService *image.UploadService
}

func NewProcessingService(db *gorm.DB, uploadService *image.UploadService) *ProcessingService {
	return &ProcessingService{db: db, uploadService: uploadService}
}

func (s *ProcessingService) TransformImage(ctx context.Context, imageID string, dto TransformationDTO) (string, error) {
	var res image.Image
	if err := s.db.First(&res, "id = ?", imageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("image not found")
		}
		return "", err
	}

	img, err := s.uploadService.GetImageByUrl(ctx, res.URL)
	if err != nil {
		return "", err
	}

	transformedImg, err := applyTransformations(img, dto)
	if err != nil {
		return "", err
	}

	uploadedURL, err := s.uploadService.UploadImageFromBytes(ctx, transformedImg, "image/jpeg")
	if err != nil {
		return "", err
	}

	if err := s.uploadService.DeleteImage(ctx, res.URL); err != nil {
		return "", err
	}

	return uploadedURL, nil
}

func applyTransformations(img []byte, dto TransformationDTO) ([]byte, error) {
	src, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	if dto.Resize != nil {
		resized := imaging.Resize(src, dto.Resize.Width, dto.Resize.Height, imaging.Lanczos)
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, resized, imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.Crop != nil {
		cropped := imaging.Crop(src, internalImg.Rect(dto.Crop.X, dto.Crop.Y, dto.Crop.Width, dto.Crop.Height))
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, cropped, imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.Rotate != nil {
		rotated := imaging.Rotate(src, dto.Rotate.Angle, internalImg.Transparent)
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, rotated, imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.Flip != nil {
		var flipped *internalImg.NRGBA
		if dto.Flip.Direction == "horizontal" {
			flipped = imaging.FlipH(src)
		} else {
			flipped = imaging.FlipV(src)
		}
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, flipped, imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.FormatConversion != nil {
		buf := new(bytes.Buffer)
		var format imaging.Format
		switch dto.FormatConversion.Format {
		case "jpeg":
			format = imaging.JPEG
		case "png":
			format = imaging.PNG
		case "tiff":
			format = imaging.TIFF
		case "bmp":
			format = imaging.BMP
		case "gif":
			format = imaging.GIF
		default:
			return nil, errors.New("unsupported format")
		}
		err = imaging.Encode(buf, src, format)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.Filter != nil {
		var filtered *internalImg.NRGBA
		filtered = imaging.AdjustSaturation(src, float64(dto.Filter.Saturation))
		filtered = imaging.AdjustBrightness(filtered, float64(dto.Filter.Brightness))
		filtered = imaging.AdjustContrast(filtered, float64(dto.Filter.Contrast))
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, filtered, imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	if dto.Watermark != nil {
		position := internalImg.Pt(dto.Watermark.Position.X, dto.Watermark.Position.Y)
		buf := new(bytes.Buffer)
		dc := gg.NewContext(200, 50)
		dc.SetRGBA(1, 1, 1, 0.5)
		dc.DrawStringAnchored(dto.Watermark.Text, 100, 25, 0.5, 0.5)
		watermark := dc.Image().(*internalImg.NRGBA)
		err = imaging.Encode(buf, imaging.Overlay(src, watermark, position, 0.5), imaging.JPEG)
		if err != nil {
			return nil, err
		}
		img = buf.Bytes()
	}

	return img, nil
}
