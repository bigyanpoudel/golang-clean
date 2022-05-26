package service

import (
	"go-clean-api/infrastructure"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSService struct {
	session *session.Session
	logger  infrastructure.Logger
	env     infrastructure.Env
}

func NewAWSService(session *session.Session,
	logger infrastructure.Logger,
	env infrastructure.Env) AWSService {
	return AWSService{
		session: session,
		logger:  logger,
		env:     env,
	}
}

func (w AWSService) Upload(bucketName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// size := fileHeader.Size
	// buffer := make([]byte, size)
	// file.Read(buffer)
	// tempFileName := "https://" + bucketName + "." + "s3-ap-south-1" + ".amazonaws.com/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	// _, err := s3.New(w.session).PutObject(&s3.PutObjectInput{
	// 	Bucket:               aws.String(bucketName),
	// 	Key:                  aws.String(tempFileName),
	// 	ACL:                  aws.String("public-read"),
	// 	Body:                 bytes.NewReader(buffer),
	// 	ContentLength:        aws.Int64(int64(size)),
	// 	ContentType:          aws.String(http.DetectContentType(buffer)),
	// 	ContentDisposition:   aws.String("attachment"),
	// 	ServerSideEncryption: aws.String("AES256"),
	// 	StorageClass:         aws.String("INTELLIGENT_TIERING"),
	// })
	// if err != nil {
	// 	return "", err
	// }

	// // up, err := uploader.Upload(&s3manager.UploadInput{
	// // 	Bucket: aws.String(w.env.AwsBucketName),
	// // 	ACL:    aws.String("public-read"),
	// // 	Key:    aws.String(fileName),
	// // 	Body:   &bytes.Reader{},
	// // })
	// w.logger.Zap.Info("file", tempFileName)
	// return tempFileName, nil
	uploader := s3manager.NewUploader(w.session)
	filename := fileHeader.Filename
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	w.logger.Zap.Info("upload.....", up, err)
	if err != nil {
		return "", err
	}

	filepath := "https://" + bucketName + "." + "s3.amazonaws.com/" + filename
	w.logger.Zap.Info("file", filepath)
	return filepath, nil

}
