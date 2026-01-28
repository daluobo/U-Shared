package middleware

import (
	"strings"

	"github.com/kataras/iris/v12"
)

var MediaFolders = []string{
	"Audiobooks",
	"Podcasts",
	"Downloads",
	"Pictures",
	"Music",
	"Movies",
	"Documents",
	"Shared",
}

func PathValidator(ctx iris.Context) {
	// 从URL参数获取路径
	path := ctx.Params().Get("path")

	// 如果没有path参数，跳过验证（可能是根目录）
	if path == "" || path == "." || path == "/" {
		ctx.Next()
		return
	}

	// 清理路径，移除开头的斜杠
	cleanPath := strings.TrimPrefix(path, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "./")

	// 分割路径获取第一个目录
	parts := strings.Split(cleanPath, "/")
	if len(parts) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "无效的路径",
			"code":  iris.StatusBadRequest,
		})
		return
	}

	firstDir := parts[0]

	// 检查是否以允许的前缀开头
	isValid := false
	for _, prefix := range MediaFolders {
		if firstDir == prefix {
			isValid = true
			break
		}
	}

	if !isValid {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"error":   "路径必须以下列目录之一开头",
			"allowed": MediaFolders,
			"path":    path,
			"code":    iris.StatusForbidden,
		})
		return
	}

	ctx.Next()
}