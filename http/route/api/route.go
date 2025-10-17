package api

import (
	"github.com/gin-gonic/gin"
	"stzbHelper/http/handle/api"
)

func Register(r *gin.RouterGroup) {
	// 获取同盟成员列表
	r.Any("getTeamUser", api.GetTeamUser)
	// 获取同盟成员分组列表
	r.Any("getTeamGroup", api.GetTeamGroup)
	// 获取任务列表
	r.Any("getTaskList", api.GetTaskList)
	// 获取任务详情
	r.Any("getTask/:tid", api.GetTask)
	// 创建任务
	r.POST("createTask", api.CreateTask)
	// 删除任务
	r.Any("deleteTask/:tid", api.DelTask)
	// 开启获取战报
	r.POST("enable/getReport", api.EnableGetReport)
	// 关闭获取战报
	r.Any("disable/getReport", api.DisableGetReport)
	// 获取战报数据
	r.Any("getReportNumByTaskId/:tid", api.GetReportNumByTaskId)
	// 统计考勤数据
	r.Any("statisticsReport/:tid", api.StatisticsReport)
	r.Any("getGroupWu", api.GetGroupWu)
	// 删除任务战报
	r.Any("deleteTaskReport/:tid", api.DelTaskReport)
	r.GET("stzb/report/list", api.ReportList)
	r.GET("stzb/player/team/get", api.GetPlayerTeam)
	// 开启获取战报详情
	r.Any("enable/getBattleReport", api.EnableGetBattleData)
	// 关闭获取战报详情
	r.Any("disable/getBattleReport", api.DisableGetBattleData)
}
