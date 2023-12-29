package tablemodel

import "gorm.io/plugin/soft_delete"

// daily_schedule 作息时间表
type DailySchedule struct {
	ID           int64                  `gorm:"primaryKey;column:id" json:"id"`
	SchoolYear   *int64                 `gorm:"column:school_year" json:"schoolYear"`                             // 学年
	GradeId      string                 `gorm:"column:grade_id" json:"gradeId"`                                   // 年级虚拟ID,形式是：学段_入学年份
	YearEndDate  *int64                 `gorm:"column:year_end_date" json:"yearEndDate"`                          // 学年结束日期
	OrgID        *int64                 `gorm:"column:org_id" json:"orgId"`                                       // 机构校区id
	StartTime    int64                  `gorm:"column:start_time" json:"startTime"`                               // 作息开始时间
	EndTime      int64                  `gorm:"column:end_time" json:"endTime"`                                   // 作息结束时间
	EffStartTime int64                  `gorm:"column:eff_start_time" json:"effStartTime"`                        // 生效开始时间
	EffEndTime   int64                  `gorm:"column:eff_end_time" json:"effEndTime"`                            // 生效结束时间
	Order        int64                  `gorm:"column:order" json:"order"`                                        // 排序字段,利用作息开始时间结束时间排序
	CreatedAt    int64                  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`          // 创建时间
	UpdatedAt    int64                  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`          // 更新时间
	DeletedAt    *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"deletedAt"` // 删除时间id
}

// class_time_schedule 排课表
type ClassTimeSchedule struct {
	ID              int64                  `gorm:"primaryKey;column:id" json:"id"`
	OrgID           *int64                 `gorm:"column:org_id" json:"orgID"`                                       // 机构校区id
	PrimaryId       *int64                 `gorm:"column:primary_id" json:"primaryId"`                               // 主讲老师id
	PrimaryName     *int64                 `gorm:"column:primary_name" json:"primaryName"`                           // 主讲老师名称
	CourseId        int64                  `gorm:"column:course_Id" json:"courseId"`                                 // 学科id
	CourseName      string                 `gorm:"column:course_name" json:"courseName"`                             // 学科名称(冗余字段)
	ClassDate       int64                  `gorm:"column:class_date" json:"classDate"`                               // 上课日期
	LiveClassId     int64                  `gorm:"column:live_class_id" json:"liveClassId"`                          // 关联课堂id
	PrepareLessonId string                 `gorm:"column:prepare_lesson_id" json:"prepareLessonId"`                  // 关联备课id
	DailyScheduleId int64                  `gorm:"column:daily_schedule_id" json:"dailyScheduleId"`                  // 关联作息时间表id
	ScheduleTitle   string                 `gorm:"column:schedule_title" json:"scheduleTitle"`                       // 授课名称,255长度但是实际限制40,
	LastStartTime   int64                  `gorm:"column:last_start_time" json:"lastStartTime"`                      // 最后作息开始时间
	LastEndTime     int64                  `gorm:"column:last_end_time" json:"lastEndTime"`                          // 最后作息结束时间
	CreatedAt       int64                  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`          // 创建时间
	UpdatedAt       int64                  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`          // 更新时间
	DeletedAt       *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"deletedAt"` // 删除时间
}

// schedule_group 排课关联班级表
type ScheduleGroup struct {
	ID         int64  `gorm:"primaryKey;column:id" json:"id"`
	GroupId    *int64 `gorm:"column:group_id" json:"groupId"`                          // 班级id
	GroupName  string `gorm:"column:group_name" json:"groupName"`                      // 班级名称(冗余字段)
	ScheduleId *int64 `gorm:"column:schedule_id" json:"scheduleId"`                    // 排课id
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"` // 创建时间
	UpdatedAt  int64  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"` // 更新时间
}
