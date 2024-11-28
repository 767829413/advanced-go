package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetObject(ctx context.Context, bucket string, key string) ([]byte, error)

	GetFile(ctx context.Context, bucket string, key string, localFile string) error

	PutObject(
		ctx context.Context,
		bucket string,
		key string,
		data []byte,
		metadata map[string]string,
	) error

	PutObjectWithMeta(
		ctx context.Context,
		bucket string,
		key string,
		data []byte,
		metadata *Metadata,
	) error

	PutFile(
		ctx context.Context,
		bucket string,
		key string,
		srcFile string,
		metadata map[string]string,
	) error

	PutObjectFromFile(
		ctx context.Context,
		bucket string,
		key string,
		filePath string,
		metadata map[string]string,
	) error

	PutFileWithMeta(
		ctx context.Context,
		bucket string,
		key string,
		srcFile string,
		metadata *Metadata,
	) error

	// 大文件 断点续传  分片数量不能超过10000,分片大小推荐  100K-1G partSize单位KB
	PutFileWithPart(
		ctx context.Context,
		bucket string,
		key string,
		srcFile string,
		metadata *Metadata,
		partSize int64,
	) error

	ListObjects(ctx context.Context, bucket string, prefix string) ([]Content, error)

	DeleteObject(ctx context.Context, bucket string, key string) error

	CopyObject(ctx context.Context, bucket string, srcKey string, destKey string) error

	SetObjectAcl(ctx context.Context, bucket string, key string, acl StorageAcl) error

	SetObjectMetaData(ctx context.Context, bucket string, key string, metadata *Metadata) error

	IsObjectExist(ctx context.Context, bucket string, key string) (bool, error)

	GetDirToken(ctx context.Context, remoteDir string) (map[string]interface{}, error)

	GetDirToken2(ctx context.Context, remoteDir string) (*StorageToken, error)

	GetDirTokenWithAction(
		ctx context.Context,
		remoteDir string,
		actions ...Action,
	) (bool, *StorageToken)
	//过期时间：秒
	SignFile(ctx context.Context, remoteDir string, expiredTime int64) (string, error)

	//过期时间：秒
	SignFile2(ctx context.Context, bucket, remoteDir string, expiredTime int64) (string, error)

	//过期时间：秒
	SignFileForDownload(
		ctx context.Context,
		remoteDir string,
		expiredTime int64,
		downloadName string,
	) (string, error)

	GetObjectMeta(ctx context.Context, bucket string, key string) (*Content, error)
	//解冻归档文件，成功就为true,异常或者失败返回false
	RestoreArchive(ctx context.Context, bucket string, key string) (bool, error)

	//判断是否归档文件
	IsArchive(ctx context.Context, bucket string, key string) (bool, error)

	//批量删除文件；可以减少调用次数，进而减少费用
	BatchDeleteObject(
		ctx context.Context,
		bucketName string,
		list []string,
	) (successList []string, e error)
}

// listObject结果对象
type Content struct {
	Key          string
	Size         int64
	ETag         string
	LastModified time.Time
}

type StorageAcl string

const (
	AclPrivate         StorageAcl = "private"
	AclPublicRead      StorageAcl = "public-read"
	AclPublicReadWrite StorageAcl = "public-read-write" //2021.1.8当前测试华为public-read-write设置未生效
	AclDefault         StorageAcl = "default"           //华为不支持，请勿使用
)

type Metadata struct {
	Mime               string
	ContentEncoding    string
	Acl                string
	ContentDisposition string
}

type StorageToken struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Bucket          string `json:"bucket"`
	Expire          int64  `json:"expire"`
	Host            string `json:"host"`
	Provider        string `json:"provider"`
	Region          string `json:"region"`
	StsToken        string `json:"stsToken"`
	UploadPath      string `json:"uploadPath"`
	Path            string `json:"path"`
	CdnDomain       string `json:"cdnDomain"`
}

type Action string

const PutObjectAction Action = "oss:PutObject"
const PutObjectAclAction Action = "oss:PutObjectAcl"
const GetObjectAction Action = "oss:GetObject"

type StorageConfig struct {
	Provider         string   `json:"provider"` //ali/huawei
	AccessKeyId      string   `json:"accessKeyId"`
	AccessKeySecret  string   `json:"accessKeySecret"`
	Endpoint         string   `json:"endpoint"`
	EndpointInternal string   `json:"endPointInternal"`
	StsEndpoint      string   `json:"stsEndPoint"`
	Bucket           string   `json:"bucket"`
	RoleArn          string   `json:"roleArn"`
	Region           string   `json:"region"`
	Root             string   `json:"root"`
	TmpRoot          string   `json:"tmpRoot"`
	Internal         bool     `json:"internal"`
	Host             string   `json:"host"`
	CdnDomain        string   `json:"cdnDomain"`
	CdnProtocol      string   `json:"cdnProtocol"`
	Path             []string `json:"path"`
	TmpPath          []string `json:"tmpPath"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
}
