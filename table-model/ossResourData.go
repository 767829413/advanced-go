package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `oss_resour_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学校名称',
	`resource_upload_count` bigint NOT NULL COMMENT '上传资源数',
	`resource_type_count_ppt` bigint NOT NULL COMMENT '资源类型数量-ppt',
	`resource_type_count_pdf` bigint NOT NULL COMMENT '资源类型数量-pdf',
	`resource_type_count_word` bigint NOT NULL COMMENT '资源类型数量-word',
	`resource_type_count_excel` bigint NOT NULL COMMENT '资源类型数量-excel',
	`resource_type_count_pic` bigint NOT NULL COMMENT '资源类型数量-pic',
	`resource_type_count_video` bigint NOT NULL COMMENT '资源类型数量-video',
	`resource_type_count_audio` bigint NOT NULL COMMENT '资源类型数量-audio',
	`resource_type_count_compress` bigint NOT NULL COMMENT '资源类型数量-compress',
	`resource_parse_document` bigint NOT NULL COMMENT '文档类资源解析总数',
	`resource_parse_video` bigint NOT NULL COMMENT '音视频类资源解析总数',
	`oss_total_capacity` bigint NOT NULL COMMENT 'oss总使用量',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `org_id_index` (`org_id`)

) ENGINE=InnoDB COMMENT = '工作台栏目表';
*/
type OssResourSummary struct {
	ID                        int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                     int64                  `gorm:"column:org_id"                                   json:"orgId"`                     // 校区ID
	OrgName                   string                 `gorm:"column:org_name"                                 json:"orgName"`                   // 机构名称
	ResourceUploadCount       int64                  `gorm:"column:resource_upload_count"                    json:"resourceUploadCount"`       // 上传资源数
	ResourceTypeCountPpt      int64                  `gorm:"column:resource_type_count_ppt"                  json:"resourceTypeCountPpt"`      // 资源类型数量-ppt
	ResourceTypeCountPdf      int64                  `gorm:"column:resource_type_count_pdf"                  json:"resourceTypeCountPdf"`      // 资源类型数量-pdf
	ResourceTypeCountWord     int64                  `gorm:"column:resource_type_count_word"                 json:"resourceTypeCountWord"`     // 资源类型数量-word
	ResourceTypeCountExcel    int64                  `gorm:"column:resource_type_count_excel"                json:"resourceTypeCountExcel"`    // 资源类型数量-excel
	ResourceTypeCountPic      int64                  `gorm:"column:resource_type_count_pic"                  json:"resourceTypeCountPic"`      // 资源类型数量-pic
	ResourceTypeCountVideo    int64                  `gorm:"column:resource_type_count_video"                json:"resourceTypeCountVideo"`    // 资源类型数量-video
	ResourceTypeCountAudio    int64                  `gorm:"column:resource_type_count_audio"                json:"resourceTypeCountAudio"`    // 资源类型数量-audio
	ResourceTypeCountCompress int64                  `gorm:"column:resource_type_count_compress"             json:"resourceTypeCountCompress"` // 资源类型数量-compress
	ResourceParseDocument     int64                  `gorm:"column:resource_parse_document"                  json:"resourceParseDocument"`     // 文档类资源解析总数
	ResourceParseVideo        int64                  `gorm:"column:resource_parse_video"                     json:"resourceParseVideo"`        // 音视频类资源解析总数
	OssTotalCapacity          int64                  `gorm:"column:oss_total_capacity"                       json:"ossTotalCapacity"`          // oss总使用量
	CreatedAt                 int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                         // 创建时间
	UpdatedAt                 int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                         // 更新时间
	DeletedAt                 *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                         // 删除时间
}

// TableName get sql table name.获取数据库表名
func (m *OssResourSummary) TableName() string {
	return "oss_resour_summary"
}
