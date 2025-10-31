package api

import (
	"fmt"
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
	taskName := c.PostForm("taskname")
	taskTime := c.PostForm("tasktime")
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

func EnableGetBattleData(c *gin.Context) {
	global.ExVar.NeedGetBattleData = true
	common.Response{}.Success(c)
}

func DisableGetBattleData(c *gin.Context) {
	global.ExVar.NeedGetBattleData = false
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

func GetGroupWu(c *gin.Context) {
	type GroupWuStats struct {
		Group       string `json:"group"`
		MemberCount int    `json:"member_count"`
		TotalWu     int    `json:"total_wu"`
		AverageWu   int    `json:"average_wu"`
		ZeroWuCount int    `json:"zero_wu_count"`
	}

	var stats []GroupWuStats

	subQuery := model.Conn.Model(&model.TeamUser{}).
		Select("`group`, COUNT(*) as zero_wu_count").
		Where("wu = 0").
		Group("`group`")

	err := model.Conn.Model(&model.TeamUser{}).
		Select("`team_user`.`group`, SUM(wu) as total_wu, ROUND(AVG(wu)) as average_wu, IFNULL(sub.zero_wu_count, 0) as zero_wu_count, COUNT(*) as member_count").
		Joins("LEFT JOIN (?) as sub ON sub.`group` = `team_user`.`group`", subQuery).
		Group("`team_user`.`group`").
		Order("total_wu DESC").
		Scan(&stats).Error

	if err != nil {
		common.Response{Message: "查询失败: " + err.Error()}.Error(c)
		return
	}

	common.Response{Data: stats}.Success(c)
}

func DelTaskReport(c *gin.Context) {

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
		action := model.Conn.Delete(&model.Report{}, "wid = ?", task.Pos)

		if action.RowsAffected != 0 {
			common.Response{Message: "清理战报成功", Data: action.RowsAffected}.Success(c)
		} else {
			//fmt.Println(action.Error, task.Pos)
			if action.Error != nil {
				common.Response{Message: "删除失败", Data: action.Error.Error()}.Error(c)
				return
			} else {
				common.Response{Message: "清理失败,可能战报已清理"}.Error(c)
				return
			}
		}
		//common.Response{Message: "清理战报成功"}.Success(c)
	} else {
		common.Response{Message: "清理任务战报失败"}.Error(c)
		return
	}
}

func ReportList(c *gin.Context) {
	var reportList []model.BattleReport
	nextid := c.Query("nextid")
	atkname := c.Query("atkname")
	atkunionname := c.Query("atkunionname")
	atkhp := c.Query("atkhp")
	atklevel := c.Query("atklevel")
	atkstar := c.Query("atkstar")
	stype := c.Query("type")
	nonpc := c.Query("nonpc")
	//no0army := c.Query("no0army")
	query := model.Conn.Limit(30).Order("time DESC")
	if nextid != "" {
		nexidInt, err := strconv.Atoi(nextid)
		if err != nil {
			return
		}
		if nexidInt > 0 {
			query.Where("id < ?", nexidInt)
		}
	} else {
		common.Response{Message: "参数错误"}.Error(c)
	}

	if stype == "1" || stype == "" {
		if atkname != "" {
			query.Where("attack_name LIKE ? OR defend_name LIKE ?", "%"+atkname+"%", "%"+atkname+"%")
		}

		if atkunionname != "" {
			query.Where("attack_union_name LIKE ? OR defend_union_name LIKE ?", "%"+atkunionname+"%", "%"+atkunionname+"%")
		}

		if atkhp != "" {
			query.Where("attack_hp >= ? OR defend_hp >= ?", atkhp, atkhp)
		}

		if atklevel != "" {
			query.Where("(attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?) OR (defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?)", atklevel, atklevel, atklevel, atklevel, atklevel, atklevel)
		}

		if atkstar != "" {
			query.Where("attack_total_star >= ? OR defend_total_star >= ?", atkstar, atkstar)
		}
	} else if stype == "2" {
		if atkname != "" {
			query.Where("attack_name LIKE ?", "%"+atkname+"%")
		}

		if atkunionname != "" {
			query.Where("attack_union_name LIKE", "%"+atkunionname+"%")
		}

		if atkhp != "" {
			query.Where("attack_hp >= ?", atkhp)
		}

		if atklevel != "" {
			query.Where("(attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?)", atklevel, atklevel, atklevel)
		}

		if atkstar != "" {
			query.Where("attack_total_star >= ?", atkstar)
		}
	} else if stype == "3" {
		if atkname != "" {
			query.Where("defend_name LIKE ?", "%"+atkname+"%")
		}

		if atkunionname != "" {
			query.Where("defend_union_name LIKE", "%"+atkunionname+"%")
		}

		if atkhp != "" {
			query.Where("defend_hp >= ?", atkhp)
		}

		if atklevel != "" {
			query.Where("(defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?)", atklevel, atklevel, atklevel)
		}

		if atkstar != "" {
			query.Where("defend_total_star >= ?", atkstar)
		}
	} else if stype == "4" {
		if atkname != "" {
			query.Where("attack_name LIKE ? OR defend_name LIKE ?", "%"+atkname+"%", "%"+atkname+"%")
		}

		if atkunionname != "" {
			query.Where("attack_union_name LIKE ? OR defend_union_name LIKE ?", "%"+atkunionname+"%", "%"+atkunionname+"%")
		}

		if atkhp != "" {
			query.Where("attack_hp >= ? AND defend_hp >= ?", atkhp, atkhp)
		}

		if atklevel != "" {
			query.Where("(attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?) AND (defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?)", atklevel, atklevel, atklevel, atklevel, atklevel, atklevel)
		}

		if atkstar != "" {
			query.Where("attack_total_star >= ? AND defend_total_star >= ?", atkstar, atkstar)
		}
	}

	if nonpc == "1" {
		query.Where("npc = 0")
	}

	//query.Where("npc = 0")

	query.Find(&reportList)

	var count int64

	model.Conn.Model(&model.BattleReport{}).Count(&count)

	common.Response{Data: gin.H{
		"report": reportList,
		"total":  count,
	}}.Success(c)
	//fmt.Println(reportList)
	//c.JSON(200, reportList)

}

func GetPlayerTeam(c *gin.Context) {
	name := c.Query("atkname")
	uname := c.Query("atkunionname")
	if name == "" && uname == "" {
		name = ""
	}
	var results []struct {
		PlayerName   string `json:"player_name" gorm:"player_name"`
		BattleID     int    `json:"battle_id" gorm:"battle_id"`
		Hero1ID      int    `json:"hero1_id" gorm:"hero1_id"`
		Hero2ID      int    `json:"hero2_id" gorm:"hero2_id"`
		Hero3ID      int    `json:"hero3_id" gorm:"hero3_id"`
		Hero1Level   int    `json:"hero1_level" gorm:"hero1_level"`
		Hero2Level   int    `json:"hero2_level" gorm:"hero2_level"`
		Hero3Level   int    `json:"hero3_level" gorm:"hero3_level"`
		Hero1Star    int    `json:"hero1_star" gorm:"hero1_star"`
		Hero2Star    int    `json:"hero2_star" gorm:"hero2_star"`
		Hero3Star    int    `json:"hero3_star" gorm:"hero3_star"`
		TotalStar    int    `json:"total_star" gorm:"total_star"`
		Hp           int    `json:"hp" gorm:"hp"`
		AllSkillInfo string `json:"all_skill_info" gorm:"all_skill_info"`
		Role         string `json:"role" gorm:"role"`
		Time         int    `json:"time" gorm:"time"`
		Gear         string `json:"gear" gorm:"gear"`
		HeroType     string `json:"hero_type" gorm:"hero_type"`
	}

	query := `WITH ranked_data AS (
		SELECT 
			attack_name AS player_name,
			attack_hero1_id AS hero1_id,
			attack_hero2_id AS hero2_id,
			attack_hero3_id AS hero3_id,
			attack_hero1_level AS hero1_level,
			attack_hero2_level AS hero2_level,
			attack_hero3_level AS hero3_level,
			attack_hero1_star AS hero1_star,
			attack_hero2_star AS hero2_star,
			attack_hero3_star AS hero3_star,
			attack_total_star AS total_star,
			attack_hp AS hp,
			attacker_gear_info AS gear,
			attack_hero_type AS hero_type,
			time,
			all_skill_info,
			battle_id,
			'attack' AS role,
			ROW_NUMBER() OVER (
				PARTITION BY attack_name, attack_hero1_id
				ORDER BY attack_hero1_level DESC, time DESC
			) AS rn
		FROM 
			battle_report
		WHERE 
			attack_hero1_id != 0
			AND attack_hero2_id != 0
			AND attack_hero3_id != 0
			AND attack_hero1_level >= 15
			AND attack_hero2_level >= 15
			AND attack_hero3_level >= 15
			AND attack_hp >= 10000
			AND attack_name LIKE '%` + name + `%'
			AND attack_union_name LIKE '%` + uname + `%'
			AND npc = 0
			AND all_skill_info != "" AND all_skill_info IS NOT NULL 
	
		UNION ALL
	
		SELECT 
			defend_name AS player_name,
			defend_hero1_id AS hero1_id,
			defend_hero2_id AS hero2_id,
			defend_hero3_id AS hero3_id,
			defend_hero1_level AS hero1_level,
			defend_hero2_level AS hero2_level,
			defend_hero3_level AS hero3_level,
			defend_hero1_star AS hero1_star,
			defend_hero2_star AS hero2_star,
			defend_hero3_star AS hero3_star,
			defend_total_star AS total_star,
			defend_hp AS hp,
			defender_gear_info AS gear,
			defend_hero_type AS hero_type,
			time,
			all_skill_info,
			battle_id,
			'defend' AS role,
			ROW_NUMBER() OVER (
				PARTITION BY defend_name, defend_hero1_id
				ORDER BY defend_hero1_level DESC, time DESC
			) AS rn
		FROM 
			battle_report
		WHERE 
			defend_hero1_id != 0
			AND defend_hero2_id != 0
			AND defend_hero3_id != 0
			AND defend_hero1_level >= 15
			AND defend_hero2_level >= 15
			AND defend_hero3_level >= 15
			AND defend_hp >= 10000
			AND defend_name LIKE '%` + name + `%'
			AND defend_union_name LIKE '%` + uname + `%'
			AND npc = 0
			AND all_skill_info != "" AND all_skill_info IS NOT NULL 
	),
	deduplicated_data AS (
    SELECT 
        player_name,
        hero1_id,
        hero2_id,
        hero3_id,
        hero1_level,
        hero2_level,
        hero3_level,
        hero1_star,
        hero2_star,
        hero3_star,
        total_star,
        hp,
		gear,
		hero_type,
        time,
        all_skill_info,
        battle_id,
        role,
        ROW_NUMBER() OVER (
            PARTITION BY player_name, hero1_id, hero2_id, hero3_id
            ORDER BY time DESC
        ) AS dedup_rn
    FROM 
        ranked_data
    WHERE 
        rn = 1
)
SELECT 
    player_name,
    hero1_id,
    hero2_id,
    hero3_id,
    hero1_level,
    hero2_level,
    hero3_level,
    hero1_star,
    hero2_star,
    hero3_star,
    total_star,
    hp,
	gear,
	hero_type,
    time,
    all_skill_info,
    battle_id,
    role
FROM 
    deduplicated_data
WHERE 
    dedup_rn = 1
ORDER BY 
    player_name, time DESC;`
	fmt.Println(model.Conn.Raw(query).Scan(&results).Error) // 自动映射到结构体

	// 使用结果
	fmt.Println("找到记录:", len(results))

	common.Response{Data: results}.Success(c)
}

func Example(c *gin.Context) {
	common.Response{Message: "This is example func"}.Success(c)
}
