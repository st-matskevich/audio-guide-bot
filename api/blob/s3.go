package blob

import (
	"context"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3BlobProvider struct {
	client     *minio.Client
	bucketName string
}

func (provider *S3BlobProvider) ReadBlob(name string, writer io.Writer) error {
	object, err := provider.client.GetObject(context.Background(), provider.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}

	return nil
}

func (provider *S3BlobProvider) WriteBlob(name string, reader io.Reader) error {
	_, err := provider.client.PutObject(context.Background(), provider.bucketName, name, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateS3BlobProvider(URL string) (BlobProvider, error) {
	urlObject, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	host := urlObject.Host
	accessID := urlObject.User.Username()
	accessSecret, _ := urlObject.User.Password()
	bucket := strings.Trim(urlObject.Path, "/")

	params := urlObject.Query()
	useSSL, _ := strconv.ParseBool(params.Get("ssl"))

	s3Client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessID, accessSecret, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	provider := S3BlobProvider{
		client:     s3Client,
		bucketName: bucket,
	}

	return &provider, nil
}
