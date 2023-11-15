package utils

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

const (
	ACCESS_KEY = "0vjg88Ci21r-2GZllTv8vzdKoSxwBuvK469Vr298"
	SECRET_KEY = "2STD28f6zGVvk46xNy6UDi_aacUd_wtS3Mxys49C"
	TARGE_TURL = "rpd659d3i.hn-bkt.clounddn.com"
	BUCKET     = "douyinhaha"
)

func GeneratePlayURL(key string) string {
	return "http" + TARGE_TURL + "/" + key
}

// UpLoadFile 该函数用于上传文件到七牛云
func UpLoadFile(title string, video multipart.FileHeader, videoSize int64) (string, error) {

	// 获取上传凭证uploadToken
	mac := qbox.NewMac(ACCESS_KEY, SECRET_KEY)
	putPolicy := storage.PutPolicy{Scope: BUCKET}
	putPolicy.Expires = 600
	uploadToken := putPolicy.UploadToken(mac)
	File, err := video.Open()
	if err != nil {
		return "", nil
	}

	// 设置UpLoadFile的配置文件
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuanan,
	}
	// new一个putExtra对象
	putExtra := storage.PutExtra{}
	// 获取文件上传加载器
	formUploader := storage.NewFormUploader(&cfg)
	// 创建ret
	ret := storage.PutRet{}

	err = formUploader.Put(context.Background(), ret, uploadToken, title, File, videoSize, &putExtra)
	if err != nil {
		return "", err
	}

	return GeneratePlayURL(title), nil
}
