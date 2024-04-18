package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `prepare_lesson_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) NOT NULL COMMENT '学校名称',
	`summary_start_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`summary_end_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`prepare_lesson_count` bigint NOT NULL COMMENT '创建备课总数',
	`prepare_lesson_board_count` bigint NOT NULL COMMENT '备课创建英飞画布数',
	`prepare_lesson_board_meeting_count` bigint NOT NULL COMMENT '备课开启英飞画布会议数',
	`prepare_lesson_meeting_duration` bigint NOT NULL COMMENT '备课英飞画布会议总时长:毫秒',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `summary_start_date` (`summary_start_date`) USING BTREE COMMENT '统计日期索引',
	KEY `summary_end_date` (`summary_end_date`) USING BTREE COMMENT '统计日期索引',
	UNIQUE KEY `org_date_indx` (`org_id`,`summary_start_date`,`summary_end_date`) USING BTREE

) ENGINE=InnoDB COMMENT = '备课汇总信息表';
*/
type PrepareLessonSummary struct {
	ID                             int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                          int64                  `gorm:"column:org_id"                                   json:"orgId"`                          // 校区ID
	OrgName                        string                 `gorm:"column:org_name"                                 json:"orgName"`                        // 机构名称
	SummaryStartDate               int64                  `gorm:"column:summary_start_date"                       json:"summaryStartDate"`               // 统计开始日期；精确到天,单位毫秒
	SummaryEndDate                 int64                  `gorm:"column:summary_end_date"                         json:"summaryEndDate"`                 // 统计开始日期；精确到天,单位毫秒
	PrepareLessonCount             int64                  `gorm:"column:prepare_lesson_count"                     json:"prepareLessonCount"`             // 创建备课总数
	PrepareLessonBoardCount        int64                  `gorm:"column:prepare_lesson_board_count"               json:"prepareLessonBoardCount"`        // 备课创建英飞画布数
	PrepareLessonBoardMeetingCount int64                  `gorm:"column:prepare_lesson_board_meeting_count"       json:"prepareLessonBoardMeetingCount"` // 备课开启英飞画布会议数
	PrepareLessonMeetingDuration   int64                  `gorm:"column:prepare_lesson_meeting_duration"          json:"prepareLessonMeetingDuration"`   // 备课英飞画布会议总时长_毫秒
	CreatedAt                      int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                              // 创建时间
	UpdatedAt                      int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                              // 更新时间
	DeletedAt                      *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                              // 删除时间
}

func (p *PrepareLessonSummary) TableName() string {
	return "prepare_lesson_summary"
}
