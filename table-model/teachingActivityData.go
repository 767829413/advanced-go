package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `teaching_activity_summary` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`org_id` bigint NOT NULL COMMENT '校区ID',
	`org_name` varchar(50) NOT NULL COMMENT '学校名称',
	`summary_start_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`summary_end_date` bigint NOT NULL COMMENT '统计日期；精确到天,单位毫秒',
	`teach_activity_count` bigint NOT NULL COMMENT '创建活动总数',
	`teach_activity_duration` bigint NOT NULL COMMENT '教研活动总时长（毫秒）',
	`activity_average_duration` bigint NOT NULL COMMENT '教研活动平均时长（毫秒）',
	`type_realtime_count` bigint NOT NULL COMMENT '教研类型数量统计:实时听评课',
	`type_obs_count` bigint NOT NULL COMMENT '教研类型数量统计:第三方推流',
	`type_record_playback_count` bigint NOT NULL COMMENT '教研类型数量统计:录播听评课',
	`type_yxtboard_count` bigint NOT NULL COMMENT '教研类型数量统计:观摩研讨总数',
	`max_online_num` bigint NOT NULL COMMENT '历史最大在线人数总数',
	`accumulate_online_num` bigint NOT NULL COMMENT '历史累计在线人数总数',
	`register_num` bigint NOT NULL COMMENT '历史注册人数总数',
	`cumulation_duration_count` bigint NOT NULL COMMENT '总参会时长',
	`entry_num_count` bigint NOT NULL COMMENT '总入会次数',
	`comment_count` bigint NOT NULL COMMENT '点评总数',
	`teach_duration` bigint NOT NULL COMMENT '授课时长',
	`evaluate_duration` bigint NOT NULL COMMENT '评课时长',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	KEY `summary_start_date` (`summary_start_date`) USING BTREE COMMENT '统计日期索引',
	KEY `summary_end_date` (`summary_end_date`) USING BTREE COMMENT '统计日期索引',
	UNIQUE KEY `org_date_indx` (`org_id`,`summary_start_date`,`summary_end_date`) USING BTREE

) ENGINE=InnoDB COMMENT = '教研汇总信息表';
*/
type TeachingActivitySummary struct {
	ID                      int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID                   int64                  `gorm:"column:org_id"                                   json:"orgId"`                   // 校区ID
	OrgName                 string                 `gorm:"column:org_name"                                 json:"orgName"`                 // 机构名称
	SummaryStartDate        int64                  `gorm:"column:summary_start_date"                       json:"summaryStartDate"`        // 统计开始日期；精确到天,单位毫秒
	SummaryEndDate          int64                  `gorm:"column:summary_end_date"                         json:"summaryEndDate"`          // 统计开始日期；精确到天,单位毫秒
	TeachActivityCount      int64                  `gorm:"column:teach_activity_count"                     json:"teachActivityCount"`      //创建活动总数
	TeachActivityDuration   int64                  `gorm:"column:teach_activity_duration"                  json:"teachActivityDuration"`   //教研活动总时长（毫秒）
	ActivityAverageDuration int64                  `gorm:"column:activity_average_duration"                json:"activityAverageDuration"` //教研活动平均时长（毫秒）
	TypeRealtimeCount       int64                  `gorm:"column:type_realtime_count"                      json:"typeRealtimeCount"`       //教研类型数量统计:实时听评课
	TypeObsCount            int64                  `gorm:"column:type_obs_count"                           json:"typeObsCount"`            //教研类型数量统计:第三方推流
	TypeRecordPlaybackCount int64                  `gorm:"column:type_record_playback_count"               json:"typeRecordPlaybackCount"` //教研类型数量统计:录播听评课
	TypeYxtboardCount       int64                  `gorm:"column:type_yxtboard_count"                      json:"typeYxtboardCount"`       //教研类型数量统计:观摩研讨总数
	MaxOnlineNum            int64                  `gorm:"column:max_online_num"                           json:"maxOnlineNum"`            //历史最大在线人数总数
	AccumulateOnlineNum     int64                  `gorm:"column:accumulate_online_num"                    json:"accumulateOnlineNum"`     //历史累计在线人数总数
	RegisterNum             int64                  `gorm:"column:register_num"                             json:"registerNum"`             //历史注册人数总数
	CumulationDurationCount int64                  `gorm:"column:cumulation_duration_count"                json:"cumulationDurationCount"` //总参会时长
	EntryNumCount           int64                  `gorm:"column:entry_num_count"                          json:"entryNumCount"`           //总入会次数
	CommentCount            int64                  `gorm:"column:comment_count"                            json:"commentCount"`            //点评总数
	TeachDuration           int64                  `gorm:"column:teach_duration"                           json:"teachDuration"`           //授课时长
	EvaluateDuration        int64                  `gorm:"column:evaluate_duration"                        json:"evaluateDuration"`        //评课时长
	CreatedAt               int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`                       // 创建时间
	UpdatedAt               int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`                       // 更新时间
	DeletedAt               *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`                       // 删除时间
}

func (*TeachingActivitySummary) TableName() string {
	return "teaching_activity_summary"
}
