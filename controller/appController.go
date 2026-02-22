package controller

import (
	"fmt"
	"net/url"
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

	latestFileName, err := findLatestAPKByPackageName(folderPath, packageName)
	if err != nil {
		golog.Error("Error finding latest APK:", err)
		return ResponseError(err)
	}

	url := fmt.Sprintf("http://%s/shared/shared/%s", c.Ctx.Host(), url.PathEscape(latestFileName))

	return ResponseData(url)
}

func findLatestAPKByPackageName(dir, pkgName string) (string, error) {
	var latestPath string
	var latestModTime time.Time

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".apk") {
			continue
		}

		// 去掉后缀，按下划线分割
		base := strings.TrimSuffix(name, ".apk")
		parts := strings.Split(base, "_")
		if len(parts) != 5 {
			continue // 不符合命名规范
		}

		// 检查 package_name 部分是否匹配
		if parts[3] != pkgName {
			continue
		}

		// 获取文件修改时间
		info, err := entry.Info()
		if err != nil {
			// 无法获取信息则跳过该文件
			continue
		}
		modTime := info.ModTime()

		if latestPath == "" || modTime.After(latestModTime) {
			latestPath = name
			latestModTime = modTime
		}
	}

	if latestPath == "" {
		return "", fmt.Errorf("未找到包名为 %s 的 APK 文件", pkgName)
	}
	return latestPath, nil
}
