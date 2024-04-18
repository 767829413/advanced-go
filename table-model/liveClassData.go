package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `live_class_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) NOT NULL COMMENT '学校名称',
	`summary_start_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`summary_end_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`total_class_num` bigint NOT NULL COMMENT '课堂总数',
	`total_class_real_num` bigint NOT NULL COMMENT '课堂实际总数',
	`total_class_time` bigint NOT NULL COMMENT '课堂总时长',
	`total_class_real_time` bigint NOT NULL COMMENT '课堂真实总时长',
	`total_class_live_num` bigint NOT NULL COMMENT '直播课堂总数',
	`total_class_live_time` bigint NOT NULL COMMENT '课堂直播总时长',
	`total_class_user_num` bigint NOT NULL COMMENT '课堂用户总数',
	`total_class_real_user_num` bigint NOT NULL COMMENT '课堂真实用户总数',
	`total_class_interactive_user_num` bigint NOT NULL COMMENT '课堂互动用户总数',
	`total_class_link_user_num` bigint NOT NULL COMMENT '链接访客总数',
	`total_class_user_time` bigint NOT NULL COMMENT '课堂用户参与总时长',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `summary_start_date` (`summary_start_date`) USING BTREE COMMENT '统计日期索引',
	KEY `summary_end_date` (`summary_end_date`) USING BTREE COMMENT '统计日期索引',
	UNIQUE KEY `org_date_indx` (`org_id`,`summary_start_date`,`summary_end_date`) USING BTREE

) ENGINE=InnoDB COMMENT = '授课汇总记录表';
*/
type LiveClassSummary struct {
	ID                           int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                        int64                  `gorm:"column:org_id"                                   json:"orgId"`                        // 校区ID
	OrgName                      string                 `gorm:"column:org_name"                                 json:"orgName"`                      // 机构名称
	SummaryStartDate             int64                  `gorm:"column:summary_start_date"                       json:"summaryStartDate"`             // 统计开始日期；精确到天,单位毫秒
	SummaryEndDate               int64                  `gorm:"column:summary_end_date"                         json:"summaryEndDate"`               // 统计开始日期；精确到天,单位毫秒
	TotalClassNum                int64                  `gorm:"column:total_class_num"                          json:"totalClassNum"`                // 课堂总数
	TotalClassRealNum            int64                  `gorm:"column:total_class_real_num"                     json:"totalClassRealNum"`            // 课堂实际总数
	TotalClassTime               int64                  `gorm:"column:total_class_time"                         json:"totalClassTime"`               // 课堂总时长
	TotalClassRealTime           int64                  `gorm:"column:total_class_real_time"                    json:"totalClassRealTime"`           // 课堂真实总时长
	TotalClassLiveNum            int64                  `gorm:"column:total_class_live_num"                     json:"totalClassLiveNum"`            // 直播课堂总数
	TotalClassLiveTime           int64                  `gorm:"column:total_class_live_time"                    json:"totalClassLiveTime"`           // 课堂直播总时长
	TotalClassUserNum            int64                  `gorm:"column:total_class_user_num"                     json:"totalClassUserNum"`            // 课堂用户总数
	TotalClassRealUserNum        int64                  `gorm:"column:total_class_real_user_num"                json:"totalClassRealUserNum"`        // 课堂真实用户总数
	TotalClassInteractiveUserNum int64                  `gorm:"column:total_class_interactive_user_num"         json:"totalClassInteractiveUserNum"` // 课堂互动用户总数
	TotalClassLinkUserNum        int64                  `gorm:"column:total_class_link_user_num"                json:"totalClassLinkUserNum"`        // 链接访客总数
	TotalClassUserTime           int64                  `gorm:"column:total_class_user_time"                    json:"totalClassUserTime"`           // 课堂用户参与总时长
	CreatedAt                    int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                            // 创建时间
	UpdatedAt                    int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                            // 更新时间
	DeletedAt                    *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                            // 删除时间
}

func (*LiveClassSummary) TableName() string {
	return "live_class_summary"
}
