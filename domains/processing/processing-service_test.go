package processing

import (
	"bytes"
	internalImg "image"
	"image/color"
	"image/jpeg"
	"testing"

	"github.com/disintegration/imaging"
)

func createTestImage(width, height int) ([]byte, error) {
	img := internalImg.NewRGBA(internalImg.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	return buf.Bytes(), err
}

func TestApplyTransformations_Resize(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Resize: &ResizeDTO{
			Width:  50,
			Height: 50,
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("Failed to decode result image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected dimensions 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestApplyTransformations_NoTransformations(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{} // Empty DTO

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestApplyTransformations_InvalidFormat(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		FormatConversion: &FormatConversionDTO{
			Format: "invalid",
		},
	}

	_, err = applyTransformations(originalImg, dto)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestApplyTransformations_Crop(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Crop: &CropDTO{
			X:      10,
			Y:      10,
			Width:  50,
			Height: 50,
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("Failed to decode result image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected dimensions 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestApplyTransformations_Rotate(t *testing.T) {
	originalImg, err := createTestImage(100, 50)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Rotate: &RotateDTO{
			Angle: 90,
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("Failed to decode result image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 100 {
		t.Errorf("Expected dimensions 50x100, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestApplyTransformations_FlipHorizontal(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Flip: &FlipDTO{
			Direction: "horizontal",
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestApplyTransformations_FlipVertical(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Flip: &FlipDTO{
			Direction: "vertical",
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestApplyTransformations_FormatConversion(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		FormatConversion: &FormatConversionDTO{
			Format: "png",
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestApplyTransformations_Filter(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Filter: &FilterDTO{
			Saturation: 50,
			Brightness: 10,
			Contrast:   20,
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestApplyTransformations_Watermark(t *testing.T) {
	originalImg, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	dto := TransformationDTO{
		Watermark: &WatermarkDTO{
			Text: "Test",
			Position: Point{
				X: 10,
				Y: 10,
			},
			Opacity: 50,
		},
	}

	result, err := applyTransformations(originalImg, dto)
	if err != nil {
		t.Fatalf("applyTransformations failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}
