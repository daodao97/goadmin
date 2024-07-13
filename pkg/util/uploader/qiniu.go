package uploader

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuUploader struct {
	AccessKey      string
	SecretKey      string
	Bucket         string
	Zone           storage.Zone
	fileNameFormat string
	domain         string
}

func NewQiniuUploader(accessKey string, secretKey string, bucket string, zone storage.Zone, fileNameFormat string, domain string) *QiniuUploader {
	return &QiniuUploader{
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		Bucket:         bucket,
		Zone:           zone,
		fileNameFormat: fileNameFormat,
		domain:         domain,
	}
}

func (qu *QiniuUploader) Upload(fileName string, file *multipart.FileHeader) (location string, err error) {
	// 获取当前时间
	now := time.Now()

	// 解析文件名中的变量
	key := parseFileName(fileName, qu.fileNameFormat, now)
	putPolicy := storage.PutPolicy{
		Scope: qu.Bucket + ":" + key,
	}

	mac := auth.New(qu.AccessKey, qu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &qu.Zone
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	client := http.Client{}
	formUploader := storage.NewFormUploaderEx(&cfg, &storage.Client{Client: &client})
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{},
	}

	reader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file %s", err)
	}
	defer reader.Close()

	err = formUploader.Put(context.Background(), &ret, upToken, key, reader, file.Size, &putExtra)
	if err != nil {
		return "", fmt.Errorf("qiniu upload failed: %v", err)
	}

	return filepath.Join(qu.domain, ret.Key), nil
}
