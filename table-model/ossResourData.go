package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `oss_resource_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) NOT NULL COMMENT '学校名称',
	`summary_start_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`summary_end_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`resource_upload_count` bigint NOT NULL COMMENT '上传资源总数',
	`resource_upload_size` bigint NOT NULL COMMENT '上传资源总大小',
	`resource_upload_parse` bigint NOT NULL COMMENT '上传资源总解析时间',
	`resource_upload_view` bigint NOT NULL COMMENT '上传资源总浏览次数',
	`resource_type_count_ppt` bigint NOT NULL COMMENT '上传资源类型数量_ppt',
	`resource_type_size_ppt` bigint NOT NULL COMMENT '上传资源类型大小_ppt',
	`resource_type_parse_ppt` bigint NOT NULL COMMENT '上传资源类型解析时间_ppt',
	`resource_type_view_ppt` bigint NOT NULL COMMENT '上传资源类型观看次数_ppt',
	`resource_type_count_pdf` bigint NOT NULL COMMENT '上传资源类型数量_pdf',
	`resource_type_size_pdf` bigint NOT NULL COMMENT '上传资源类型大小_pdf',
	`resource_type_parse_pdf` bigint NOT NULL COMMENT '上传资源类型解析时间_pdf',
	`resource_type_view_pdf` bigint NOT NULL COMMENT '上传资源类型观看次数_pdf',
	`resource_type_count_word` bigint NOT NULL COMMENT '上传资源类型数量_word',
	`resource_type_size_word` bigint NOT NULL COMMENT '上传资源类型大小_word',
	`resource_type_parse_word` bigint NOT NULL COMMENT '上传资源类型解析时间_word',
	`resource_type_view_word` bigint NOT NULL COMMENT '上传资源类型观看次数_word',
	`resource_type_count_excel` bigint NOT NULL COMMENT '上传资源类型数量_excel',
	`resource_type_size_excel` bigint NOT NULL COMMENT '上传资源类型大小_excel',
	`resource_type_parse_excel` bigint NOT NULL COMMENT '上传资源类型解析时间_excel',
	`resource_type_view_excel` bigint NOT NULL COMMENT '上传资源类型观看次数_excel',
	`resource_type_count_pic` bigint NOT NULL COMMENT '上传资源类型数量_pic',
	`resource_type_size_pic` bigint NOT NULL COMMENT '上传资源类型大小_pic',
	`resource_type_parse_pic` bigint NOT NULL COMMENT '上传资源类型解析时间_pic',
	`resource_type_view_pic` bigint NOT NULL COMMENT '上传资源类型观看次数_pic',
	`resource_type_count_video` bigint NOT NULL COMMENT '上传资源类型数量_video',
	`resource_type_size_video` bigint NOT NULL COMMENT '上传资源类型大小_video',
	`resource_type_parse_video` bigint NOT NULL COMMENT '上传资源类型解析时间_video',
	`resource_type_view_video` bigint NOT NULL COMMENT '上传资源类型观看次数_video',
	`resource_type_count_audio` bigint NOT NULL COMMENT '上传资源类型数量_audio',
	`resource_type_size_audio` bigint NOT NULL COMMENT '上传资源类型大小_audio',
	`resource_type_parse_audio` bigint NOT NULL COMMENT '上传资源类型解析时间_audio',
	`resource_type_view_audio` bigint NOT NULL COMMENT '上传资源类型观看次数_audio',
	`resource_type_count_compress` bigint NOT NULL COMMENT '上传资源类型数量_compress',
	`resource_type_size_compress` bigint NOT NULL COMMENT '上传资源类型大小_compress',
	`resource_type_parse_compress` bigint NOT NULL COMMENT '上传资源类型解析时间_compress',
	`resource_type_view_compress` bigint NOT NULL COMMENT '上传资源类型观看次数_compress',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `summary_start_date` (`summary_start_date`) USING BTREE COMMENT '统计日期索引',
	KEY `summary_end_date` (`summary_end_date`) USING BTREE COMMENT '统计日期索引',
	UNIQUE KEY `org_date_indx` (`org_id`,`summary_start_date`,`summary_end_date`) USING BTREE

) ENGINE=InnoDB COMMENT = 'oss资源汇总信息表';
*/
type OssResourceSummary struct {
	ID                        int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                     int64                  `gorm:"column:org_id"                                   json:"orgId"`                     // 校区ID
	OrgName                   string                 `gorm:"column:org_name"                                 json:"orgName"`                   // 机构名称
	SummaryStartDate          int64                  `gorm:"column:summary_start_date"                       json:"summaryStartDate"`          // 统计开始日期；精确到天,单位毫秒
	SummaryEndDate            int64                  `gorm:"column:summary_end_date"                         json:"summaryEndDate"`            // 统计开始日期；精确到天,单位毫秒
	ResourceUploadCount       int64                  `gorm:"column:resource_upload_count"                    json:"resourceUploadCount"`       // 上传资源总数
	ResourceUploadSize        int64                  `gorm:"column:resource_upload_size"                     json:"resourceUploadSize"`        // 上传资源总大小
	ResourceUploadParse       int64                  `gorm:"column:resource_upload_parse"                    json:"resourceUploadParse"`       // 上传资源总解析时间
	ResourceUploadView        int64                  `gorm:"column:resource_upload_view"                     json:"resourceUploadView"`        // 上传资源总浏览次数
	ResourceTypeCountPpt      int64                  `gorm:"column:resource_type_count_ppt"                  json:"resourceTypeCountPpt"`      // 上传资源类型数量_ppt
	ResourceTypeSizePpt       int64                  `gorm:"column:resource_type_size_ppt"                   json:"resourceTypeSizePpt"`       // 上传资源类型大小_ppt
	ResourceTypeParsePpt      int64                  `gorm:"column:resource_type_parse_ppt"                  json:"resourceTypeParsePpt"`      // 上传资源类型解析时间_ppt
	ResourceTypeViewPpt       int64                  `gorm:"column:resource_type_view_ppt"                   json:"resourceTypeViewPpt"`       // 上传资源类型观看次数_ppt
	ResourceTypeCountPdf      int64                  `gorm:"column:resource_type_count_pdf"                  json:"resourceTypeCountPdf"`      // 上传资源类型数量_pdf
	ResourceTypeSizePdf       int64                  `gorm:"column:resource_type_size_pdf"                   json:"resourceTypeSizePdf"`       // 上传资源类型大小_pdf
	ResourceTypeParsePdf      int64                  `gorm:"column:resource_type_parse_pdf"                  json:"resourceTypeParsePdf"`      // 上传资源类型解析时间_pdf
	ResourceTypeViewPdf       int64                  `gorm:"column:resource_type_view_pdf"                   json:"resourceTypeViewPdf"`       // 上传资源类型观看次数_pdf
	ResourceTypeCountWord     int64                  `gorm:"column:resource_type_count_word"                 json:"resourceTypeCountWord"`     // 上传资源类型数量_word
	ResourceTypeSizeWord      int64                  `gorm:"column:resource_type_size_word"                  json:"resourceTypeSizeWord"`      // 上传资源类型大小_word
	ResourceTypeParseWord     int64                  `gorm:"column:resource_type_parse_word"                 json:"resourceTypeParseWord"`     // 上传资源类型解析时间_word
	ResourceTypeViewWord      int64                  `gorm:"column:resource_type_view_word"                  json:"resourceTypeViewWord"`      // 上传资源类型观看次数_word
	ResourceTypeCountExcel    int64                  `gorm:"column:resource_type_count_excel"                json:"resourceTypeCountExcel"`    // 上传资源类型数量_excel
	ResourceTypeSizeExcel     int64                  `gorm:"column:resource_type_size_excel"                 json:"resourceTypeSizeExcel"`     // 上传资源类型大小_excel
	ResourceTypeParseExcel    int64                  `gorm:"column:resource_type_parse_excel"                json:"resourceTypeParseExcel"`    // 上传资源类型解析时间_excel
	ResourceTypeViewExcel     int64                  `gorm:"column:resource_type_view_excel"                 json:"resourceTypeViewExcel"`     // 上传资源类型观看次数_excel
	ResourceTypeCountPic      int64                  `gorm:"column:resource_type_count_pic"                  json:"resourceTypeCountPic"`      // 上传资源类型数量_pic
	ResourceTypeSizePic       int64                  `gorm:"column:resource_type_size_pic"                   json:"resourceTypeSizePic"`       // 上传资源类型大小_pic
	ResourceTypeParsePic      int64                  `gorm:"column:resource_type_parse_pic"                  json:"resourceTypeParsePic"`      // 上传资源类型解析时间_pic
	ResourceTypeViewPic       int64                  `gorm:"column:resource_type_view_pic"                   json:"resourceTypeViewPic"`       // 上传资源类型观看次数_pic
	ResourceTypeCountVideo    int64                  `gorm:"column:resource_type_count_video"                json:"resourceTypeCountVideo"`    // 上传资源类型数量_video
	ResourceTypeSizeVideo     int64                  `gorm:"column:resource_type_size_video"                 json:"resourceTypeSizeVideo"`     // 上传资源类型大小_video
	ResourceTypeParseVideo    int64                  `gorm:"column:resource_type_parse_video"                json:"resourceTypeParseVideo"`    // 上传资源类型解析时间_video
	ResourceTypeViewVideo     int64                  `gorm:"column:resource_type_view_video"                 json:"resourceTypeViewVideo"`     // 上传资源类型观看次数_video
	ResourceTypeCountAudio    int64                  `gorm:"column:resource_type_count_audio"                json:"resourceTypeCountAudio"`    // 上传资源类型数量_audio
	ResourceTypeSizeAudio     int64                  `gorm:"column:resource_type_size_audio"                 json:"resourceTypeSizeAudio"`     // 上传资源类型大小_audio
	ResourceTypeParseAudio    int64                  `gorm:"column:resource_type_parse_audio"                json:"resourceTypeParseAudio"`    // 上传资源类型解析时间_audio
	ResourceTypeViewAudio     int64                  `gorm:"column:resource_type_view_audio"                 json:"resourceTypeViewAudio"`     // 上传资源类型观看次数_audio
	ResourceTypeCountCompress int64                  `gorm:"column:resource_type_count_compress"             json:"resourceTypeCountCompress"` // 上传资源类型数量_compress
	ResourceTypeSizeCompress  int64                  `gorm:"column:resource_type_size_compress"              json:"resourceTypeSizeCompress"`  // 上传资源类型大小_compress
	ResourceTypeParseCompress int64                  `gorm:"column:resource_type_parse_compress"             json:"resourceTypeParseCompress"` // 上传资源类型解析时间_compress
	ResourceTypeViewCompress  int64                  `gorm:"column:resource_type_view_compress"              json:"resourceTypeViewCompress"`  // 上传资源类型观看次数_compress
	CreatedAt                 int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                         // 创建时间
	UpdatedAt                 int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                         // 更新时间
	DeletedAt                 *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                         // 删除时间
}

func (o *OssResourceSummary) TableName() string {
	return "oss_resource_summary"
}
