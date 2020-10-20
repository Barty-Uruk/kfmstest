package service

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	CreateFileLink(folderName, filename string) string
	UploadFile(ctx context.Context, req UploadRequest) (UploadResponse, error)
	DownloadFile(ctx context.Context, req DownloadRequest) (DownloadResponse, error)
}
type amazonS3 struct {
	Logger         *log.Logger
	Downloader     *s3manager.Downloader
	Host           string
	S3             *s3.S3
	Session        *session.Session
	Bucket         string
	RootFolderName string
}

func (s *amazonS3) CreateFileLink(folderName, filename string) string {
	s.RootFolderName = "storage"
	return s.RootFolderName + "/" + folderName + "/" + filename
}
func NewService(log *log.Logger, c *configs.AmazonS3) (Service, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(endpoints.UsEast2RegionID),
		Endpoint: aws.String(c.Address),
	}))
	svc := s3.New(sess, newAWSConfig(c))
	return &amazonS3{
		Logger:  log,
		Session: sess,
		Downloader: s3manager.NewDownloader(session.Must(session.NewSession(&aws.Config{
			Region:   aws.String(endpoints.UsEast2RegionID),
			Endpoint: aws.String(c.Address),
		}))),
		RootFolderName: c.RootFolderName,
		Bucket:         c.Bucket,
		S3:             svc,
	}, nil
}
func newAWSConfig(c *configs.AmazonS3) *aws.Config {
	conf := aws.NewConfig()
	conf.Endpoint = aws.String(c.Address)
	conf.Region = aws.String(endpoints.UsEast2RegionID)
	return conf
}

func (s *amazonS3) DownloadFile(ctx context.Context, req DownloadRequest) (DownloadResponse, error) {
	buff := aws.NewWriteAtBuffer([]byte{})

	_, err := s.Downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(req.FileName),
		})
	if err != nil {
		return DownloadResponse{}, err
	}
	bf := bytes.NewBuffer(buff.Bytes())
	return DownloadResponse{File: bf}, nil
}
func (s *amazonS3) UploadFile(ctx context.Context, req UploadRequest) (UploadResponse, error) {

	fileLink := s.CreateFileLink(req.FolderName, req.FileName)

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
	DownloadFile endpoint.Endpoint
	UploadFile   endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		UploadFile:   makeUploadFileEndpoint(s),
		DownloadFile: makeDownloadFileEndpoint(s),
	}
}

func makeUploadFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadRequest)
		res, err := s.UploadFile(ctx, req)
		return res, err
	}
}

func makeDownloadFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DownloadRequest)
		res, err := s.DownloadFile(ctx, req)
		return res, err
	}
}
