package tablemodel

import "gorm.io/plugin/soft_delete"

// daily_schedule 作息时间表
type DailySchedule struct {
	ID            int64                  `gorm:"primaryKey;column:id" json:"id"`
	SchoolYear    *int64                 `gorm:"column:school_year" json:"school_year"`                            // 学年
	Grade         string                 `gorm:"column:grade" json:"grade"`                                        // 年级
	YearStartDate *int64                 `gorm:"column:year_start_date" json:"year_start_date"`                    // 学年开始日期
	YearEndDate   *int64                 `gorm:"column:year_end_date" json:"year_end_date"`                        // 学年结束日期
	OrgID         *int64                 `gorm:"column:org_id" json:"orgId"`                                       // 机构校区id
	StartTime     int64                  `gorm:"column:start_time" json:"start_time"`                              // 作息开始时间
	EndTime       int64                  `gorm:"column:end_time" json:"end_time"`                                  // 作息结束时间
	EffEndTime    int64                  `gorm:"column:eff_end_time" json:"eff_end_time"`                          // 生效结束时间
	Order         int64                  `gorm:"column:order" json:"order"`                                        // 排序字段
	CreatedAt     int64                  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`          // 创建时间
	UpdatedAt     int64                  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`          // 更新时间
	DeletedAt     *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"deletedAt"` // 删除时间id
}

// class_time_schedule 排课表
type ClassTimeSchedule struct {
	ID              int64                  `gorm:"primaryKey;column:id" json:"id"`
	OrgID           *int64                 `gorm:"column:org_id" json:"org_id"`                                      // 机构校区id
	// TeacherId       *int64                 `gorm:"column:teacher_id" json:"teacher_id"`                              // 任课老师
	PrimaryId       *int64                 `gorm:"column:primary_id" json:"primary_id"`                              // 主讲老师
	CourseId        int64                  `gorm:"column:start_time" json:"start_time"`                              // 学科id
	CourseName      string                 `gorm:"column:course_name" json:"course_name"`                            // 学科名称(冗余字段)
	ClassDate       int64                  `gorm:"column:class_date" json:"class_date"`                              // 上课日期
	LiveClassId     int64                  `gorm:"column:live_class_id" json:"live_class_id"`                        // 关联课堂id
	WhiteBoardId    int64                  `gorm:"column:white_board_id" json:"white_board_id"`                      // 关联白板id
	DailyScheduleId int64                  `gorm:"column:daily_schedule_id" json:"daily_schedule_id"`                // 关联作息时间表id
	ScheduleTitle   string                 `gorm:"column:schedule_title" json:"schedule_title"`                      // 授课名称,255长度但是实际限制40,
	CreatedAt       int64                  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`          // 创建时间
	UpdatedAt       int64                  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`          // 更新时间
	DeletedAt       *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"deletedAt"` // 删除时间
}

// schedule_group 排课关联班级表
type ScheduleGroup struct {
	ID         int64  `gorm:"primaryKey;column:id" json:"id"`
	GroupId    *int64 `gorm:"column:group_id" json:"group_id"`                         // 班级id
	GroupName  string `gorm:"column:group_name" json:"group_name"`                     // 班级名称(冗余字段)
	ScheduleId *int64 `gorm:"column:schedule_id" json:"schedule_id"`                   // 排课id
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"` // 创建时间
	UpdatedAt  int64  `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"` // 更新时间
}
