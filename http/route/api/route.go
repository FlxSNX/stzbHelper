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
	r.Any("getTask/:tid", api.Example)
	// 创建任务
	r.POST("createTask", api.CreateTask)
	// 删除任务
	r.POST("deleteTask", api.Example)
}
