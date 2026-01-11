package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type sharedFile struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type SharedController struct {
	Ctx iris.Context
}

func (c *SharedController) GetBy(folderName string) mvc.Result {
	folderPath := "./shared/" + folderName

	folder, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ResponseError(err)
	}
	defer folder.Close()

	fileList, err := folder.ReadDir(-1)
	if err != nil {
		fmt.Println("Error reading files:", err)
		return ResponseError(err)
	}

	fs := make([]sharedFile, 0)
	for _, file := range fileList {

		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		var t = "folder"
		if !file.IsDir() {
			t = "file"
		}

		info, err := file.Info()
		if err != nil {
			return ResponseError(err)
		}

		fs = append(fs, sharedFile{
			Name: file.Name(),
			Type: t,
			Size: info.Size(),
		})
	}

	return ResponseData(fs)
}

func (c *SharedController) GetDownloadBy(folderName, fileName string) {
	filePath := filepath.Join("./shared/", folderName, fileName)

	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			golog.Errorf("File does not exist:", filePath)
			c.Ctx.StatusCode(iris.StatusNotFound)
			c.Ctx.WriteString("File not found")
			return
		} else {
			golog.Errorf("Error checking file:", err)
			c.Ctx.StatusCode(iris.StatusInternalServerError)
			c.Ctx.WriteString("Internal server error")
			return
		}
	}
	c.Ctx.SendFile(filePath, fileName)
}

func (c *SharedController) PostUpload() mvc.Result {
	folderName := c.Ctx.FormValue("folder")
	fileName := c.Ctx.FormValue("name")
	file, _, err := c.Ctx.FormFile("file")
	if err != nil {
		return ResponseError(err)
	}
	defer file.Close()

	filePath := filepath.Join("./shared/", folderName, fileName)

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			golog.Errorf("File does not exist:", filePath)
		} else {
			golog.Errorf("Error checking file:", err)
			return ResponseError(err)
		}
	} else {
		err = os.Remove(filePath)
		if err != nil {
			golog.Errorf("Error deleting file:", err)
			return ResponseError(err)
		} else {
			golog.Infof("File deleted:", filePath)
		}
	}

	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return ResponseError(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return ResponseError(err)
	}

	return ResponseOk()
}
