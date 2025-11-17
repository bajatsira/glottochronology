package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinio() (*MinioClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	access := os.Getenv("MINIO_ACCESS")
	secret := os.Getenv("MINIO_SECRET")
	bucket := os.Getenv("MINIO_BUCKET")
	if endpoint == "" || access == "" || secret == "" || bucket == "" {
		return nil, fmt.Errorf("minio: missing env variables (MINIO_*)")
	}

	useSSL := false
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioClient{Client: client, Bucket: bucket}, nil
}

func (m *MinioClient) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.Client.PutObject(ctx, m.Bucket, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (m *MinioClient) Remove(ctx context.Context, objectName string) error {
	return m.Client.RemoveObject(ctx, m.Bucket, objectName, minio.RemoveObjectOptions{})
}
