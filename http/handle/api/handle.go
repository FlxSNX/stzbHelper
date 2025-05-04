package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"stzbHelper/http/common"
	"stzbHelper/model"
)

func GetTeamUser(c *gin.Context) {
	var teamUsers []model.TeamUser
	query := model.Conn
	group := c.Query("group")
	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	query.Find(&teamUsers)
	common.Response{Data: teamUsers}.Success(c)
}

func GetTeamGroup(c *gin.Context) {
	var groups []string
	model.Conn.Model(&model.TeamUser{}).Select("group").Distinct("group").Pluck("group", &groups)
	common.Response{Data: groups}.Success(c)
}

func CreateTask(c *gin.Context) {
	//group := c.PostForm("group")
	//if group == "" {
	//	common.Response{}.Error(c)
	//}

	// 获取普通字段
	taskName := c.PostForm("taskname")
	taskTime := c.PostForm("tasktime")
	// 获取数组字段
	targetGroup := c.PostFormArray("targetgroup")
	taskPos := c.PostFormArray("taskpos")

	taskPosFormat := model.ToTaskPos(taskPos)
	if taskPosFormat == 0 {
		common.Response{Message: "任务坐标格式错误"}.Error(c)
		return
	}

	taskTimeFormat, err := strconv.Atoi(taskTime)
	if err != nil {
		common.Response{Message: "任务时间格式错误"}.Error(c)
		return
	}

	var users []model.TeamUser
	model.Conn.Where("`group` IN ?", targetGroup).Find(&users)
	taskUserList := model.TeamUserListToTaskUserList(users)
	//fmt.Println(test)
	//common.Response{Data: taskUserList}.Success(c)
	//return

	if len(users) <= 0 {
		common.Response{Message: "创建出错:目标人数为0"}.Error(c)
		return
	}

	task := model.Task{
		Status:          0,
		Name:            taskName,
		Time:            taskTimeFormat,
		Pos:             taskPosFormat,
		Target:          targetGroup,
		TargetUserNum:   len(users),
		CompleteUserNum: 0,
		UserList:        taskUserList,
	}

	add := model.Conn.Create(&task)

	if add.RowsAffected != 0 {
		common.Response{Message: "创建成功", Data: add.RowsAffected}.Success(c)
	} else {
		common.Response{Message: "创建失败", Data: add.Error.Error()}.Error(c)
	}
}

func GetTaskList(c *gin.Context) {
	var taskList []model.Task

	model.Conn.Omit("user_list").Find(&taskList)
	common.Response{Data: taskList}.Success(c)
}

func Example(c *gin.Context) {
	common.Response{Message: "This is example func"}.Success(c)
}
