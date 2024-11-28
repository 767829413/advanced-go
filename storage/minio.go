package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
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
	// Initialize minio client object.
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyId, c.AccessKeySecret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &minioStorage{config: c, client: client}, nil
}

func (m *minioStorage) GetObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	return buf.Bytes(), err
}

func (m *minioStorage) GetFile(
	ctx context.Context,
	bucket string,
	key string,
	localFile string,
) error {
	return m.client.FGetObject(
		ctx,
		bucket,
		key,
		localFile,
		minio.GetObjectOptions{},
	)
}

func (m *minioStorage) PutObject(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	metadata map[string]string,
) error {
	opt, err := mapToPutObjOptions(metadata)
	if err != nil {
		return err
	}
	_, err = m.client.PutObject(
		ctx,
		bucket,
		key,
		bytes.NewBuffer(data),
		int64(len(data)),
		opt,
	)
	return err
}

func (m *minioStorage) PutObjectWithMeta(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	metadata *Metadata,
) error {
	_, err := m.client.PutObject(
		ctx,
		bucket,
		key,
		bytes.NewBuffer(data),
		int64(len(data)),
		metadataToPutObjOptions(metadata),
	)
	return err
}

func (m *minioStorage) PutFile(
	ctx context.Context,
	bucket string,
	key string,
	localFile string,
	metadata map[string]string,
) error {
	opt, err := mapToPutObjOptions(metadata)
	if err != nil {
		return err
	}
	_, err = m.client.FPutObject(
		ctx,
		bucket,
		key,
		localFile,
		opt,
	)
	return err
}

func (m *minioStorage) PutFileWithMeta(
	ctx context.Context,
	bucket string,
	key string,
	srcFile string,
	metadata *Metadata,
) error {
	_, err := m.client.FPutObject(
		ctx,
		bucket,
		key,
		srcFile,
		metadataToPutObjOptions(metadata),
	)
	return err
}

func (m *minioStorage) PutObjectFromFile(
	ctx context.Context,
	bucket string,
	key string,
	localFile string,
	metadata map[string]string,
) error {
	opt, err := mapToPutObjOptions(metadata)
	if err != nil {
		return err
	}
	// 使用FPutObject上传文件
	_, err = m.client.FPutObject(ctx, bucket, key, localFile, opt)
	if err != nil {
		return fmt.Errorf("failed to put object from file: %w", err)
	}
	return nil
}

// 在单个 PUT 操作中上传小于 128MiB 的对象。对于大于 128MiB 的对象，PutObject 会根据实际文件大小将对象无缝上传为 128MiB 或更大的部分。对象的最大上传大小为 5TB
func (m *minioStorage) PutFileWithPart(
	ctx context.Context,
	bucket string,
	key string,
	srcFile string,
	metadata *Metadata,
	partSize int64,
) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("PutFileWithPart os.Open error %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("PutFileWithPart file.Stat error %w", err)
	}

	putOpt := metadataToPutObjOptions(metadata)
	// 这里是直接覆盖
	putOpt.ContentType = "application/octet-stream"
	_, err = m.client.PutObject(
		ctx,
		bucket,
		key,
		file,
		fileStat.Size(),
		putOpt,
	)
	if err != nil {
		return fmt.Errorf("PutFileWithPart m.client.PutObject error %w", err)
	}
	return nil
}

func (m *minioStorage) ListObjects(
	ctx context.Context,
	bucket string,
	prefix string,
) ([]Content, error) {
	res := make([]Content, 0, 32)
	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	// List all objects from a bucket-name with a matching prefix.
	for object := range objectCh {
		if object.Err != nil {
			return res, object.Err
		}
		// protect
		if len(res) >= 100000 {
			return res, errors.New("too much data, break")
		}
		res = append(res, *objectInfoToContent(&object))
	}
	return res, nil
}

func (m *minioStorage) DeleteObject(ctx context.Context, bucket string, key string) error {
	return m.client.RemoveObject(
		ctx,
		bucket,
		key,
		minio.RemoveObjectOptions{GovernanceBypass: true},
	)
}

func (m *minioStorage) BatchDeleteObject(
	ctx context.Context,
	bucket string,
	filelist []string,
) (successList []string, err error) {
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
				objectsCh <- minio.ObjectInfo{
					Key: object,
				}
			}
		}()

		// 执行批量删除
		fmt.Println("Deleting objects...")
		for err := range m.client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
			if err.Err != nil {
				fmt.Printf("Error detected during deletion: %v\n", err)
			}
		}
		fmt.Println("Deletion completed for this batch")
		successList = append(successList, toDelFileList...)
	}

	return successList, nil
}

func (m *minioStorage) CopyObject(
	ctx context.Context,
	bucket string,
	srcKey string,
	destKey string,
) error {
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

func (m *minioStorage) IsObjectExist(ctx context.Context, bucket string, key string) (bool, error) {
	res, err := m.GetObjectMeta(ctx, bucket, key)
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

func (m *minioStorage) SignFile(
	ctx context.Context,
	remoteDir string,
	expiredTime int64,
) (string, error) {
	return m.SignFile2(ctx, m.config.Bucket, remoteDir, expiredTime)
}

func (m *minioStorage) SignFile2(
	ctx context.Context,
	bucket, remoteDir string,
	expiredTime int64,
) (string, error) {
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
		return "", err
	}
	return presignedURL.String(), nil
}

// 过期时间:秒
func (m *minioStorage) SignFileForDownload(
	ctx context.Context,
	remoteFilePath string,
	expiredTime int64,
	downLoadFilename string,
) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set(
		"response-content-disposition",
		"attachment; filename=\""+url.PathEscape(downLoadFilename)+"\"",
	)

	// Generates a presigned url which expires in a day.
	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		m.config.Bucket,
		remoteFilePath,
		time.Duration(expiredTime)*time.Second,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *minioStorage) GetObjectMeta(
	ctx context.Context,
	bucket string,
	key string,
) (*Content, error) {
	objInfo, err := m.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return objectInfoToContent(&objInfo), nil
}

// 获取临时token
// https://github.com/minio/minio/blob/master/docs/sts/assume-role.md
// 使用了minio的账号密码实现,相当于最大权限
func (m *minioStorage) GetDirToken2(ctx context.Context, remoteDir string) (*StorageToken, error) {
	expires := 6 * 3600 * time.Second
	expireDeadLine := time.Now().Add(expires)
	// Initialize credential options
	var stsOpts credentials.STSAssumeRoleOptions
	stsOpts.AccessKey = m.config.Username
	stsOpts.SecretKey = m.config.Password

	li, err := credentials.NewSTSAssumeRole("http://"+m.config.Endpoint, stsOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create STS assume role credential: %w", err)
	}

	to, err := li.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get STS credentials: %w", err)
	}
	return &StorageToken{
		AccessKeyID:     to.AccessKeyID,
		AccessKeySecret: to.SecretAccessKey,
		StsToken:        to.SessionToken,
		Bucket:          m.config.Bucket,
		Region:          m.config.Region,
		Provider:        m.config.Provider,
		Expire:          expireDeadLine.UnixMilli(),
		UploadPath:      remoteDir,
		Host:            m.config.Host,
		Path:            remoteDir,
		CdnDomain:       m.config.CdnDomain,
	}, nil
}

func (m *minioStorage) GetDirToken(ctx context.Context, remoteDir string) (map[string]any, error) {
	t, err := m.GetDirToken2(ctx, remoteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get dir token: %w", err)
	}
	if t == nil {
		return nil, fmt.Errorf("dir token is nil")
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
	return res, nil

}

func (m *minioStorage) GetDirTokenWithAction(
	ctx context.Context,
	remoteDir string,
	actions ...Action,
) (bool, *StorageToken) {
	t, err := m.GetDirToken2(ctx, remoteDir)
	if err != nil {
		return false, nil
	}
	if t == nil {
		return false, nil
	}
	return true, t
}

// TODO: minio not support
func (m *minioStorage) RestoreArchive(bucket string, key string) (bool, error) {
	return false, errors.New("cant support RestoreArchive")
}

// TODO: minio not support
func (m *minioStorage) IsArchive(bucket string, key string) (bool, error) {
	return false, errors.New("cant support IsArchive")
}

func mapToPutObjOptions(metadata map[string]string) (minio.PutObjectOptions, error) {
	ops := minio.PutObjectOptions{}
	if len(metadata) == 0 {
		return ops, nil
	}
	for key, value := range metadata {
		if len(value) == 0 {
			continue
		}
		switch key {
		case "Content-Encoding":
			ops.ContentEncoding = value
		case "Content-Disposition":
			ops.ContentDisposition = value
		case "mime":
			ops.ContentType = value
		default:
			return ops, fmt.Errorf("minio cant support metadata %s %s", key, value)
		}
	}
	return ops, nil
}

func metadataToPutObjOptions(metadata *Metadata) minio.PutObjectOptions {
	ops := minio.PutObjectOptions{}
	if metadata != nil {
		if len(metadata.Mime) != 0 {
			ops.ContentType = metadata.Mime
		}
		if len(metadata.ContentEncoding) != 0 {
			ops.ContentEncoding = metadata.ContentEncoding
		}
		if len(metadata.ContentDisposition) != 0 {
			ops.ContentDisposition = metadata.ContentDisposition
		}
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
