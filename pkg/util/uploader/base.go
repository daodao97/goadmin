package uploader

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type Uploader interface {
	Upload(fileName string, file *multipart.FileHeader) (location string, err error)
}

// parseFileName 解析文件名中的变量并替换为实际的值
// 支持 {year} {month} {day} {hour} {minute} {second} {since_second} {since_millisecond} {random} {filename} {.suffix} {suffix} {mimetype} 等变量。
// 比如: 上传的图片为 uPic.jpg, 设定为 “uPic/{filename}{.suffix}”, 则会保存到 “uPic/uPic.jpg”
func parseFileName(fileName string, fileNameFormat string, now time.Time) string {
	if fileNameFormat == "" {
		return fileName
	}
	replacer := strings.NewReplacer(
		"{year}", fmt.Sprintf("%04d", now.Year()),
		"{month}", fmt.Sprintf("%02d", now.Month()),
		"{day}", fmt.Sprintf("%02d", now.Day()),
		"{hour}", fmt.Sprintf("%02d", now.Hour()),
		"{minute}", fmt.Sprintf("%02d", now.Minute()),
		"{second}", fmt.Sprintf("%02d", now.Second()),
		"{since_second}", fmt.Sprintf("%d", now.Unix()),
		"{since_millisecond}", fmt.Sprintf("%d", now.UnixNano()/int64(time.Millisecond)),
		"{random}", fmt.Sprintf("%d", now.Nanosecond()),
		"{filename}", strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		"{.suffix}", filepath.Ext(fileName),
		"{suffix}", strings.TrimPrefix(filepath.Ext(fileName), "."),
		"{mimetype}", getMimeType(fileName),
	)

	return replacer.Replace(fileNameFormat)
}

// getMimeType 获取文件的 MIME 类型
func getMimeType(fileName string) string {
	// 获取文件后缀名
	ext := filepath.Ext(fileName)

	// 使用 mime.TypeByExtension 获取后缀名对应的 MIME 类型
	mimeType := mime.TypeByExtension(ext)

	// 如果未找到对应的 MIME 类型，返回 "application/octet-stream" 作为默认值
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType
}
