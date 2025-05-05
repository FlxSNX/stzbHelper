package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"stzbHelper/global"
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
		if add.Error != nil {
			common.Response{Message: "创建失败", Data: add.Error.Error()}.Error(c)
		} else {
			common.Response{Message: "创建失败"}.Error(c)
		}
	}
}

func GetTaskList(c *gin.Context) {
	var taskList []model.Task

	model.Conn.Omit("user_list").Order("id DESC").Find(&taskList)
	common.Response{Data: taskList}.Success(c)
}

func DelTask(c *gin.Context) {
	id := c.Param("tid")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	action := model.Conn.Delete(&model.Task{}, idInt)

	if action.RowsAffected != 0 {
		common.Response{Message: "删除成功", Data: action.RowsAffected}.Success(c)
	} else {
		if action.Error != nil {
			common.Response{Message: "删除失败", Data: action.Error.Error()}.Error(c)
		} else {
			common.Response{Message: "删除失败"}.Error(c)
		}
	}
}

func EnableGetReport(c *gin.Context) {
	pos := c.PostForm("pos")

	posInt, err := strconv.Atoi(pos)
	if err != nil {
		common.Response{Message: "坐标格式错误"}.Error(c)
		return
	}

	global.ExVar.NeededReportPos = posInt
	global.ExVar.NeedGetReport = true
	common.Response{}.Success(c)
}

func DisableGetReport(c *gin.Context) {
	global.ExVar.NeededReportPos = 0
	global.ExVar.NeedGetReport = false
	common.Response{}.Success(c)
}

func GetReportNumByTaskId(c *gin.Context) {
	tid := c.Param("tid")

	if tid == "" {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	tidint, err := strconv.Atoi(tid)

	if err != nil {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	var task model.Task

	query := model.Conn.Last(&task, tidint)

	if query.Error == nil {
		var taskNum int64

		model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos}).Count(&taskNum)

		common.Response{Data: gin.H{
			"count": taskNum,
		}}.Success(c)
	} else {
		common.Response{Message: "获取任务失败"}.Error(c)
		return
	}
}

func StatisticsReport(c *gin.Context) {
	tid := c.Param("tid")

	if tid == "" {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	tidint, err := strconv.Atoi(tid)

	if err != nil {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	var task model.Task

	query := model.Conn.Last(&task, tidint)
	//fmt.Println(task, query.Error)
	if query.Error == nil {
		task.CompleteUserNum = 0
		for id, t := range task.UserList {
			var Num int64
			model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos, AttackName: t.Name}).Count(&Num)
			//fmt.Print(t.Name, "总战报数量:", Num, " ")
			//查询攻城次数
			var AtkNum int64
			model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos, AttackName: t.Name}).Where("garrison = ?", 0).Count(&AtkNum)
			//fmt.Print(t.Name, "主力次数:", AtkNum, " ")
			//查询拆迁次数
			var DisNum int64
			model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos, AttackName: t.Name, Garrison: 1}).Count(&DisNum)
			//fmt.Println(t.Name, "拆迁次数:", DisNum)
			//主力队伍数量
			var AtkTeamNum int64
			model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos, AttackName: t.Name}).Where("garrison = ?", 0).Group("attack_base_heroid").Count(&AtkTeamNum)
			//fmt.Print(t.Name, "主力队伍数量:", AtkTeamNum, " ")
			//拆迁队伍数量
			var DisTeamNum int64
			model.Conn.Model(model.Report{}).Where(model.Report{Wid: task.Pos, AttackName: t.Name, Garrison: 1}).Group("attack_base_heroid").Count(&DisTeamNum)
			//fmt.Println(t.Name, "拆迁队伍数量:", DisTeamNum, " ")
			task.UserList[id].AtkNum = int(AtkNum)
			task.UserList[id].DisNum = int(DisNum)
			task.UserList[id].AtkTeamNum = int(AtkTeamNum)
			task.UserList[id].DisTeamNum = int(DisTeamNum)
			if AtkNum != 0 || DisNum != 0 {
				task.CompleteUserNum++
			}
		}
		save := model.Conn.Save(&task)
		if save.RowsAffected != 0 {
			common.Response{Message: "统计考勤数据成功", Data: save.RowsAffected}.Success(c)
		} else {
			if save.Error != nil {
				common.Response{Message: "统计考勤数据失败", Data: save.Error.Error()}.Error(c)
			} else {
				common.Response{Message: "统计考勤数据失败"}.Error(c)
			}
		}
	} else {
		common.Response{Message: "获取任务失败"}.Error(c)
		return
	}
}

func GetTask(c *gin.Context) {
	tid := c.Param("tid")

	if tid == "" {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	tidint, err := strconv.Atoi(tid)

	if err != nil {
		common.Response{Message: "任务ID错误"}.Error(c)
		return
	}

	var task model.Task

	query := model.Conn.Last(&task, tidint)
	//fmt.Println(task, query.Error)
	if query.Error == nil {
		common.Response{Data: task}.Success(c)
	} else {
		common.Response{Message: "获取任务失败"}.Error(c)
		return
	}
}

func Example(c *gin.Context) {
	common.Response{Message: "This is example func"}.Success(c)
}
