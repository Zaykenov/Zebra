package utils

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateItemImage(c *gin.Context) (string, error) {
	// Image handling
	c.Request.ParseMultipartForm(10 << 20)

	fileName := ""
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := c.Request.FormFile("image")
	if file != nil {
		defer file.Close()
	}

	if file != nil {
		defer file.Close()
	}

	locationImage, exists := os.LookupEnv("LocationItemDocker")

	if !exists {
		return "", errors.New("enviroment variable is not set")
	}

	if err != nil {
		if err.Error() != "http: no such file" {

			return "", err
		}
	} else {
		// Create a temporary file within our temp-images directory that follows
		tempFile, err := ioutil.TempFile(locationImage, "upload-*.jpeg")
		if err != nil {
			return "", err
		}
		defer tempFile.Close()
		// read all of the contents of our uploaded file into a
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}
		tempFile.Write(fileBytes)

		fileName = filepath.Base(tempFile.Name())
	}
	c.Request.ParseMultipartForm(0)

	return fileName, nil
}

func GenCode() string {
	str := ""
	for i := 0; i < 4; i++ {
		str = str + strconv.Itoa(rand.Intn(10))
	}
	return str
}

func CalculatePageCount(totalCount int64, pageSize int64) int64 {
	if totalCount%pageSize == 0 {
		return totalCount / pageSize
	}
	return totalCount/pageSize + 1
}
