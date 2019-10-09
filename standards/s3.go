package MagicCapKernelStandards

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"magiccap-uploaders-kernel/utils"
	"net/http"
)

// S3Spec defines the S3 spec for this uploader.
type S3Spec struct {
	AccessKeyID string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Endpoint string `json:"endpoint"`
	BucketName string `json:"bucket_name"`
	BucketURL string `json:"bucket_url"`
}

// S3Init defines the S3 standard.
func S3Init(Structure UploaderStructure) (*Uploader, error) {
	b, err := json.Marshal(Structure.Spec)
	if err != nil {
		return nil, err
	}
	var spec S3Spec
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	return &Uploader{
		Description:   Structure.Description,
		Name:          Structure.Name,
		ConfigOptions: Structure.Config,
		Icon:          Structure.Icon,
		Upload: func(Config map[string]interface{}, Data []byte, Filename string) (string, error) {
			AccessKeyID, err := utils.SubString(spec.AccessKeyID, Config, Filename)
			if err != nil {
				return "", err
			}
			SecretAccessKey, err := utils.SubString(spec.SecretAccessKey, Config, Filename)
			if err != nil {
				return "", err
			}
			Endpoint, err := utils.SubString(spec.Endpoint, Config, Filename)
			if err != nil {
				return "", err
			}
			BucketName, err := utils.SubString(spec.BucketName, Config, Filename)
			if err != nil {
				return "", err
			}
			BucketURL, err := utils.SubString(spec.BucketURL, Config, Filename)
			if err != nil {
				return "", err
			}

			StaticCredential := credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, "")
			s3sess := session.Must(session.NewSession(&aws.Config{
				Endpoint: &Endpoint,
				Credentials: StaticCredential,
			}))
			svc := s3.New(s3sess)
			h := http.DetectContentType(Data)
			UploadParams := &s3.PutObjectInput{
				Bucket: &BucketName,
				Key: &Filename,
				ContentType: &h,
				Body: bytes.NewReader(Data),
				ACL: aws.String("public-read"),
				ContentLength: aws.Int64(int64(len(Data))),
			}
			_, err = svc.PutObject(UploadParams)
			if err != nil {
				return "", err
			}

			return BucketURL + Filename, nil
		},
	}, nil
}
