package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioHelper         *minioStorage
	configMinioTestPath = "/home/fangyuan/work/rongke/statistical_service/test/minio_conf.json"
	// 本地测试文件
	sourceMinioLocalTestFile = "./tmp.txt"
	configMinio              = &StorageConfig{}
)

// 初始化客户端配置
func initMinioClient() {
	data, err := os.ReadFile(
		configMinioTestPath,
	)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &configMinio)
	fmt.Printf("config: %+v\n", configMinio)
	minioHelper, err = newMinioStorager(configMinio)
	if err != nil {
		panic(err)
	}
}

func TestGetObject(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestGetObject",
		[]byte("这是什么东西"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(context.Background(), configMinio.Bucket, "TestGetObject")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestGetFile(t *testing.T) {
	initMinioClient()
	err := minioHelper.GetFile(context.Background(), configMinio.Bucket, "12344.png", "./test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPutObject(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestPutObject",
		[]byte("这是谁的部下"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(context.Background(), configMinio.Bucket, "TestPutObject")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestPutObjectWithMeta(t *testing.T) {
	initMinioClient()
	metadata := &Metadata{
		Mime:               "text/plain",
		ContentDisposition: "TestPutObjectWithMeta",
	}
	err := minioHelper.PutObjectWithMeta(
		context.Background(),
		configMinio.Bucket,
		"PutObjectWithMeta",
		[]byte("吾乃长山赵子龙"),
		metadata,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(
		context.Background(),
		configMinio.Bucket,
		"PutObjectWithMeta",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestPutFile(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutFile(
		context.Background(),
		configMinio.Bucket,
		"TestPutFile",
		sourceMinioLocalTestFile,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(context.Background(), configMinio.Bucket, "TestPutFile")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestPutFileWithMeta(t *testing.T) {
	initMinioClient()
	metadata := &Metadata{
		Mime:               "text/plain",
		ContentDisposition: "TestPutFileWithMeta",
	}
	minioHelper.PutFileWithMeta(
		context.Background(),
		configMinio.Bucket,
		"TestPutFileWithMeta",
		sourceMinioLocalTestFile,
		metadata,
	)
	data, err := minioHelper.GetObject(
		context.Background(),
		configMinio.Bucket,
		"TestPutFileWithMeta",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestPutFileWithPart(t *testing.T) {
	initMinioClient()
	metadata := &Metadata{
		ContentDisposition: "iso file",
	}
	// 定义 MiB 的字节数
	const MiB int64 = 1 << 20 // 1,048,576 字节

	// 计算 128 MiB 的字节数
	partSize := int64(128 * MiB)
	// 测试文件大小为 5.3G
	err := minioHelper.PutFileWithPart(
		context.Background(),
		configMinio.Bucket,
		"test_vedio.iso",
		"/home/fangyuan/下载/bak/Win10_22H2_China_GGK_Chinese_Simplified_x64.iso",
		metadata,
		partSize,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListObjects(t *testing.T) {
	initMinioClient()
	datas, err := minioHelper.ListObjects(context.Background(), configMinio.Bucket, "zppt10401")
	if err != nil {
		t.Fatal(err)
	}
	for i := range datas {
		t.Logf("%v", datas[i])
	}
}

func TestDeleteObject(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestDeleteObject",
		[]byte("这是什么东西"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(context.Background(), configMinio.Bucket, "TestDeleteObject")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
	err = minioHelper.DeleteObject(context.Background(), configMinio.Bucket, "TestDeleteObject")
	if err != nil {
		t.Fatal(err)
	}
	_, err = minioHelper.GetObject(context.Background(), configMinio.Bucket, "TestDeleteObject")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestBatchDeleteObject(t *testing.T) {
	initMinioClient()
	fileList := []string{
		"PutObjectWithMeta",
		"TestGetObject",
		"TestPutFile",
		"TestPutFileWithMeta",
		"TestPutObject",
	}
	successFiles, err := minioHelper.BatchDeleteObject(
		context.Background(),
		configMinio.Bucket,
		fileList,
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range successFiles {
		t.Log("success delete file: ", file)
	}
}

func TestCopyObject(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestCopyObject",
		[]byte("这是什么东西"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = minioHelper.CopyObject(
		context.Background(),
		configMinio.Bucket,
		"TestCopyObject",
		"TestCopyObject_COPY",
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := minioHelper.GetObject(
		context.Background(),
		configMinio.Bucket,
		"TestCopyObject_COPY",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestIsObjectExist(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestIsObjectExist",
		[]byte("这是什么东西"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	isObjectExist, err := minioHelper.IsObjectExist(
		context.Background(),
		configMinio.Bucket,
		"TestIsObjectExist",
	)
	if err != nil {
		t.Fatal(err)
	}
	if isObjectExist {
		t.Log("TestIsObjectExist: 存在")
	}
	err = minioHelper.DeleteObject(context.Background(), configMinio.Bucket, "TestIsObjectExist")
	if err != nil {
		t.Fatal(err)
	}
	isObjectExist, err = minioHelper.IsObjectExist(
		context.Background(),
		configMinio.Bucket,
		"TestIsObjectExist",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !isObjectExist {
		t.Log("TestIsObjectExist: 不存在")
	}
}

// 过期时间:秒
func TestSignFileForDownload(t *testing.T) {
	initMinioClient()
	downloadUrl, err := minioHelper.SignFileForDownload(
		context.Background(),
		"zppt10401/zppt10401.pptx",
		3600,
		"爱多福多寿.pptx",
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("downloadUrl: %s", downloadUrl)
}

func TestGetObjectMeta(t *testing.T) {
	initMinioClient()
	err := minioHelper.PutObject(
		context.Background(),
		configMinio.Bucket,
		"TestGetObjectMeta",
		[]byte("这是什么东西"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	content, err := minioHelper.GetObjectMeta(
		context.Background(),
		configMinio.Bucket,
		"TestGetObjectMeta",
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("content: %v", content)
}

func TestSignFile(t *testing.T) {
	initMinioClient()
	url, err := minioHelper.SignFile(context.Background(), "zppt10401/zppt10401_s.jpg", 3600)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("url: %s", url)
}

func TestSignFile2(t *testing.T) {
	initMinioClient()
	url, err := minioHelper.SignFile2(
		context.Background(),
		configMinio.Bucket,
		"zppt10401/zppt10401_s.jpg",
		3600,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("url: %s", url)
}

// 下述方法采用兼容支持策略

// GetDirToken,GetDirTokenWithAction都是依赖使用GetDirToken2实现,就不统一测试了
func TestGetDirToken2(t *testing.T) {
	initMinioClient()
	token, err := minioHelper.GetDirToken2(context.Background(), "tmp-file")
	if err != nil {
		t.Fatal(err)
	}

	creds := credentials.NewStaticV4(token.AccessKeyID, token.AccessKeySecret, token.StsToken)
	// Use generated credentials to authenticate with MinIO server
	minioClient, err := minio.New(
		configMinio.Endpoint,
		&minio.Options{
			Creds:  creds,
			Secure: false,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("TestGetDirToken2\nTestGetDirToken2")
	_, err = minioClient.PutObject(
		context.Background(),
		configMinio.Bucket,
		"tmp-file/TestCopyObject",
		bytes.NewBuffer(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf(
		"Calling list objects on bucket named `%s` with temp creds:\n===\n",
		configMinio.Bucket,
	)
	objCh := minioClient.ListObjects(
		context.Background(),
		configMinio.Bucket,
		minio.ListObjectsOptions{
			Prefix:    "",
			Recursive: true,
		},
	)
	for obj := range objCh {
		if obj.Err != nil {
			t.Fatalf("Listing error: %v", obj.Err)
		}
		fmt.Printf(
			"Key: %s\nSize: %d\nLast Modified: %s\n===\n",
			obj.Key,
			obj.Size,
			obj.LastModified,
		)
	}
}
