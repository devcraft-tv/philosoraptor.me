package uploader

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

const dataType string = "image/jpeg"
const fileType string = "jpg"
const acl s3.ACL = s3.PublicRead
const s3RootUrl = "https://s3.amazonaws.com"

type S3Uploader struct {
	AccessId        string
	SecretAccessKey string
	Bucket          string
	Path            string
}

func (uploader S3Uploader) Upload(data []byte, fileName string) (url string, err error) {
	awsAuth, err := aws.GetAuth(uploader.AccessId, uploader.SecretAccessKey)
	s3Client := s3.New(awsAuth, aws.USEast)
	bucket := s3.Bucket{s3Client, uploader.Bucket}
	filePath := uploader.Path + fileName + "." + fileType

	bucket.Put(filePath, data, dataType, acl)
	url = s3RootUrl + "/" + uploader.Bucket + "/" + filePath

	return
}
