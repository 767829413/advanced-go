package main

import (
	"fmt"
)

func main() {
	// 假设我们想要获取 ID 为 5 的项目的路径
	path := getPath(5)
	fmt.Println("Path:", path)
}

type FileSystemItem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	FileSystemType int    `json:"file_system_type"`
	FileUsage      int    `json:"file_usage"`
	ParentID       int    `json:"parent_id"`
	CreatorID      int    `json:"creator_id"`
	RawfileFileid  string `json:"rawfile_fileid"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	DeletedAt      int64  `json:"deleted_at"`
	Cover          string `json:"cover"`
	RelationType   int    `json:"relation_type"`
}

// 假设我们有一个函数来获取所有 FileSystemItem
func getAllFileSystemItems() map[int]*FileSystemItem {
	// 这个函数应该返回一个 map，其中 key 是 ID，value 是 FileSystemItem 指针
	// 实现细节取决于你的数据存储方式

	return nil
}

func getPath(itemID int) []string {
	items := getAllFileSystemItems()
	var path []string

	var buildPath func(int)
	buildPath = func(id int) {
		if item, ok := items[id]; ok {
			path = append(path, item.Name)
			if item.ParentID != 0 {
				buildPath(item.ParentID)
			}
		}
	}

	buildPath(itemID)

	// 反转路径
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - 1 - i
		path[i], path[j] = path[j], path[i]
	}

	return path
}
