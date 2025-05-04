package model

import (
	"fmt"
	"strconv"
)

type Task struct {
	Id              int                  `json:"id" gorm:"column:id"`
	Status          int                  `json:"status" gorm:"column:status"`
	Name            string               `json:"name" gorm:"column:name"`
	Time            int                  `json:"time" gorm:"column:time"`
	Pos             int                  `json:"pos" gorm:"column:pos"`
	Target          []string             `json:"target" gorm:"column:target;serializer:json"`
	TargetUserNum   int                  `json:"target_user_num" gorm:"column:target_user_num"`
	CompleteUserNum int                  `json:"complete_user_num" gorm:"column:complete_user_num"`
	UserList        map[int]TaskUserList `json:"user_list,omitempty" gorm:"column:user_list;serializer:json"`
}

type TaskUserList struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	AtkNum     int    `json:"atk_num"`      //主力次数
	DisNum     int    `json:"dis_num"`      //拆迁次数
	AtkTeamNum int    `json:"atk_team_num"` //主力数量
	DisTeamNum int    `json:"dis_team_num"` //拆迁数量
}

func (Task) TableName() string {
	return "task"
}

func TeamUserListToTaskUserList(data []TeamUser) map[int]TaskUserList {
	taskUserList := map[int]TaskUserList{}
	for _, user := range data {
		taskUserList[user.Id] = TaskUserList{
			Id:         user.Id,
			Name:       user.Name,
			Group:      user.Group,
			AtkNum:     0,
			DisNum:     0,
			AtkTeamNum: 0,
			DisTeamNum: 0,
		}
	}
	return taskUserList
}

func ToTaskPos(pos []string) int {
	if len(pos) != 2 {
		return 0
	}
	// 转换第一个部分
	part1, err := strconv.Atoi(pos[0])
	if err != nil {
		return 0
	}
	// 转换第二个部分并确保是4位数
	part2, err := strconv.Atoi(pos[1])
	if err != nil {
		return 0
	}
	// 格式化为4位数，不足补0
	part2Str := fmt.Sprintf("%04d", part2)
	// 拼接两部分并转换为整数
	resultStr := fmt.Sprintf("%d%s", part1, part2Str)
	result, err := strconv.Atoi(resultStr)
	if err != nil {
		return 0
	}
	return result
}
