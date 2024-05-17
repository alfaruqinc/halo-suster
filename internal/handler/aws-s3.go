package handler

import (
	"context"
	"health-record/internal/domain"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AWSS3Handler interface {
	UploadImage() fiber.Handler
}

type awsS3 struct{}

func NewAWSS3() AWSS3Handler {
	return &awsS3{}
}

var (
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsS3BucketName    = os.Getenv("AWS_S3_BUCKET_NAME")
	awsRegion          = os.Getenv("AWS_REGION")
)

func (s *awsS3) UploadImage() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			err := domain.NewErrInternalServerError(err.Error())
			return ctx.Status(err.Status()).JSON(err)
		}

		fileSize := file.Size
		if fileSize < 80000 {
			err := domain.NewErrBadRequest("filsize should more than 10KB")
			return ctx.Status(err.Status()).JSON(err)
		}
		if fileSize > 8e7 {
			err := domain.NewErrBadRequest("filsize should less than 10MB")
			return ctx.Status(err.Status()).JSON(err)
		}

		clientFile, err := file.Open()
		if err != nil {
			err := domain.NewErrInternalServerError(err.Error())
			return ctx.Status(err.Status()).JSON(err)
		}
		defer clientFile.Close()
		buff := make([]byte, 512)
		if _, err := clientFile.Read(buff); err != nil {
			err := domain.NewErrInternalServerError(err.Error())
			return ctx.Status(err.Status()).JSON(err)
		}

		onlyImg := []string{"image/jpeg", "image/jpg"}
		formatFile := http.DetectContentType(buff)
		if !slices.Contains(onlyImg, formatFile) {
			err := domain.NewErrBadRequest("image format is not allowed")
			return ctx.Status(err.Status()).JSON(err)
		}
		fileExtension := filepath.Ext(file.Filename)

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")))
		if err != nil {
			err := domain.NewErrInternalServerError(err.Error())
			return ctx.Status(err.Status()).JSON(err)
		}

		client := s3.NewFromConfig(cfg)

		client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(awsS3BucketName),
			Key:    aws.String(uuid.New().String() + fileExtension),
			Body:   clientFile,
		})

		respUpload := domain.NewStatusOK("success upload image", map[string]string{"imageUrl": "test"})

		return ctx.Status(respUpload.Status()).JSON(respUpload)
	}
}
