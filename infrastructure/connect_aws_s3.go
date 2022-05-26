package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectAws(env Env, logger Logger) *session.Session {
	AccessKeyID := env.AwsAccessKey
	SecretAccessKey := env.AwsAccessKey
	MyRegion := env.AwsRegion
	logger.Zap.Info(".... config........ ", AccessKeyID, SecretAccessKey, MyRegion)
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				"AKIAXLFZ3XXC56MBUFXN",
				"HwhfF7z3CfLKbxGzOsNkWclieN7p1xEKCRMWRo9v",
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {

		logger.Zap.Panic(err)
	}
	logger.Zap.Info(".... AWS is initialized........ 🎉🎉🎉 ")
	return sess
}
