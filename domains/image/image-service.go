package image

import "gorm.io/gorm"

type ImageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{db: db}
}

func (s *ImageService) SaveImage(image *Image) (*Image, error) {
	if err := s.db.Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

func (s *ImageService) GetImageByID(id uint) (*Image, error) {
	var image Image
	if err := s.db.Preload("User").First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (s *ImageService) ListImagesByUser(userID uint) ([]Image, error) {
	var images []Image
	if err := s.db.Where("user_id = ?", userID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (s *ImageService) DeleteImage(id uint) error {
	if err := s.db.Delete(&Image{}, id).Error; err != nil {
		return err
	}
	return nil
}
