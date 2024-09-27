package maiMeng

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	apiVersion = "api/v2"
)

// 文档地址: https://www.tapd.cn/57919927/markdown_wikis/show/#1157919927001002316@toc4
type MaimengClient struct {
	BaseURL      string
	TenantID     string
	AccessSecret string
	SchoolID     string
}

func NewMaimengClient(baseURL, tenantID, accessSecret, schoolID string) *MaimengClient {
	return &MaimengClient{
		BaseURL:      baseURL,
		TenantID:     tenantID,
		AccessSecret: accessSecret,
		SchoolID:     schoolID,
	}
}

func (c *MaimengClient) generateSignature(endpoint string, exp int64, extenal *string) string {
	var stringToSign string
	if extenal == nil {
		stringToSign = fmt.Sprintf("api:%s:%s:%d:%s", endpoint, c.SchoolID, exp, c.AccessSecret)
	} else {
		stringToSign = fmt.Sprintf("api:%s:%s:%s:%d:%s", endpoint, c.SchoolID, *extenal, exp, c.AccessSecret)
	}
	hash := sha256.Sum256([]byte(stringToSign))
	return hex.EncodeToString(hash[:])
}

func (c *MaimengClient) callAPI(
	method, endpoint string,
	request interface{},
	extenal *string,
	urlParams *url.Values,
) ([]byte, error) {
	var jsonData []byte
	if request != nil {
		d, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		jsonData = d
	}

	resp, err := c.doRequest(method, endpoint, jsonData, extenal, urlParams)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *MaimengClient) doRequest(
	method, endpoint string,
	jsonData []byte,
	extenal *string,
	urlParams *url.Values,
) ([]byte, error) {
	exp := time.Now().Add(10 * time.Minute).Unix()
	signature := c.generateSignature(endpoint, exp, extenal)
	urlStr := fmt.Sprintf("%s/%s/%s/%s", c.BaseURL, apiVersion, c.SchoolID, endpoint)
	if urlParams == nil {
		urlParams = &url.Values{}
	}
	urlParams.Add("tenant_id", c.TenantID)
	urlParams.Add("exp", strconv.FormatInt(exp, 10))
	urlParams.Add("signature", signature)

	fullURL := fmt.Sprintf("%s?%s", urlStr, urlParams.Encode())

	client := resty.New()
	r := client.R()
	r.SetHeader("Content-Type", "application/json")
	if len(jsonData) > 0 {
		r.SetBody(jsonData)
	}
	resp, err := r.Execute(method, fullURL)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

type BaseResp struct {
	Code    int64  `json:"code"` // code 0表示成功
	Message string `json:"message"`
}

type CreateSchoolRequest struct {
	SchoolName string `json:"school_name"` // 学校名称
}

type CreateSchoolResponse struct {
	BaseResp
}

// 创建学校
// POST /api/v2/{SCHOOL_ID}/create_school
// 第三方系统发起的请求中包含学校名、学段、年级、学科信息，麦盟收到后，初始化学校信息
// 待签名字符串: string_to_sign = "api:create_school:" + $SCHOOL_ID + ":" + exp + ":" + $AccessSecret
func (c *MaimengClient) CreateSchool(request *CreateSchoolRequest) (*CreateSchoolResponse, error) {
	resp, err := c.callAPI("POST", "create_school", request, nil, nil)
	if err != nil {
		return nil, err
	}

	var response CreateSchoolResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CreateLectureRequest 代表创建课程请求的结构体
// CreateLectureRequest 代表创建课程请求的结构体
type CreateLectureRequest struct {
	// GradeName 年级，如"四年级"，选填，建议填写，如不填，会影响横向对比报表
	GradeName string `json:"grade_name"`

	// ClassName 班级，如"二班"，选填，建议填写，如不填，会影响横向对比报表
	ClassName string `json:"class_name"`

	// SubjectName 学科，如"英语"，选填，建议填写，如不填，默认按"语文"处理，对于真实上的是英语课的情况会有较大影响
	SubjectName string `json:"subject_name"`

	// CourseName 课程名称，如"小英雄雨来"，选填，建议填写
	CourseName string `json:"course_name"`

	// TeacherName 老师姓名，如"陈小静"，选填，建议填写，如不填，默认按"匿名老师"处理，会影响横向对比报表
	TeacherName string `json:"teacher_name"`

	// ExternalTeacherID 客户端系统的老师ID，选填，建议填写，用以处理同名老师的情况
	ExternalTeacherID string `json:"external_teacher_id"`

	// SchoolYear 学年，如"2022-2023学年"，选填，如不填，默认当前学年
	SchoolYear string `json:"school_year"`

	// SchoolTerm 学期，如"上学期"，选填，如传入，只能为"上学期"或"下学期"之一，如不填，以当前学期为准
	SchoolTerm string `json:"school_term"`

	// ClassroomName 教室名称，如"一楼101室"，选填
	ClassroomName string `json:"classroom_name"`

	// TeacherVideoURL 老师画面视频的下载地址，必填，应为http(s)地址
	TeacherVideoURL string `json:"teacher_video_url"`

	// StudentVideoURL 学生画面视频的下载地址，选填，建议填写，应为http(s)地址
	StudentVideoURL string `json:"student_video_url"`

	// MixedVideoURL 混流画面视频的下载地址，选填，必须有音频流，应为http(s)地址
	MixedVideoURL string `json:"mixed_video_url"`

	// StartedAt 课堂开始时间点，UTC 1970年1月1日0点后的秒数，选填，建议填写
	StartedAt int64 `json:"started_at"`

	// EndedAt 课堂结束时间，UTC 1970年1月1日0点后的秒数，选填，建议填写
	EndedAt int64 `json:"ended_at"`

	// Callback 分析完成后的回调地址，选填，格式如 "http://ip:port/url"
	Callback string `json:"callback"`

	// TeachingPlanURL 教案下载地址，选填，建议填写，应为docx格式
	TeachingPlanURL string `json:"teaching_plan_url"`

	// SchoolName 学校名称，选填，如"中心小学"
	SchoolName string `json:"school_name"`
}

type CreateLectureResponse struct {
	BaseResp
	Data *CreateLectureData `json:"data"`
}

type CreateLectureData struct {
	ReportSingleUrl string `json:"report_single_url"`
}

// 创建 AI 分析课堂
// POST /api/v2/{SCHOOL_ID}/create_lecture/{LECTURE_ID}
// 当 LECTURE_ID 在麦盟系统不存在，就在麦盟系统创建，否则为更新
// 客户端发起 POST 请求，请求体包含年级、班级、学科、老师等课堂基本信息和录制好的音视频、教案的下载地址信息。麦盟服务器收到请求后，把它加入分析任务队列 ，分析完成后，如果客户端当时指定过回调地址(callback)，则向该地址发送 POST 请求，通知该堂课分析完成，客户端可以更新自有数据库中课堂的状态，以在界面中提供查看报告的链接
// 待签名字符串: string_to_sign = "api:create_lecture:" + $SCHOOL_ID + ":" + $LECTURE_ID + ":" + exp + ":" + $AccessSecret
func (c *MaimengClient) CreateLecture(
	lectureId string,
	request *CreateLectureRequest,
) (*CreateLectureResponse, error) {
	resp, err := c.callAPI("POST", "create_lecture", request, &lectureId, nil)
	if err != nil {
		return nil, err
	}

	var response CreateLectureResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

type DeleteLectureResponse struct {
	BaseResp
}

// 删除 AI 分析课堂
// DELETE /api/v2/{SCHOOL_ID}/delete_lecture/{LECTURE_ID}
// 用以在麦盟端删除某个指定 LECTURE_ID 的 AI 分析课堂。如果这节课还没开始或已结束，则正常删除，如果在分析中，则停止分析任务后删除
// 待签名字符串：string_to_sign = "api:delete_lecture:" + $SCHOOL_ID + ":" + $LECTURE_ID + ":" + exp + ":" + $AccessSecret
func (c *MaimengClient) DeleteLecture(
	lectureId string,
) (*DeleteLectureResponse, error) {
	resp, err := c.callAPI("DELETE", "delete_lecture", nil, &lectureId, nil)
	if err != nil {
		return nil, err
	}

	var response DeleteLectureResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

type QueryLectureResponse struct {
	BaseResp
	Data *LectureStatusData `json:"data"`
}

type LectureStatusData struct {
	SchoolId         string `json:"school_id"`          // 学校 ID
	LectureId        string `json:"lecture_id"`         // 课堂 ID
	Status           string `json:"status"`             // 分析状态: 待分析 分析中 分析完成 分析失败
	LectureReportUrl string `json:"lecture_report_url"` // 单课报告地址链接，只有 status 为"分析完成"时才应该展示给用户
}

// 查询 AI 分析课堂状态接口
// GET /api/v2/{SCHOOL_ID}/query_lecture/{LECTURE_ID}
// URL 中，包含待查询的课堂 ID，本接口返回该课堂的分析状态
// 待签名字符串: string_to_sign = "api:query_lecture:" + $SCHOOL_ID + ":" + $LECTURE_ID + ":" + exp + ":" + $AccessSecret
func (c *MaimengClient) QueryLecture(
	lectureId string,
) (*QueryLectureResponse, error) {
	resp, err := c.callAPI("GET", "query_lecture", nil, &lectureId, nil)
	if err != nil {
		return nil, err
	}

	var response QueryLectureResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

type ListLectureRequest struct {
	Status string // 任务状态，为待分析/分析中/分析完成/分析失败之一，如不传，默认为全部
	Page   string // 分页页号，Int，从 1 开始，如不传，默认为 1
	Limit  string // 每页返回记录数，最大不能超过 200，如不传，默认为 20
}

type ListLectureResponse struct {
	BaseResp
	Data *ListLectureData `json:"data"`
}

type ListLectureData struct {
	Items      []*LectureStatusData `json:"items"`      // 学校 ID
	Pagination *Pagination          `json:"pagination"` // 课堂 ID
}

type Pagination struct {
	Total int64 `json:"total"` // 总数
}

// AI 分析课堂列表
// GET /api/v2/{SCHOOL_ID}/list_lecture
// 表示查询解析任务列表
// 待签名字符串: string_to_sign = "api:list_lecture:" + $SCHOOL_ID + ":" + exp + ":" + $AccessSecret
func (c *MaimengClient) ListLecture(
	lectureId string,
	request *ListLectureRequest,
) (*ListLectureResponse, error) {
	urlParams := &url.Values{}
	if len(request.Status) > 0 {
		urlParams.Add("status", request.Status)
	}
	if len(request.Page) > 0 {
		urlParams.Add("page", request.Page)
	}
	if len(request.Limit) > 0 {
		urlParams.Add("limit", request.Limit)
	}

	resp, err := c.callAPI("GET", "list_lecture", nil, &lectureId, urlParams)
	if err != nil {
		return nil, err
	}

	var response ListLectureResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
