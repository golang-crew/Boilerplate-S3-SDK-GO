package common

import (
	"io"
	"mime/multipart"
	"os"
	filePath "path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

const (
	s3Region string = viper.GetString(`aws.S3Region`)
	s3Bucket string = viper.GetString(`aws.S3Bucket`)
)

func GenerateDirPath(dirs ...string) string {
	return filePath.Join(dirs...) + "/"
}

/**
* @Param: fileFullPath
* @Param: s3Path
* @Param: fileName
 */
func FileUploadToS3(localFilePath string, s3Path string, fileName string) (result *s3manager.UploadOutput, err error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s3Region),
		})
	if err != nil {
		return
	}

	svc := s3manager.NewUploader(sess)
	result, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Path + fileName),
		Body:   file,
	})

	return result, err
}

func FileSaveToLocal(file *multipart.FileHeader, dirPath string) (err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return
	}

	dst, err := os.Create(dirPath + file.Filename)
	if err != nil {
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return
	}

	return
}
