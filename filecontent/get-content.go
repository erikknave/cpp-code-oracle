package filecontent

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetFileContent(objectKey string) (string, error) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	minioURL := os.Getenv("MINIO_URL")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-rack-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioAccessKey, minioSecretKey, "")),
	)
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(minioURL)
		o.UsePathStyle = true
	})

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func GetFilePartContent(fileKey string, start int64, end int64) (string, error) {

	fileContent, err := GetFileContent(fileKey)
	if err != nil {
		return "", err
	}
	filePartContent := fileContent[start:end]

	return filePartContent, nil
}
