package tableModel

import "gorm.io/plugin/soft_delete"

/*
CREATE TABLE `org_portal_categorier_article` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`unique_id` varchar(255) NOT NULL COMMENT '唯一uuid',
	`org_id` bigint NOT NULL COMMENT '机构id',
	`column_id` bigint NOT NULL COMMENT '对应栏目id',
	`categorie_id` bigint NOT NULL COMMENT '对应栏目分类id',
	`title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
	`summaries` varchar(500) NOT DEFAULT '' NULL COMMENT '摘要',
	`content` longtext NOT NULL COMMENT '内容',
	`view_times` bigint DEFAULT '0' COMMENT '浏览次数',
	`created_by` bigint DEFAULT NULL COMMENT '创建人',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `unique_id_indx` (`unique_id`) USING BTREE,
	KEY `categorie_id` (`categorie_id`) USING BTREE COMMENT '对应栏目分类id索引',
	KEY `column_id` (`column_id`) USING BTREE COMMENT '对应栏目id索引',

) ENGINE=InnoDB COMMENT='机构分类文章关联表';
*/
type OrgPortalCategorierArticle struct {
	ID          int64                  `gorm:"primaryKey;column:id"                            json:"id"`        //主键
	UniqueId    string                 `gorm:"column:unique_id"                                json:"uniqueId"`  //唯一id
	OrgID       int64                  `gorm:"column:org_id"                                   json:"-"`         //机构id
	CategorieId int64                  `gorm:"column:categorie_id"                             json:"-"`         //对应栏目分类id
	ColumnId    int64                  `gorm:"column:column_id"                                json:"-"`         //对应栏目id
	Title       string                 `gorm:"column:title"                                    json:"title"`     //标题
	Summaries   string                 `gorm:"column:summaries"                                json:"summaries"` //摘要
	Content     string                 `gorm:"column:content"                                  json:"content"`   //内容
	ViewTimes   int64                  `gorm:"column:view_times"                               json:"viewTimes"` //浏览次数
	CreatedBy   int64                  `gorm:"column:created_by"                               json:"-"`         //创建人
	CreatedAt   int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`         //创建时间
	UpdatedAt   int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`         //更新时间
	DeletedAt   *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`         //删除时间
}

func (m *OrgPortalCategorierArticle) TableName() string {
	return "org_portal_categorier_article"
}

/*
CREATE TABLE `org_portal_home_page_config` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`unique_id` varchar(255) NOT NULL COMMENT '唯一uuid',
	`org_id` bigint NOT NULL COMMENT '机构id',
	`column_id` bigint NOT NULL COMMENT '对应栏目id',
	`allow_display` int DEFAULT '1' COMMENT '默认展示,1 展示 2 不展示',
	`name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
	`synopsis` varchar(500) NOT NULL DEFAULT '' NULL COMMENT '简介',
	`logo` varchar(255) NOT NULL DEFAULT '' COMMENT 'logo图',
	`backgrounds` varchar(255) NOT NULL DEFAULT '' COMMENT '背景图',
	`created_at` bigint NOT NULL COMMENT '创建时间',
	`updated_at` bigint NOT NULL COMMENT '更新时间',
	`deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `unique_id_indx` (`unique_id`) USING BTREE,
	KEY `categorie_id` (`categorie_id`) USING BTREE COMMENT '对应栏目分类id索引',
	KEY `column_id` (`column_id`) USING BTREE COMMENT '对应栏目id索引',

) ENGINE=InnoDB COMMENT='机构主页栏目配置表';
*/
type OrgPortalHomePageConfig struct {
	ID           int64                  `gorm:"primaryKey;column:id"                            json:"id"`           // 主键
	UniqueId     string                 `gorm:"column:unique_id"                                json:"uniqueId"`     //唯一id
	OrgID        int64                  `gorm:"column:org_id"                                   json:"-"`            //机构id
	CategorieId  int64                  `gorm:"column:categorie_id"                             json:"-"`            //对应栏目分类id
	ColumnId     int64                  `gorm:"column:column_id"                                json:"-"`            //对应栏目id
	AllowDisplay int                    `gorm:"column:allow_display"                            json:"allowDisplay"` //默认展示,1 展示 2 不展示
	Name         string                 `gorm:"column:name"                                     json:"name"`         //名称
	Synopsis     string                 `gorm:"column:synopsis"                                 json:"synopsis"`     //简介
	Logo         string                 `gorm:"column:logo"                                     json:"logo"`         //logo图
	Backgrounds  string                 `gorm:"column:backgrounds"                              json:"backgrounds"`  //背景图
	CreatedBy    int64                  `gorm:"column:created_by"                               json:"-"`            //创建人
	CreatedAt    int64                  `gorm:"column:created_at;autoCreateTime:milli"          json:"-"`            //创建时间
	UpdatedAt    int64                  `gorm:"column:updated_at;autoUpdateTime:milli"          json:"-"`            //更新时间
	DeletedAt    *soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:null" json:"-"`            //删除时间
}

func (m *OrgPortalHomePageConfig) TableName() string {
	return "org_portal_home_page_config"
}
