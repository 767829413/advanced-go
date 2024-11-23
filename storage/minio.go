package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/767829413/advanced-go/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	MAX_MINIO_BATCH_DELETE_COUNT = 1000
	ProviderMinio                = "minio"
)

type minioStorage struct {
	client *minio.Client
	config *StorageConfig
}

func newMinioStorager(c *StorageConfig) (*minioStorage, error) {
	endpoint := util.If(c.Internal, c.EndpointInternal, c.Endpoint)
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			c.AccessKeyId,
			c.AccessKeySecret,
			"",
		),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &minioStorage{config: c, client: client}, nil
}

func (m *minioStorage) GetObject(bucket string, key string) ([]byte, error) {
	obj, err := m.client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	return buf.Bytes(), err
}

func (m *minioStorage) GetFile(bucket string, key string, localFile string) error {
	return m.client.FGetObject(
		context.Background(),
		bucket,
		key,
		localFile,
		minio.GetObjectOptions{},
	)
}

func (m *minioStorage) PutObject(
	bucket string,
	key string,
	data []byte,
	metadata map[string]string,
) error {
	_, err := m.client.PutObject(
		context.Background(),
		bucket,
		key,
		bytes.NewBuffer(data),
		int64(len(data)),
		mapToPutObjOptions(metadata),
	)
	return err
}

func (m *minioStorage) PutObjectWithMeta(
	bucket string,
	key string,
	data []byte,
	metadata *Metadata,
) error {
	_, err := m.client.PutObject(
		context.Background(),
		bucket,
		key,
		bytes.NewBuffer(data),
		int64(len(data)),
		metadataToPutObjOptions(metadata),
	)
	return err
}

func (m *minioStorage) PutFile(
	bucket string,
	key string,
	localFile string,
	metadata map[string]string,
) error {
	_, err := m.client.FPutObject(
		context.Background(),
		bucket,
		key,
		localFile,
		mapToPutObjOptions(metadata),
	)
	return err
}

func (m *minioStorage) PutFileWithMeta(
	bucket string,
	key string,
	srcFile string,
	metadata *Metadata,
) error {
	_, err := m.client.FPutObject(
		context.Background(),
		bucket,
		key,
		srcFile,
		metadataToPutObjOptions(metadata),
	)
	return err
}

// PutObjectFromFile 注意小文件可以，大文件不能走这个函数
func (m *minioStorage) PutObjectFromFile(
	bucket string,
	key string,
	localFile string,
	metadata map[string]string,
) error {
	return nil
}

// PutFileWithPart 支持分片上传，支持5GB以上的文件上传
// 注意 断点续传上传接口传入的文件大小至少要100K以上
func (m *minioStorage) PutFileWithPart(
	bucket string,
	key string,
	srcFile string,
	metadata *Metadata,
	partSize int64,
) error {
	return nil
}

func (m *minioStorage) ListObjects(bucket string, prefix string) ([]Content, error) {
	ctx := context.Background()
	res := make([]Content, 0, 32)
	// List all objects from a bucket-name with a matching prefix.
	for object := range m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if object.Err != nil {
			return res, object.Err
		}
		// protect
		if len(res) >= 100000 {
			return res, errors.New("minio ListObjects too much data, break")
		}
		res = append(res, *objectInfoToContent(&object))
	}
	return res, nil
}

func (m *minioStorage) DeleteObject(bucket string, key string) error {
	ctx := context.Background()
	return m.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (m *minioStorage) BatchDeleteObject(
	bucket string,
	filelist []string,
) (successList []string, err error) {
	ctx := context.Background()
	successList = make([]string, 0, len(filelist))

	for len(filelist) > 0 {
		var splitPos int
		if len(filelist) > MAX_MINIO_BATCH_DELETE_COUNT {
			splitPos = MAX_MINIO_BATCH_DELETE_COUNT
		} else {
			splitPos = len(filelist)
		}
		toDelFileList := filelist[:splitPos]
		filelist = filelist[splitPos:]

		objectsCh := make(chan minio.ObjectInfo)

		// 发送对象以供删除
		go func() {
			defer close(objectsCh)
			for _, object := range toDelFileList {
				objectsCh <- minio.ObjectInfo{Key: object}
			}
		}()

		// 执行批量删除
		fmt.Println("Deleting objects...")
		for err := range m.client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{}) {
			if err.Err != nil {
				fmt.Printf("Error detected during deletion: %v\n", err)
			} else {
				successList = append(successList, err.ObjectName)
			}
		}
		fmt.Println("Deletion completed for this batch")
	}

	return successList, nil
}

func (m *minioStorage) CopyObject(bucket string, srcKey string, destKey string) error {
	ctx := context.Background()
	// Source object
	srcOpts := minio.CopySrcOptions{
		Bucket: bucket,
		Object: srcKey,
	}

	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: bucket,
		Object: destKey,
	}

	// Copy object call
	_, err := m.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *minioStorage) SetObjectAcl(bucket string, key string, acl StorageAcl) error {
	// return errors.New("minio not support SetObjectAcl")
	return nil
}

func (m *minioStorage) SetObjectMetaData(bucket string, key string, metadata *Metadata) error {
	// return errors.New("minio not support SetObjectMeta")
	return nil
}

func (m *minioStorage) IsObjectExist(bucket string, key string) (bool, error) {
	res, err := m.GetObjectMeta(bucket, key)
	switch err := err.(type) {
	case minio.ErrorResponse:
		if err.Code == "NoSuchKey" {
			return false, nil
		} else {
			return false, err
		}
	default:
		return res != nil, err
	}

}

func (m *minioStorage) GetDirToken2(remoteDir string) *StorageToken {
	expires := 6 * 3600 * time.Second
	li, err := credentials.NewSTSAssumeRole(
		"http://"+m.config.Endpoint,
		credentials.STSAssumeRoleOptions{
			AccessKey:       "rw_client",
			SecretKey:       "#$infi0831",
			DurationSeconds: int(expires),
		},
	)
	if err != nil {
		return nil
	}
	to, err := li.Get()
	if err != nil {
		return nil
	}
	return &StorageToken{
		AccessKeyID:     to.AccessKeyID,
		AccessKeySecret: to.SecretAccessKey,
		StsToken:        to.SessionToken,
		Bucket:          m.config.Bucket,
		Region:          m.config.Region,
		Provider:        m.config.Provider,
		Expire:          time.Now().Add(expires).UnixMilli(),
		UploadPath:      remoteDir,
		Host:            m.config.Host,
		Path:            remoteDir,
		CdnDomain:       m.config.CdnDomain,
	}
}

func (m *minioStorage) GetDirToken(remoteDir string) map[string]any {
	t := m.GetDirToken2(remoteDir)
	if t == nil {
		return nil
	}
	res := map[string]any{
		"accessKeyId":     t.AccessKeyID,
		"accessKeySecret": t.AccessKeySecret,
		"stsToken":        t.StsToken,
		"bucket":          t.Bucket,
		"region":          t.Region,
		"provider":        t.Provider,
		"expire":          t.Expire,
		"uploadPath":      t.UploadPath,
		"host":            t.Host,
	}
	return res

}

func (m *minioStorage) GetDirTokenWithAction(
	remoteDir string,
	actions ...Action,
) (bool, *StorageToken) {

	return true, nil
}

func (m *minioStorage) SignFile(remoteDir string, expiredTime int64) string {
	return m.SignFile2(m.config.Bucket, remoteDir, expiredTime)
}

func (m *minioStorage) SignFile2(bucket, remoteDir string, expiredTime int64) string {
	ctx := context.Background()
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	// Generates a presigned url which expires in a day.
	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		bucket,
		remoteDir,
		time.Duration(expiredTime)*time.Second,
		reqParams,
	)
	if err != nil {
		return ""
	}
	return presignedURL.String()
}

// 过期时间:秒
func (m *minioStorage) SignFileForDownload(
	remoteFilePath string,
	expiredTime int64,
	downLoadFilename string,
) string {
	ctx := context.Background()
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set(
		"response-content-disposition",
		"attachment; filename=\""+url.PathEscape(downLoadFilename)+"\"",
	)

	// Generates a presigned url which expires in a day.
	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		"mybucket",
		"myobject",
		time.Duration(expiredTime)*time.Second,
		reqParams,
	)
	if err != nil {
		return ""
	}
	return presignedURL.String()
}

func (m *minioStorage) GetObjectMeta(bucket string, key string) (*Content, error) {
	ctx := context.Background()
	objInfo, err := m.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return objectInfoToContent(&objInfo), nil
}

func mapToPutObjOptions(metadata map[string]string) minio.PutObjectOptions {
	ops := minio.PutObjectOptions{}
	if len(metadata) == 0 {
		return ops
	}
	for key, value := range metadata {
		switch key {
		case "Content-Encoding":
			ops.ContentEncoding = value
		case "Content-Disposition":
			ops.ContentDisposition = value
		case "mime":
			ops.ContentType = value
		default:
			// do nothing
			// logger.Warn("cant support metadata %s %s", metaKey, metaValue)
		}
	}
	return ops
}

func metadataToPutObjOptions(metadata *Metadata) minio.PutObjectOptions {
	ops := minio.PutObjectOptions{}
	if metadata != nil {
		ops.ContentType = metadata.Mime
		ops.ContentEncoding = metadata.ContentEncoding
		ops.ContentDisposition = metadata.ContentDisposition
	}
	return ops
}

func objectInfoToContent(obj *minio.ObjectInfo) *Content {
	res := &Content{
		Key:          obj.Key,
		Size:         obj.Size,
		ETag:         obj.ETag,
		LastModified: obj.LastModified,
	}
	return res
}
