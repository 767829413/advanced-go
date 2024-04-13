package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `live_class_summary` (

	    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	    `org_id` bigint NOT NULL COMMENT '校区ID',
	    `org_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学校名称',
		`live_class_type` int NOT NULL DEFAULT '1' COMMENT '课堂类型：0 普通课堂-融课； 1 教研说课评课-融课 2 录课评课-融课，已经弃用 3 集体会议-融课; 4实时听课评课-融课与网页； 5第三方推流-网页； 6录播听评课-网页  7 观摩研讨 8 画布授课',
		`total_class_num` bigint NOT NULL COMMENT '课堂总数',
		`total_class_real_num` bigint NOT NULL COMMENT '课堂实际总数',
		`total_class_time` bigint NOT NULL COMMENT '课堂总时长',
		`total_class_real_time` bigint NOT NULL COMMENT '课堂真实总时长',
		`total_class_live_time` bigint NOT NULL COMMENT '课堂直播总时长',
		`total_class_user_num` bigint NOT NULL COMMENT '课堂用户总数',
		`total_class_real_user_num` bigint NOT NULL COMMENT '课堂真实用户总数',
		`total_class_user_time` bigint NOT NULL COMMENT '课堂用户参与总时长',
		`created_at` bigint NOT NULL COMMENT '创建时间',
		`updated_at` bigint NOT NULL COMMENT '更新时间',
		`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
		PRIMARY KEY (`id`) USING BTREE,
		KEY `org_id_index` (`org_id`)

) ENGINE=InnoDB COMMENT = '工作台栏目表';
*/
type LiveClassSummary struct {
	ID                    int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                 int64                  `gorm:"column:org_id"                                   json:"orgId"`                 // 校区ID
	OrgName               string                 `gorm:"column:org_name"                                 json:"orgName"`               // 机构名称
	LiveClassType         int                    `gorm:"column:live_class_type"                          json:"liveClassType"`         // 课堂类型
	TotalClassNum         int64                  `gorm:"column:total_class_num"                          json:"totalClassNum"`         // 课堂总数
	TotalClassRealNum     int64                  `gorm:"column:total_class_real_num"                     json:"totalClassRealNum"`     // 课堂实际总数
	TotalClassTime        int64                  `gorm:"column:total_class_time"                         json:"totalClassTime"`        // 课堂总时长
	TotalClassRealTime    int64                  `gorm:"column:total_class_real_time"                    json:"totalClassRealTime"`    // 课堂真实总时长
	TotalClassLiveTime    int64                  `gorm:"column:total_class_live_time"                    json:"totalClassLiveTime"`    // 课堂直播总时长
	TotalClassUserNum     int64                  `gorm:"column:total_class_user_num"                     json:"totalClassUserNum"`     // 课堂用户总数
	TotalClassRealUserNum int64                  `gorm:"column:total_class_real_user_num"                json:"totalClassRealUserNum"` // 课堂真实用户总数
	TotalClassUserTime    int64                  `gorm:"column:total_class_user_time"                    json:"totalClassUserTime"`    // 课堂用户参与总时长
	CreatedAt             int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                     // 创建时间
	UpdatedAt             int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                     // 更新时间
	DeletedAt             *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                     // 删除时间
}

// TableName FileSystem's table name
func (*LiveClassSummary) TableName() string {
	return "live_class_summary"
}
