package tablemodel

import "gorm.io/plugin/soft_delete"

const (
	FolderGroupOwner     int = 1 // 分享文件夹所有者
	FolderGroupMember    int = 2 // 组员
	FolderGroupNonMember int = 0 // 非组员
)

// FileSystem mapped from table <file_system>
type FolderMemberGroup struct {
	ID        int64                  `gorm:"column:id;primaryKey;autoIncrement:true"         json:"id"`
	OrgID     int64                  `gorm:"column:org_id"                                   json:"orgId"`     // 校区ID
	UserId    int64                  `gorm:"column:user_id"                                  json:"userId"`    // 文件成员用户id
	FolderId  int64                  `gorm:"column:folder_id"                                json:"folderId"`  // 文件夹的id
	ShareType int                    `gorm:"column:share_type"                               json:"shareType"` // 操作类型 1 可编辑 2 可批注 3 可查看
	CreatedAt int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`         // 创建时间
	UpdatedAt int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`         // 更新时间
	DeletedAt *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`         // 删除时间
}

// TableName FileSystem's table name
func (*FolderMemberGroup) TableName() string {
	return "folder_member_group"
}
