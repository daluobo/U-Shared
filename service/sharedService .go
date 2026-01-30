package service

import (
	"fmt"
	"u-shared/model"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const BasePath = "./shared/"

type SharedService struct{}

func NewSharedService() *SharedService {
	return &SharedService{}
}

func (s *SharedService) ListFiles(path string) (*[]model.SharedFile, error) {
	fullPath := filepath.Join(BasePath, path)
	if !s.isSafePath(fullPath) {
		return nil, os.ErrPermission
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	SharedFiles := make([]model.SharedFile, 0, len(files))
	for _, file := range files {
		info, _ := file.Info()
		SharedFiles = append(SharedFiles, model.SharedFile{
			Name:    file.Name(),
			Path:    strings.ReplaceAll(filepath.Join(path, file.Name()), "\\", "/"),
			Size:    info.Size(),
			ModTime: info.ModTime().UnixNano() / int64(time.Millisecond),
			IsDir:   file.IsDir(),
		})
	}

	return &SharedFiles, nil
}

func (s *SharedService) DownloadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(BasePath, path)
	if !s.isSafePath(fullPath) {
		return nil, os.ErrPermission
	}
	return os.ReadFile(fullPath)
}

func (s *SharedService) UploadFile(path string, file io.Reader, overwrite bool) error {
	fullPath := filepath.Join(BasePath, path)

	// 安全检查
	if !s.isSafePath(fullPath) {
		return os.ErrPermission
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err == nil {
		// 文件存在
		if !overwrite {
			return fmt.Errorf("文件已存在: %s", path)
		}
		// 如果允许覆盖，继续执行
	} else if !os.IsNotExist(err) {
		// 其他错误
		return err
	}

	// 创建或覆盖文件
	out, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

func (s *SharedService) DeleteFile(path string) error {
	fullPath := filepath.Join(BasePath, path)
	if !s.isSafePath(fullPath) {
		return os.ErrPermission
	}
	return os.RemoveAll(fullPath)
}

func (s *SharedService) GetFileInfo(path string) (os.FileInfo, error) {
	fullPath := filepath.Join(BasePath, path)
	if !s.isSafePath(fullPath) {
		return nil, os.ErrPermission
	}
	return os.Stat(fullPath)
}

func (s *SharedService) isSafePath(path string) bool {
	absBase, _ := filepath.Abs(BasePath)
	absPath, _ := filepath.Abs(path)
	return strings.HasPrefix(absPath, absBase)
}
