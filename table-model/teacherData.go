package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `teacher_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) NOT NULL COMMENT '学校名称',
	`summary_start_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`summary_end_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`user_id` bigint NOT NULL COMMENT '老师ID',
	`user_name` varchar(50) NOT NULL COMMENT '老师名称',
	`create_class_num` bigint NOT NULL COMMENT '创建授课数量',
	`create_realtime_count` bigint NOT NULL COMMENT '创建教研类型数量统计_实时听评课',
	`create_obs_count` bigint NOT NULL COMMENT '创建教研类型数量统计_第三方推流',
	`create_record_playback_count` bigint NOT NULL COMMENT '创建教研类型数量统计_录播听评课',
	`create_yxtboard_count` bigint NOT NULL COMMENT '创建教研类型数量统计_观摩研讨总数',
	`create_prepare_lesson_count` bigint NOT NULL COMMENT '创建备课数量',
	`join_prepare_lesson_count` bigint NOT NULL COMMENT '参与备课数',
	`join_prepare_lesson_meeting_count` bigint NOT NULL COMMENT '参与备课会议总数',
	`join_prepare_lesson_meeting_duration` bigint NOT NULL COMMENT '参与备课会议总时长',
	`join_teaching_activity_count` bigint NOT NULL COMMENT '参与教研数',
	`join_live_class_count` bigint NOT NULL COMMENT '参与课堂数',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`),
	KEY `summary_start_date` (`summary_start_date`) USING BTREE COMMENT '统计日期索引',
	KEY `summary_end_date` (`summary_end_date`) USING BTREE COMMENT '统计日期索引',
	UNIQUE KEY `org_date_indx` (`org_id`,`summary_start_date`,`summary_end_date`) USING BTREE

) ENGINE = InnoDB COMMENT = '老师汇总信息表';
*/
type TeacherSummary struct {
	ID                               int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                            int64                  `gorm:"column:org_id"                                   json:"orgId"`                            // 校区ID
	OrgName                          string                 `gorm:"column:org_name"                                 json:"orgName"`                          // 机构名称
	SummaryStartDate                 int64                  `gorm:"column:summary_start_date"                       json:"summaryStartDate"`                 // 统计开始日期；精确到天,单位毫秒
	SummaryEndDate                   int64                  `gorm:"column:summary_end_date"                         json:"summaryEndDate"`                   // 统计开始日期；精确到天,单位毫秒
	UserId                           int64                  `gorm:"column:user_id"                                  json:"userId"`                           //老师ID
	UserName                         string                 `gorm:"column:user_name"                                json:"userName"`                         //老师名称
	CreateClassNum                   int64                  `gorm:"column:create_class_num"                         json:"createClassNum"`                   //创建授课数量
	CreateRealtimeCount              int64                  `gorm:"column:create_realtime_count"                    json:"createRealtimeCount"`              //创建教研类型数量统计_实时听评课
	CreateObsCount                   int64                  `gorm:"column:create_obs_count"                         json:"createObsCount"`                   //创建教研类型数量统计_第三方推流
	CreateRecordPlaybackCount        int64                  `gorm:"column:create_record_playback_count"             json:"createRecordPlaybackCount"`        //创建教研类型数量统计_录播听评课
	CreateYxtboardCount              int64                  `gorm:"column:create_yxtboard_count"                    json:"createYxtboardCount"`              //创建教研类型数量统计_观摩研讨总数
	CreatePrepareLessonCount         int64                  `gorm:"column:create_prepare_lesson_count"              json:"createPrepareLessonCount"`         //创建备课数量
	JoinPrepareLessonCount           int64                  `gorm:"column:join_prepare_lesson_count"                json:"joinPrepareLessonCount"`           //参与备课数
	JoinPrepareLessonMeetingCount    int64                  `gorm:"column:join_prepare_lesson_meeting_count"        json:"joinPrepareLessonMeetingCount"`    //参与备课会议总数
	JoinPrepareLessonMeetingDuration int64                  `gorm:"column:join_prepare_lesson_meeting_duration"     json:"joinPrepareLessonMeetingDuration"` //参与备课会议总时长
	JoinTeachingActivityCount        int64                  `gorm:"column:join_teaching_activity_count"             json:"joinTeachingActivityCount"`        //参与教研数
	JoinLiveClassCount               int64                  `gorm:"column:join_live_class_count"                    json:"joinLiveClassCount"`               //参与课堂数
	CreatedAt                        int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                                // 创建时间
	UpdatedAt                        int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                                // 更新时间
	DeletedAt                        *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                                // 删除时间
}

func (*TeacherSummary) TableName() string {
	return "teacher_summary"
}
