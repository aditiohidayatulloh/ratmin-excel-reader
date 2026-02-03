package excelreader

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newMinioClient(cfg MinioConfig) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
}

func getObject(
	ctx context.Context,
	cfg MinioConfig,
	bucket string,
	object string,
) (io.ReadCloser, error) {

	client, err := newMinioClient(cfg)
	if err != nil {
		return nil, err
	}

	obj, err := client.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
