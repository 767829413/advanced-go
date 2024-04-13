package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `prepare_lesson_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学校名称',
	`course_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学科名称',
	`prepare_lesson_count` bigint NOT NULL COMMENT '创建备课总数',
	`prepare_lesson_board_count` bigint NOT NULL COMMENT '备课创建白板数',
	`prepare_lesson_board_meeting_count` bigint NOT NULL COMMENT '备课开启白板会议数',
	`prepare_lesson_meeting_duration` bigint NOT NULL COMMENT '备课会议总时长',
	`prepare_lesson_pv` bigint NOT NULL COMMENT '备课PV',
	`prepare_lesson_uv` bigint NOT NULL COMMENT '备课UV',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `org_id_index` (`org_id`)

) ENGINE=InnoDB COMMENT = '工作台栏目表';
*/
type PrepareLessonSummary struct {
	ID                             int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                          int64                  `gorm:"column:org_id"                                   json:"orgId"`                          // 校区ID
	OrgName                        string                 `gorm:"column:org_name"                                 json:"orgName"`                        // 机构名称
	CourseName                     string                 `gorm:"column:course_name"                              json:"courseName"`                     // 学科名称
	PrepareLessonCount             int64                  `gorm:"column:prepare_lesson_count"                     json:"prepareLessonCount"`             // 创建备课总数
	PrepareLessonBoardCount        int64                  `gorm:"column:prepare_lesson_board_count"               json:"prepareLessonBoardCount"`        // 备课创建白板数
	PrepareLessonBoardMeetingCount int64                  `gorm:"column:prepare_lesson_board_meeting_count"       json:"prepareLessonBoardMeetingCount"` // 备课开启白板会议数
	PrepareLessonMeetingDuration   int64                  `gorm:"column:prepare_lesson_meeting_duration"          json:"prepareLessonMeetingDuration"`   // 备课会议总时长:毫秒
	PrepareLessonPV                int64                  `gorm:"column:prepare_lesson_pv"                        json:"prepareLessonPV"`                // 备课PV
	PrepareLessonUV                int64                  `gorm:"column:prepare_lesson_uv"                        json:"prepareLessonUV"`                // 备课UV
	CreatedAt                      int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                              // 创建时间
	UpdatedAt                      int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                              // 更新时间
	DeletedAt                      *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                              // 删除时间
}

// TableName get sql table name.获取数据库表名
func (m *PrepareLessonSummary) TableName() string {
	return "prepare_lesson_summary"
}
