package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-crew/Bolierplate-CRUD-Gingonic/common"
	"github.com/spf13/viper"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func uploadFileToS3(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No request file"})
		return
	}

	dirPath := common.GenerateDirPath("temp", "firstSubDir", "secondSubDir")
	s3Path := common.GenerateDirPath("s3Path", "firstSubdir", "secondSubdir")
	localFilePathAndName := dirPath + file.Filename

	err = common.FileSaveToLocal(file, dirPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot save file to local"})
		return
	}

	res, err := common.FileUploadToS3(localFilePathAndName, s3Path, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail to upload s3"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"S3URL": res.Location})
}

func viperInit() {
	viper.SetConfigFile(`./config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()

	viperInit()

	r.GET("/ping", ping)
	r.POST("/", uploadFileToS3)

	r.Run(":8080")
}
