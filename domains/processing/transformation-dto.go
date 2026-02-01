package processing

type TransformationDTO struct {
	Resize           *ResizeDTO           `json:"resize,omitempty"`
	Crop             *CropDTO             `json:"crop,omitempty"`
	Rotate           *RotateDTO           `json:"rotate,omitempty"`
	Flip             *FlipDTO             `json:"flip,omitempty"`
	Watermark        *WatermarkDTO        `json:"watermark,omitempty"`
	FormatConversion *FormatConversionDTO `json:"format_conversion,omitempty"`
	Filter           *FilterDTO           `json:"filter,omitempty"`
}

type ResizeDTO struct {
	Width  int `json:"width" binding:"required,min=1"`
	Height int `json:"height" binding:"required,min=1"`
}

type CropDTO struct {
	X      int `json:"x" binding:"required,min=0"`
	Y      int `json:"y" binding:"required,min=0"`
	Width  int `json:"width" binding:"required,min=1"`
	Height int `json:"height" binding:"required,min=1"`
}

type RotateDTO struct {
	Angle float64 `json:"angle" binding:"required"`
}

type FlipDTO struct {
	Direction string `json:"direction" binding:"required,oneof=horizontal vertical"`
}

type WatermarkDTO struct {
	Text     string `json:"text" binding:"required"`
	Position Point  `json:"position" binding:"required,oneof=top-left top-right bottom-left bottom-right center"`
	Opacity  int    `json:"opacity" binding:"required,min=0,max=100"`
}

type Point struct {
	X int `json:"x" binding:"required"`
	Y int `json:"y" binding:"required"`
}

type FormatConversionDTO struct {
	Format string `json:"format" binding:"required,oneof=jpeg png webp tiff bmp gif"`
}

type FilterDTO struct {
	Saturation int `json:"saturation" binding:"required,min=-100,max=100"`
	Brightness int `json:"brightness" binding:"required,min=-100,max=100"`
	Contrast   int `json:"contrast" binding:"required,min=-100,max=100"`
}
