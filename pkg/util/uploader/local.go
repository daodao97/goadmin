package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// LocalUploader 实现了 Uploader 接口，用于保存上传文件到本地
type LocalUploader struct {
	uploadPath     string
	fileNameFormat string
	domain         string
}

// NewLocalUploader 创建一个 LocalUploader 实例
func NewLocalUploader(uploadPath, fileNameFormat, domain string) *LocalUploader {
	return &LocalUploader{uploadPath: uploadPath, fileNameFormat: fileNameFormat, domain: domain}
}

// Upload 将文件保存到本地目录中，并返回保存后的文件路径
func (u *LocalUploader) Upload(fileName string, file *multipart.FileHeader) (location string, err error) {
	if err := os.MkdirAll(u.uploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("无法创建上传目录：%v", err)
	}

	// 获取当前时间
	now := time.Now()

	// 解析文件名中的变量
	saveFileName := parseFileName(fileName, u.fileNameFormat, now)

	// 构造保存文件的路径
	filePath := filepath.Join(u.uploadPath, saveFileName)

	// 检查目标目录是否存在，如果不存在，则创建
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("无法创建目标目录：%v", err)
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("保存文件失败：%v", err)
	}

	return filepath.Join(u.domain, filePath), nil
}
