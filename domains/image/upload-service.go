package image

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"time"
	"vixel/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadService struct {
	client *minio.Client
}

func NewUploadService() *UploadService {
	c, err := minio.New(config.Config.MINIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.MINIOAccessKey, config.Config.MINIOSecretKey, ""),
		Secure: config.Config.MINIOUseSSL,
	})
	if err != nil {
		panic(err)
	}
	return &UploadService{client: c}
}

func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	bucketName := config.Config.MINIOBucketName
	if ok, err := s.client.BucketExists(ctx, bucketName); err != nil || !ok {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: config.Config.MINIORegion,
		}); err != nil {
			return "", err
		}

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, bucketName)
		if err := s.client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
			return "", err
		}
	}

	imageName := fmt.Sprintf("vixel-%v-%v", rand.Intn(999999), time.Now().UnixNano())
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	uploadInfo, err := s.client.PutObject(ctx, bucketName, imageName, f, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	protocol := "http"
	if config.Config.MINIOUseSSL {
		protocol = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/%s/%s", protocol, config.Config.MINIOEndpoint, bucketName, uploadInfo.Key)

	return imageURL, nil
}
