package service

import (
	"context"
	"fmt"

	"github.com/Barty-Uruk/kfmstest/configs"
	log "github.com/Barty-Uruk/kfmstest/pkg/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	SayHello(ctx context.Context, name string) (string, error)
	CreateFileLink(folderName, filename string) string
	UploadFile(ctx context.Context, req UploadRequest) (UploadResponse, error)
}
type amazonS3 struct {
	Logger         *log.Logger
	Host           string
	S3             *s3.S3
	Bucket         string
	RootFolderName string
}

func (s *amazonS3) SayHello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("%s, %s", s.Host, name), nil
}
func (s *amazonS3) CreateFileLink(folderName, filename string) string {
	s.RootFolderName = "storage"
	return s.RootFolderName + "/" + folderName + "/" + filename
}
func NewService(log *log.Logger, c *configs.AmazonS3) (Service, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast2RegionID),
	}))
	svc := s3.New(sess, newAWSConfig(c.Address))
	return &amazonS3{
		Logger:         log,
		RootFolderName: c.RootFolderName,
		Bucket:         c.Bucket,
		S3:             svc,
	}, nil
}
func newAWSConfig(endpoint string) *aws.Config {
	conf := aws.NewConfig()
	conf.Endpoint = &endpoint
	return conf
}
func (s *amazonS3) UploadFile(ctx context.Context, req UploadRequest) (UploadResponse, error) {

	fileLink := s.CreateFileLink(req.FolderName, req.FileName)
	_ = "https://ams3.digitaloceanspaces.com"

	// f, err := os.Open("image_2020-06-29_12-14-31.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// cred, _ := sess.Config.Credentials.Get()
	// fmt.Println(cred.AccessKeyID, cred.SecretAccessKey)
	_, err := s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket), // Bucket to be used
		Key:    aws.String(fileLink), // Name of the file to be saved
		Body:   req.File,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			return UploadResponse{}, fmt.Errorf("upload canceled due to timeout, %v\n", err)
		} else {
			return UploadResponse{}, fmt.Errorf("failed to upload object, %v\n", err)
		}
	}
	return UploadResponse{FileLink: fileLink}, nil
}

type HelloRequest struct {
	Name string `json:"name"`
}
type HelloResponse struct {
	HelloText string `json:"hello_text"`
}

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Hello      endpoint.Endpoint
	UploadFile endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Hello: makeHelloEndpoint(s),
	}
}
func makeHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HelloRequest)
		text, err := s.SayHello(ctx, req.Name)
		return HelloResponse{HelloText: text}, err
	}
}

func makeUploadFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadRequest)
		res, err := s.UploadFile(ctx, req)
		return res, err
	}
}
