package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"strconv"
	"time"
)

var (
	bucket    = ""
	accessKey = ""
	secretKey = ""
	domain    = ""
	config    = &storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      true,
		UseCdnDomains: true,
	}
	cacheKeyToken    = []byte("token")
	cacheKeyExpireAt = []byte("expireAt")
)

func NewQny(bk, ak, sk, dm string) {
	bucket = bk
	accessKey = ak
	secretKey = sk
	domain = dm
}

func getUpTokenCache() string {
	now := time.Now()
	token := LocalCache.Get(nil, cacheKeyToken)
	expiredAt := LocalCache.Get(nil, cacheKeyExpireAt)
	if token != nil && expiredAt != nil && now.Unix() < StringToInt64(string(expiredAt)) {
		return string(token)
	}

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	LocalCache.Set(cacheKeyToken, []byte(upToken))
	LocalCache.Set(cacheKeyExpireAt, []byte(strconv.Itoa(int(now.Add(59*time.Minute).Unix()))))

	return upToken
}

func UploadByteData(ctx context.Context, data []byte, key string) (string, error) {
	ret := &storage.PutRet{}
	err := storage.NewFormUploader(config).Put(ctx, ret, getUpTokenCache(), key, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}

	return ret.Key, nil
}

func UploadPartByteData(ctx context.Context, data []byte, key string) (string, error) {
	ret := storage.PutRet{}
	err := storage.NewResumeUploaderV2(config).Put(context.Background(), &ret, getUpTokenCache(), key, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}

	return ret.Key, nil
}

func GenAvatarKey(uid uint) string {
	return fmt.Sprintf("resource_det_search/%d_%d_avater", uid, time.Now().Unix())
}

func GenDocKey(docId uint, uid uint) string {
	return fmt.Sprintf("resource_det_search/doc/%d_%d_%d", uid, docId, time.Now().Unix())
}

func GenFileLink(key string) string {
	return domain + key
}
