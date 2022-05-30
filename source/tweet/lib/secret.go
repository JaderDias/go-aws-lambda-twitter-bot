package tweet

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

func GetObject(ctx context.Context, s3Client *s3.Client, bucketname, key string) ([]byte, error) {
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketname,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, resp.ContentLength)
	defer resp.Body.Close()
	var bbuffer bytes.Buffer
	for {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}
	return bbuffer.Bytes(), nil
}

func PutObject(ctx context.Context, s3Client *s3.Client, bucketname, key string, payload []byte) (*s3.PutObjectOutput, error) {
	reader := bytes.NewReader(payload)
	return s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Body:   reader,
			Bucket: &bucketname,
			Key:    &key,
		})
}
