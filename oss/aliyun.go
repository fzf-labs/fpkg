package oss

import (
	"io"
	"strings"

	aliOssSdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

// AliConfig oss 上传配置
type AliConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Host      string
	Prefix    string
	Salt      string
}

// CreateBucket oss 根据参数来创建 Bucket
func (c *AliConfig) CreateBucket() (bucket *aliOssSdk.Bucket, err error) {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := c.Endpoint
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyID := c.AccessKey
	accessKeySecret := c.SecretKey
	bucketName := c.Bucket
	// 创建OSSClient实例。
	ossClient, err := aliOssSdk.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, errors.Wrapf(err, "创建 aliyun OSSClient实例失败")
	}
	// 获取存储空间。
	bucket, err = ossClient.Bucket(bucketName)
	if err != nil {
		return nil, errors.Wrapf(err, "获取 aliyun OSS 存储空间失败")
	}
	return
}

// FileURL 获取文件的oss中的url
func (c *AliConfig) FileURL(path string, isEncrypt int32, wh string) string {
	url := strings.Trim(c.Host, "/") + "/" + strings.Trim(path, "/")
	if isEncrypt == 1 {
		url += "?encrypt=1"
	}
	if wh != "" {
		url += "&size=" + wh
	}
	return url
}

// SpliceURL 获取文件的oss中的url
func (c *AliConfig) SpliceURL(path string) string {
	if path == "" {
		return ""
	}
	return c.Host + "/" + strings.TrimPrefix(path, "/")
}

// oss 上传客户端
type aliOss struct {
	bucket *aliOssSdk.Bucket
}

func NewAliOss(c *AliConfig) Driver {
	bucket, err := c.CreateBucket()
	if err != nil {
		panic(err)
	}
	return &aliOss{
		bucket: bucket,
	}
}

// Put 上传
func (c *aliOss) Put(objectName, localFileName string) error {
	err := c.bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return errors.Wrapf(err, "put oss file fail")
	}
	return nil
}

func (c *aliOss) PutObj(objectName string, file io.Reader) error {
	err := c.bucket.PutObject(objectName, file)
	if err != nil {
		return errors.Wrapf(err, "put oss file fail")
	}
	return nil
}

// Get 下载
func (c *aliOss) Get(objectName, downloadedFileName string) error {
	err := c.bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		return errors.Wrapf(err, "get oss file fail")
	}
	return nil
}

// Del 删除
func (c *aliOss) Del(objectName string) error {
	// 删除文件。
	err := c.bucket.DeleteObject(objectName)
	if err != nil {
		return errors.Wrapf(err, "del oss file fail")
	}
	return nil
}
