package tableModel

import "gorm.io/plugin/soft_delete"

// workbench_column WorkbenchColumn workbenchColumn

const (
	liveClassColumn           = 1 //授课
	teachingActivityColumn    = 2 //教研
	personalSpaceColumn       = 3 //画布
	prepareLessonManageColumn = 4 //备课
	particularsColumn         = 5 //资料
	todoListColumn            = 6 //待办
)

type WorkbenchColumn struct {
	ID         int64                  `gorm:"primaryKey;column:id"                            json:"id"`         // 主键
	OrgID      int64                  `gorm:"column:org_id"                                   json:"-"`          // 机构id
	NameKey    string                 `gorm:"column:name_key"                                 json:"nameKey"`    // 名称key
	Name       string                 `gorm:"column:name"                                     json:"name"`       // 名称
	ColumnType int                    `gorm:"column:column_type"                              json:"columnType"` // 栏目类型 1-授课 2-教研 3-画布 4-备课 5-资料 6-待办
	Order      int                    `gorm:"column:column_order"                             json:"order"`      // 顺序
	CreatedAt  int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`          // 创建时间
	UpdatedAt  int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`          // 更新时间
	DeletedAt  *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`          // 删除时间
}

// TableName WorkbenchColumn's table name
func (*WorkbenchColumn) TableName() string {
	return "workbench_column"
}
