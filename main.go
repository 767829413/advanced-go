package main

import (
	"fmt"
	"strings"
)

var (
	notifyStr = `
	### 课堂统计：2024-04-01
	>**系统环境** : dev
	>**系统域名** : dev-s1.plaso.cn
	>**上课时间** : 10点-11点
	>**课堂总数量** : 1000
	>**课堂实际上课总数量** : 1000
	>**课堂总时长** : 232
	>**课堂真实总时长** : 200
	>**课堂直播总时长** : 200
	>**课堂用户总数** : 800
	>**课堂用户参与总时长** : 800

	### Top 10 课堂参与用户数排名
`
)

type classInfo struct {
	Name               string // 课堂名称
	TotalClassUserNum  int    // 课堂用户总数
	TotalClassUserTime int64  // 课堂用户参与总时长
	OrgId              int64  // 机构id
	OrgName            string // 机构名称
}

func main() {
	listClassInfos := []*classInfo{
		{
			Name:               "课堂1",
			TotalClassUserNum:  100,
			TotalClassUserTime: 100,
			OrgId:              1,
			OrgName:            "机构1",
		},
		{
			Name:               "课堂2",
			TotalClassUserNum:  80,
			TotalClassUserTime: 200,
			OrgId:              1,
			OrgName:            "机构1",
		},
		{
			Name:               "课堂3",
			TotalClassUserNum:  70,
			TotalClassUserTime: 300,
			OrgId:              1,
			OrgName:            "机构1",
		},
	}
	topNum := 9
	var builder strings.Builder
	builder.WriteString(notifyStr)
	for i, v := range listClassInfos {
		builder.WriteString(fmt.Sprintf(`
	>**课节名称** : %s
	>**课节用户数** : %d
	>**课节用户时长** : %d
	>**所属机构名称** : %s
`,
			v.Name, v.TotalClassUserNum, v.TotalClassUserTime, v.OrgName))
		if i == topNum {
			break
		}
	}
	fmt.Println(builder.String())
}
