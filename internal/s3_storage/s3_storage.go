package s3storage

import (
	"bytes"
	"context"
	"log"
	"os"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ikonera/codex/internal/config"
)

type IS3Manager interface{}

type S3Manager struct {
	Client        *s3.Client
	ConfigManager config.IConfigManager
}

func NewS3Manager() *S3Manager {
	cfgManager := config.NewYAMLConfigManager()
	userConfig, err := cfgManager.ReadConfig()
	if err != nil {
		log.Fatalf("Can't retrieve Codex configuration: %s\n", err.Error())
	}
	clientConf, err := awsCfg.LoadDefaultConfig(
		context.TODO(),
		awsCfg.WithCredentialsProvider(
			awsCreds.NewStaticCredentialsProvider(userConfig.Credentials.AccessKey, userConfig.Credentials.SecretKey, ""),
		),
	)
	if err != nil {
		log.Fatalf("Couldn't load aws client config: %s\n", err.Error())
	}

	return &S3Manager{
		Client: s3.NewFromConfig(
			clientConf,
		),
		ConfigManager: cfgManager,
	}
}

func (m *S3Manager) CheckForBucket(bucketName string) (bool, error) {
	_, err := m.Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *S3Manager) Upload(path, bucket string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	s3UploadInput := &s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &path,
	}

	resp, err := m.Client.CreateMultipartUpload(context.TODO(), s3UploadInput)
	if err != nil {
		return err
	}

	uploadId := resp.UploadId
	var parts []types.CompletedPart
	bufSize := 5 * 1024 * 1024
	partNb := int32(1)

	for {
		buffer := make([]byte, bufSize)
		n, err := file.Read(buffer)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		s3UploadPartInput := &s3.UploadPartInput{
			Bucket:     &bucket,
			Key:        &path,
			PartNumber: &partNb,
			UploadId:   uploadId,
			Body:       bytes.NewReader(buffer[:n]),
		}

		uploadRes, err := m.Client.UploadPart(context.TODO(), s3UploadPartInput)
		if err != nil {
			return err
		}

		parts = append(parts, types.CompletedPart{
			ETag:       uploadRes.ETag,
			PartNumber: &partNb,
		})
		partNb++
	}

	s3CompletedMultipartInput := &s3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &path,
		UploadId: uploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	}

	if _, err := m.Client.CompleteMultipartUpload(context.TODO(), s3CompletedMultipartInput); err != nil {
		return err
	}

	return nil
}
