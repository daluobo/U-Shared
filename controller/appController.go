package controller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AppController struct {
	Ctx iris.Context
}

func (c *AppController) Get() mvc.Result {
	packageName := c.Ctx.URLParam("package_name")

	folderPath := "./shared/shared"
	folder, err := os.Open(folderPath)
	if err != nil {
		golog.Error("Error reading directory:", err)
		return ResponseError(err)
	}
	defer folder.Close()

	fileList, err := folder.ReadDir(-1)
	if err != nil {
		golog.Error("Error reading files:", err)
		return ResponseError(err)
	}

	var latestModTime time.Time
	var latestFileInfo os.FileInfo
	found := false

	for _, file := range fileList {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// 检查是否以.apk结尾
		if !strings.HasSuffix(fileName, ".apk") {
			continue
		}

		// 检查是否以packageName开头
		if strings.HasPrefix(fileName, packageName) {
			// 获取去除包名后的剩余部分
			remaining := strings.TrimPrefix(fileName, packageName)

			// 关键判断：排除包含额外包名段的情况
			// 我们想要的是：包名后直接跟 _ 或其他非点字符
			// 排除包名后直接跟 . 的情况
			if len(remaining) > 0 {
				// 如果剩余部分以点开头，说明包名后面有额外的包段（如.tv_lite），跳过
				if remaining[0] == '.' {
					continue
				}

				// 或者使用更严格的检查：剩余部分必须以 _ 开头
				// 这样可以确保匹配 com.daluobo.litplayer_xxx.apk 格式
				// 如果不满足这个条件也跳过
				if remaining[0] != '_' {
					// 这里可以根据需要调整，如果想要更严格的匹配
					// 可以只接受下划线开头的
					// continue
				}
			}

			// 获取文件信息
			fileInfo, err := file.Info()
			if err != nil {
				golog.Warnf("Cannot get info for file %s: %v", fileName, err)
				continue
			}

			// 检查是否比当前最新文件更新
			if !found || fileInfo.ModTime().After(latestModTime) {
				latestModTime = fileInfo.ModTime()
				latestFileInfo = fileInfo
				found = true
			}
		}
	}

	if !found {
		return ResponseFail(fmt.Errorf("no matching APK files found for package '%s'", packageName))
	}

	return ResponseData(latestFileInfo.Name())
}
