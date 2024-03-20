package tableModel

import "gorm.io/plugin/soft_delete"

type ToDoList struct {
	ID         int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID      int64                  `gorm:"column:org_id"                                   json:"orgId"`      // 校区ID
	StartTime  int64                  `gorm:"column:start_time"                               json:"startTime"`  // 开始时间
	EndTime    int64                  `gorm:"column:end_time"                                 json:"endTime"`    // 结束时间
	ToDoName   string                 `gorm:"column:to_do_name"                               json:"toDoName"`   // 待办名称
	ToDoType   int                    `gorm:"column:to_do_type"                               json:"toDoType"`   // 待办类型 待办类型 1-备课任务:class_time_schedule 2-备课会议:infi_board_meeting 3-教研活动:teaching_activity
	RelationID int64                  `gorm:"column:relation_id"                              json:"relationID"` // 关联类型 1-备课 2-教研活动场次
	PrimaryId  *int64                 `gorm:"column:primary_id"                               json:"primaryId"`  // 主讲老师id
	CourseId   int64                  `gorm:"column:course_id"                                json:"courseId"`   // 学科id
	ExtendInfo string                 `gorm:"column:extend_info"                              json:"extendInfo"` // 扩展信息
	GradeName  string                 `gorm:"column:grade_name"                               json:"gradeName"`  // 年级名称
	CreatedAt  int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`          // 创建时间
	UpdatedAt  int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`          // 更新时间
	DeletedAt  *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`          // 删除时间
}

type ToDoListUser struct {
	ID         int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	OrgID      int64 `gorm:"column:org_id"                           json:"orgId"`      // 校区ID
	UserId     int64 `gorm:"column:user_id"                          json:"userId"`     // 用户id
	ToDoStatus int   `gorm:"column:to_do_status"                     json:"toDoStatus"` // 待办状态 1-未完成 2-已完成
	ToDoID     int64 `gorm:"column:to_do_id"                         json:"toDoId"`     // 待办id
	CreatedAt  int64 `gorm:"column:created_at;autoCreateTime:milli"  json:"-"`          // 创建时间
	UpdatedAt  int64 `gorm:"column:updated_at;autoUpdateTime:milli"  json:"-"`          // 更新时间
}

// TableName ToDoList's table name
func (*ToDoList) TableName() string {
	return "to_do_list"
}

// TableName ToDoListUser's table name
func (*ToDoListUser) TableName() string {
	return "to_do_list_user"
}
