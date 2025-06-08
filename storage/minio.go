package storage

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucket string) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, _ := client.BucketExists(ctx, bucket)
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
	}

	return &MinioClient{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, objectName string, filePath string) (string, error) {
	_, err := m.Client.FPutObject(ctx, m.Bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	url := "http://" + m.Client.EndpointURL().Host + "/" + m.Bucket + "/" + objectName
	return url, nil
}
