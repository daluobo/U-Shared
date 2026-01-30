package controller

import (
	"u-shared/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type SharedController struct {
	Ctx     iris.Context
	Service *service.SharedService
}

// GetBy 处理 GET /files/{path}
func (c *SharedController) GetBy(path string) mvc.Result {
	info, err := c.Service.GetFileInfo(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ResponseErrorCode(iris.StatusNotFound, err)
		} else if os.IsPermission(err) {
			return ResponseErrorCode(iris.StatusForbidden, err)
		}
		return ResponseErrorCode(iris.StatusInternalServerError, err)
	}

	if info.IsDir() {
		// 如果是目录，返回目录列表
		list, err := c.Service.ListFiles(path)
		if err != nil {
			return ResponseErrorCode(iris.StatusInternalServerError, err)
		}
		return ResponseData(list)

	} else {
		// 如果是文件，直接提供下载
		c.Ctx.ServeFile(filepath.Join(service.BasePath, path))
		return mvc.Response{}
	}
}

// PostBy 处理 POST /files/{path} 用于上传文件
func (c *SharedController) PostBy(path string) mvc.Result {
	file, _, err := c.Ctx.FormFile("file")
	if err != nil {
		return ResponseErrorCode(iris.StatusBadRequest, err)
	}
	defer file.Close()

	if err := c.Service.UploadFile(path, file, false); err != nil {
		return ResponseError(err)
	}

	return ResponseOk()
}

// PutBy 处理 PUT /files/{path} 用于覆盖上传文件
func (c *SharedController) PutBy(path string) mvc.Result {
	// 从请求体读取文件内容
	body := c.Ctx.Request().Body
	defer body.Close()

	if err := c.Service.UploadFile(path, body, true); err != nil {
		return ResponseError(err)
	}

	return ResponseOk()
}

// DeleteBy 处理 DELETE /files/{path}
func (c *SharedController) DeleteBy(path string) mvc.Result {
	if err := c.Service.DeleteFile(path); err != nil {
		return ResponseError(err)
	}
	return ResponseOk()
}

// HeadBy 处理 HEAD /files/{path}
func (c *SharedController) HeadBy(path string) mvc.Result {
	info, err := c.Service.GetFileInfo(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.Ctx.StatusCode(iris.StatusNotFound)
		} else if os.IsPermission(err) {
			c.Ctx.StatusCode(iris.StatusForbidden)
		} else {
			c.Ctx.StatusCode(iris.StatusInternalServerError)
		}
		return mvc.Response{}
	}

	c.Ctx.Header("Content-Length", strconv.FormatInt(info.Size(), 10))
	c.Ctx.Header("Last-Modified", info.ModTime().Format(http.TimeFormat))
	c.Ctx.StatusCode(iris.StatusOK)
	return mvc.Response{}
}

// Options 处理 OPTIONS /files
func (c *SharedController) Options() mvc.Result {
	c.Ctx.Header("Allow", "GET, HEAD, POST, PUT, DELETE, OPTIONS")
	c.Ctx.Header("DAV", "1")
	c.Ctx.StatusCode(iris.StatusOK)
	return mvc.Response{}
}

func (c *SharedController) BeforeActivation(b mvc.BeforeActivation) {
	// 1. 注册 GET 路由（支持多级路径）
	b.Handle("GET", "/{path:path}", "GetBy")

	// 2. 注册 POST 路由（支持多级路径）
	b.Handle("POST", "/{path:path}", "PostBy")

	// 3. 注册 PUT 路由（支持多级路径）
	b.Handle("PUT", "/{path:path}", "PutBy")

	// 4. 注册 DELETE 路由（支持多级路径）
	b.Handle("DELETE", "/{path:path}", "DeleteBy")

	// 5. 注册 HEAD 路由（支持多级路径）
	b.Handle("HEAD", "/{path:path}", "HeadBy")

	// 6. 注册 OPTIONS 路由（根路径）
	b.Handle("OPTIONS", "/", "Options")
}
