package goutil

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

//文件上传
//    bucket:七牛云存储空间名
//    pathType：文件存储路径的名称，例如：photo,video
//    hostName：七牛云的自定义域名
//    accessKey:
//    secretKey:
//    formFile：接收到的文件
//    limitSize：文件大小限制（单位Byte）
//    path：文件路径
//    err：error
func UploadFileToQiniuCloud(bucket, pathType, hostName, accessKey, secretKey string, formFile *multipart.FileHeader, limitSize int64) (urlPath string, err error) {
	fileSize := formFile.Size
	//限制图片上传大小
	if fileSize > limitSize {
		return "", errors.New("文件超过大小限制")
	}
	//打开上传的源文件
	srcFile, err := formFile.Open()
	defer srcFile.Close()
	if err != nil {
		return "", errors.New("文件打开出错:" + err.Error())
	}

	//创建要把保存的空文件
	fileName := formFile.Filename
	split := strings.Split(fileName, ".")
	fileType := split[len(split)-1]
	fileKey := pathType + "/" + GetRandomString(32) + "." + fileType

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	config := new(storage.Config)
	// 空间对应的机房
	config.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	config.UseHTTPS = false
	// 上传是否使用CDN上传加速
	config.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(config)

	ret := new(storage.PutRet)
	putExtra := new(storage.PutExtra)
	err = formUploader.Put(context.Background(), ret, upToken, fileKey, srcFile, fileSize, putExtra)
	if err != nil {
		return "", errors.New("文件远程保存出错：" + err.Error())
	}
	urlPath = hostName + ret.Key
	return
}
