package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"stzbHelper/global"
	"stzbHelper/model"

	"golang.org/x/sys/windows"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	global.AppCtx = ctx
	global.LogW.SetContext(ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetTeamUser() string {
	var teamUsers []model.TeamUser
	query := model.Conn
	query.Find(&teamUsers)

	return global.Response{Data: teamUsers}.Success()
}

// GetTeamGroup 获取所有不重复的分组名称
func (a *App) GetTeamGroup() string {
	var groups []string
	model.Conn.Model(&model.TeamUser{}).Distinct("group").Pluck("group", &groups)
	return global.Response{Data: groups}.Success()
}

// CreateTask 创建攻城任务
func (a *App) CreateTask(name string, tasktime int, target []string, taskpos []string) string {
	task := model.Task{
		Name:   name,
		Time:   tasktime,
		Pos:    model.ToTaskPos(taskpos),
		Target: target,
		Status: 0,
	}

	// 获取目标分组的成员
	var teamUsers []model.TeamUser
	model.Conn.Where("`group` IN ?", target).Find(&teamUsers)
	task.TargetUserNum = len(teamUsers)
	task.UserList = model.TeamUserListToTaskUserList(teamUsers)

	result := model.Conn.Create(&task)
	if result.Error != nil {
		return global.Response{Message: "创建任务失败: " + result.Error.Error()}.Error()
	}

	return global.Response{Data: task, Message: "创建任务成功"}.Success()
}

// GetTaskList 获取任务列表
func (a *App) GetTaskList() string {
	var tasks []model.Task
	model.Conn.Find(&tasks)
	return global.Response{Data: tasks}.Success()
}

// GetGroupWu 获取分组武勋统计
func (a *App) GetGroupWu() string {
	type GroupWu struct {
		Group       string  `json:"group"`
		MemberCount int     `json:"member_count"`
		TotalWu     int     `json:"total_wu"`
		AverageWu   float64 `json:"average_wu"`
		ZeroWuCount int     `json:"zero_wu_count"`
	}

	var results []GroupWu
	model.Conn.Model(&model.TeamUser{}).
		Select("`group`, count(*) as member_count, sum(wu) as total_wu, avg(wu) as average_wu, sum(case when wu = 0 then 1 else 0 end) as zero_wu_count").
		Group("`group`").
		Scan(&results)

	return global.Response{Data: results}.Success()
}

// DeleteTask 删除任务
func (a *App) DeleteTask(id int) string {
	result := model.Conn.Delete(&model.Task{}, id)
	if result.Error != nil {
		return global.Response{Message: "删除任务失败: " + result.Error.Error()}.Error()
	}
	return global.Response{Message: "删除任务成功"}.Success()
}

// EnableGetReport 开启战报获取
func (a *App) EnableGetReport(pos int) string {
	global.ExVar.NeedGetReport = true
	global.ExVar.NeededReportPos = pos
	return global.Response{Message: "开启获取战报成功"}.Success()
}

func (a *App) DisableGetReport() string {
	global.ExVar.NeedGetReport = false
	return global.Response{Message: "停止获取战报"}.Success()
}

// GetReportNumByTaskId 获取某任务的战报数量
func (a *App) GetReportNumByTaskId(id int) string {
	var task model.Task
	model.Conn.First(&task, id)
	if task.Id == 0 {
		return global.Response{Message: "任务不存在"}.Error()
	}

	var count int64
	model.Conn.Model(&model.Report{}).Where("wid = ?", task.Pos).Count(&count)

	return global.Response{Data: map[string]int64{"count": count}}.Success()
}

// StatisticsReport 统计考勤
func (a *App) StatisticsReport(id int) string {
	var task model.Task
	model.Conn.First(&task, id)
	if task.Id == 0 {
		return global.Response{Message: "任务不存在"}.Error()
	}

	// 获取该坐标的所有战报
	var reports []model.Report
	model.Conn.Where("wid = ?", task.Pos).Find(&reports)

	log.Printf("统计任务[%s]的考勤, 坐标%d, 共%d条战报", task.Name, task.Pos, len(reports))

	if task.UserList == nil {
		task.UserList = map[int]*model.TaskUserList{}
	}

	// 统计每个成员的出勤
	for _, report := range reports {
		// 根据 attack_name 匹配成员
		for _, user := range task.UserList {
			if user.Name == report.AttackName {
				// 判断是主力还是拆迁
				if report.AttackHp > 0 {
					user.AtkNum++
					user.AtkTeamNum++
				}
				break
			}
		}
	}

	// 计算实际到的人数
	completeNum := 0
	for _, user := range task.UserList {
		if user.AtkNum > 0 || user.DisNum > 0 {
			completeNum++
		}
	}
	task.CompleteUserNum = completeNum
	task.Status = 1

	model.Conn.Save(&task)

	return global.Response{Message: "统计完成"}.Success()
}

// GetTask 获取任务详情
func (a *App) GetTask(id int) string {
	var task model.Task
	model.Conn.First(&task, id)
	if task.Id == 0 {
		return global.Response{Message: "任务不存在"}.Error()
	}
	return global.Response{Data: task}.Success()
}

// DeleteTaskReport 清理任务战报
func (a *App) DeleteTaskReport(id int) string {
	var task model.Task
	model.Conn.First(&task, id)
	if task.Id == 0 {
		return global.Response{Message: "任务不存在"}.Error()
	}

	// 删除该坐标相关的战报
	model.Conn.Where("wid = ?", task.Pos).Delete(&model.Report{})

	// 重置任务的考勤数据
	task.CompleteUserNum = 0
	task.Status = 0
	for _, user := range task.UserList {
		user.AtkNum = 0
		user.DisNum = 0
		user.AtkTeamNum = 0
		user.DisTeamNum = 0
	}
	model.Conn.Save(&task)

	return global.Response{Message: "清理战报成功"}.Success()
}

// EnableGetBattleReport 开启详细战报获取
func (a *App) EnableGetBattleReport() string {
	global.ExVar.NeedGetBattleData = true
	global.ExVar.NeedGetReport = false
	return global.Response{Message: "开启获取详细战报成功"}.Success()
}

// DisableGetBattleReport 关闭详细战报获取
func (a *App) DisableGetBattleReport() string {
	global.ExVar.NeedGetBattleData = false
	return global.Response{Message: "关闭获取详细战报成功"}.Success()
}

// EnableBookData 开启主公簿数据推送
func (a *App) EnableBookData() string {
	global.ExVar.NeedPushBookData = true
	return global.Response{Message: "开启主公簿数据推送成功"}.Success()
}

// DisableBookData 关闭主公簿数据推送
func (a *App) DisableBookData() string {
	global.ExVar.NeedPushBookData = false
	return global.Response{Message: "关闭主公簿数据推送成功"}.Success()
}

// // EnableBattleCall 开启战役叫阵数据推送
// func (a *App) EnableBattleCall() string {
// 	global.ExVar.NeedPushBattleCallData = true
// 	return global.Response{Message: "开启战役叫阵数据推送成功"}.Success()
// }

// // DisableBattleCall 关闭战役叫阵数据推送
// func (a *App) DisableBattleCall() string {
// 	global.ExVar.NeedPushBattleCallData = false
// 	return global.Response{Message: "关闭战役叫阵数据推送成功"}.Success()
// }

// GetDbList 获取当前目录下的数据库文件列表
func (a *App) GetDbList() string {
	exePath, err := os.Executable()
	if err != nil {
		return global.Response{Message: "获取程序路径失败: " + err.Error()}.Error()
	}
	dir := filepath.Dir(exePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return global.Response{Message: "读取目录失败: " + err.Error()}.Error()
	}

	var dbList []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".db") {
			dbList = append(dbList, strings.TrimSuffix(entry.Name(), ".db"))
		}
	}

	return global.Response{Data: dbList}.Success()
}

// CreateDb 创建新数据库并连接
func (a *App) CreateDb(name string) string {
	if name == "" {
		return global.Response{Message: "数据库名称不能为空"}.Error()
	}
	exePath, err := os.Executable()
	if err != nil {
		return global.Response{Message: "获取程序路径失败: " + err.Error()}.Error()
	}
	dir := filepath.Dir(exePath)
	dbPath := filepath.Join(dir, name)

	model.InitDB(dbPath)
	if model.Conn == nil {
		return global.Response{Message: "创建数据库失败，请检查日志"}.Error()
	}
	databaseSelected = true
	return global.Response{Message: "数据库创建成功"}.Success()
}

// SelectDb 选择并初始化数据库
func (a *App) SelectDb(name string) string {
	exePath, err := os.Executable()
	if err != nil {
		return global.Response{Message: "获取程序路径失败: " + err.Error()}.Error()
	}
	dir := filepath.Dir(exePath)
	dbPath := filepath.Join(dir, name)

	model.InitDB(dbPath)
	if model.Conn == nil {
		return global.Response{Message: "数据库连接失败，请检查日志"}.Error()
	}
	databaseSelected = true
	return global.Response{Message: "数据库连接成功"}.Success()
}

// GetLogs 获取历史日志
func (a *App) GetLogs() string {
	return global.Response{Data: global.LogW.GetLogs()}.Success()
}

// GetVersion 获取当前版本号
func (a *App) GetVersion() string {
	return global.Response{Data: global.Version}.Success()
}

// CheckNpcap 检测 Npcap 是否已安装
func (a *App) CheckNpcap() string {
	dll := windows.NewLazySystemDLL("wpcap.dll")
	err := dll.Load()
	installed := err == nil
	log.Printf("Npcap installed: %v", installed)
	return global.Response{Data: map[string]bool{"installed": installed}}.Success()
}

// CheckUpdate 检查是否有新版本
func (a *App) CheckUpdate() string {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/FlxSNX/stzbHelper/releases/latest")
	if err != nil {
		return global.Response{Message: "检查更新失败: " + err.Error()}.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return global.Response{Data: map[string]interface{}{"hasUpdate": false, "message": "暂无发行版本"}}.Success()
	}

	if resp.StatusCode != 200 {
		return global.Response{Message: "检查更新失败，状态码: " + fmt.Sprint(resp.StatusCode)}.Error()
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return global.Response{Message: "解析更新信息失败: " + err.Error()}.Error()
	}

	hasUpdate := release.TagName != global.Version
	return global.Response{Data: map[string]interface{}{
		"hasUpdate":  hasUpdate,
		"latestVer":  release.TagName,
		"name":       release.Name,
		"body":       release.Body,
		"url":        release.HTMLURL,
		"currentVer": global.Version,
	}}.Success()
}

// GetPlayerTeam 查询玩家队伍
func (a *App) GetPlayerTeam(name string, uname string, idu string, page int, pageSize int) string {
	type PlayerTeam struct {
		PlayerName   string `json:"player_name"`
		BattleID     int    `json:"battle_id"`
		Hero1ID      int    `json:"hero1_id"`
		Hero2ID      int    `json:"hero2_id"`
		Hero3ID      int    `json:"hero3_id"`
		Hero1Level   int    `json:"hero1_level"`
		Hero2Level   int    `json:"hero2_level"`
		Hero3Level   int    `json:"hero3_level"`
		Hero1Star    int    `json:"hero1_star"`
		Hero2Star    int    `json:"hero2_star"`
		Hero3Star    int    `json:"hero3_star"`
		TotalStar    int    `json:"total_star"`
		Hp           int    `json:"hp"`
		AllSkillInfo string `json:"all_skill_info"`
		Role         string `json:"role"`
		Time         int    `json:"time"`
		Gear         string `json:"gear"`
		HeroType     string `json:"hero_type"`
		Idu          string `json:"idu"`
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	namePattern := "%" + name + "%"
	unamePattern := "%" + uname + "%"
	iduPattern := "%" + idu + "%"

	baseQuery := `WITH ranked_data AS (
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
			attack_idu AS idu,
			time,
			all_skill_info,
			battle_id,
			'attack' AS role,
			ROW_NUMBER() OVER (
				PARTITION BY attack_name, attack_hero1_id
				ORDER BY attack_hero1_level DESC, time DESC
			) AS rn
		FROM battle_report
		WHERE attack_hero1_id != 0 AND attack_hero2_id != 0 AND attack_hero3_id != 0
			AND attack_hero1_level >= 15 AND attack_hero2_level >= 15 AND attack_hero3_level >= 15
			AND attack_hp >= 10000
			AND attack_name LIKE ? AND attack_union_name LIKE ? AND attack_idu LIKE ?
			AND npc = 0 AND all_skill_info != "" AND all_skill_info IS NOT NULL
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
			defend_idu AS idu,
			time,
			all_skill_info,
			battle_id,
			'defend' AS role,
			ROW_NUMBER() OVER (
				PARTITION BY defend_name, defend_hero1_id
				ORDER BY defend_hero1_level DESC, time DESC
			) AS rn
		FROM battle_report
		WHERE defend_hero1_id != 0 AND defend_hero2_id != 0 AND defend_hero3_id != 0
			AND defend_hero1_level >= 15 AND defend_hero2_level >= 15 AND defend_hero3_level >= 15
			AND defend_hp >= 10000
			AND defend_name LIKE ? AND defend_union_name LIKE ? AND defend_idu LIKE ?
			AND npc = 0 AND all_skill_info != "" AND all_skill_info IS NOT NULL
	),
	deduplicated_data AS (
		SELECT *, ROW_NUMBER() OVER (
			PARTITION BY player_name, hero1_id, hero2_id, hero3_id
			ORDER BY time DESC
		) AS dedup_rn
		FROM ranked_data WHERE rn = 1
	)`

	args := []interface{}{
		namePattern, unamePattern, iduPattern,
		namePattern, unamePattern, iduPattern,
	}

	// 查询总数
	var total int64
	countQuery := baseQuery + ` SELECT COUNT(*) FROM deduplicated_data WHERE dedup_rn = 1`
	if err := model.Conn.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return global.Response{Message: "查询失败: " + err.Error()}.Error()
	}

	// 分页查询
	offset := (page - 1) * pageSize
	dataQuery := baseQuery + ` SELECT player_name, hero1_id, hero2_id, hero3_id, hero1_level, hero2_level, hero3_level,
		hero1_star, hero2_star, hero3_star, total_star, hp, gear, hero_type, idu,
		time, all_skill_info, battle_id, role
		FROM deduplicated_data WHERE dedup_rn = 1
		ORDER BY player_name, time DESC
		LIMIT ? OFFSET ?`

	var results []PlayerTeam
	if err := model.Conn.Raw(dataQuery, append(args, pageSize, offset)...).Scan(&results).Error; err != nil {
		return global.Response{Message: "查询失败: " + err.Error()}.Error()
	}

	log.Printf("查询玩家队伍: name=%s, union=%s, idu=%s, page=%d, total=%d, 结果: %d条", name, uname, idu, page, total, len(results))
	return global.Response{Data: map[string]interface{}{
		"list":     results,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}}.Success()
}

// GetTeamWinRate 查询队伍胜率统计
func (a *App) GetTeamWinRate(name string, uname string, idu string, page int, pageSize int, minLevel int, minHp int) string {
	type TeamWinRate struct {
		PlayerName   string  `json:"player_name"`
		Hero1Id      int64   `json:"hero1_id"`
		Hero2Id      int64   `json:"hero2_id"`
		Hero3Id      int64   `json:"hero3_id"`
		Hero1Level   int64   `json:"hero1_level"`
		Hero2Level   int64   `json:"hero2_level"`
		Hero3Level   int64   `json:"hero3_level"`
		Hero1Star    int64   `json:"hero1_star"`
		Hero2Star    int64   `json:"hero2_star"`
		Hero3Star    int64   `json:"hero3_star"`
		TotalStar    int64   `json:"total_star"`
		TotalBattles int64   `json:"total_battles"`
		WinCount     int64   `json:"win_count"`
		LossCount    int64   `json:"loss_count"`
		DrawCount    int64   `json:"draw_count"`
		WinRate      float64 `json:"win_rate"`
		LastTime     int64   `json:"last_time"`
		Idu          string  `json:"idu"`
		AllSkillInfo string  `json:"all_skill_info"`
		Role         string  `json:"role"`
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	namePattern := "%" + name + "%"
	unamePattern := "%" + uname + "%"
	iduPattern := "%" + idu + "%"

	// 攻方: result IN (1,2,3,4,10,18,19) 胜, result=0 负, result IN (6,7,8,13) 平
	// 守方: result=0 胜, result IN (1,2,3,4,10,18,19) 负, result IN (6,7,8,13) 平
	baseQuery := `WITH battle_stats AS (
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
			attack_idu AS idu,
			time,
			all_skill_info,
			'attack' AS role,
			CASE WHEN result = 0 THEN 1 ELSE 0 END AS loss,
			CASE WHEN result IN (6,7,8,13) THEN 1 ELSE 0 END AS draw,
			CASE WHEN result IN (1,2,3,4,10,18,19) THEN 1 ELSE 0 END AS win
		FROM battle_report
		WHERE attack_hero1_id != 0 AND attack_hero2_id != 0 AND attack_hero3_id != 0
			AND attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?
			AND attack_hp >= ?
			AND defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?
			AND defend_hp >= ?
			AND LENGTH(all_skill_info) - LENGTH(REPLACE(all_skill_info, ';', '')) = 6
			AND LENGTH(REPLACE(all_skill_info, ',0,', ',')) = LENGTH(all_skill_info)
			AND attack_name LIKE ? AND attack_union_name LIKE ? AND attack_idu LIKE ?
			AND npc = 0 AND result IN (0,1,2,3,4,6,7,8,10,13,18,19)
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
			defend_idu AS idu,
			time,
			all_skill_info,
			'defend' AS role,
			CASE WHEN result IN (1,2,3,4,10,18,19) THEN 1 ELSE 0 END AS loss,
			CASE WHEN result IN (6,7,8,13) THEN 1 ELSE 0 END AS draw,
			CASE WHEN result = 0 THEN 1 ELSE 0 END AS win
		FROM battle_report
		WHERE defend_hero1_id != 0 AND defend_hero2_id != 0 AND defend_hero3_id != 0
			AND defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?
			AND defend_hp >= ?
			AND attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?
			AND attack_hp >= ?
			AND LENGTH(all_skill_info) - LENGTH(REPLACE(all_skill_info, ';', '')) = 6
			AND LENGTH(REPLACE(all_skill_info, ',0,', ',')) = LENGTH(all_skill_info)
			AND defend_name LIKE ? AND defend_union_name LIKE ? AND defend_idu LIKE ?
			AND npc = 0 AND result IN (0,1,2,3,4,6,7,8,10,13,18,19)
	),
	aggregated AS (
		SELECT
			player_name, hero1_id, hero2_id, hero3_id,
			MAX(hero1_level) AS hero1_level,
			MAX(hero2_level) AS hero2_level,
			MAX(hero3_level) AS hero3_level,
			MAX(hero1_star) AS hero1_star,
			MAX(hero2_star) AS hero2_star,
			MAX(hero3_star) AS hero3_star,
			MAX(total_star) AS total_star,
			MAX(idu) AS idu,
			MAX(time) AS last_time,
			SUBSTR(MAX(time || '_' || all_skill_info), INSTR(MAX(time || '_' || all_skill_info), '_') + 1) AS all_skill_info,
			SUBSTR(MAX(time || '_' || role), INSTR(MAX(time || '_' || role), '_') + 1) AS role,
			SUM(win) AS win_count,
			SUM(loss) AS loss_count,
			SUM(draw) AS draw_count,
			COUNT(*) AS total_battles
		FROM battle_stats
		GROUP BY player_name, hero1_id, hero2_id, hero3_id
	)`

	args := []interface{}{
		minLevel, minLevel, minLevel, minHp, minLevel, minLevel, minLevel, minHp, namePattern, unamePattern, iduPattern,
		minLevel, minLevel, minLevel, minHp, minLevel, minLevel, minLevel, minHp, namePattern, unamePattern, iduPattern,
	}

	// 查询总数
	var total int64
	countQuery := baseQuery + ` SELECT COUNT(*) FROM aggregated`
	if err := model.Conn.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return global.Response{Message: "查询失败: " + err.Error()}.Error()
	}

	// 分页查询
	offset := (page - 1) * pageSize
	dataQuery := baseQuery + ` SELECT player_name, hero1_id, hero2_id, hero3_id,
		hero1_level, hero2_level, hero3_level, hero1_star, hero2_star, hero3_star,
		total_star, idu, last_time, all_skill_info, role,
		win_count, loss_count, draw_count, total_battles,
		ROUND(CAST(win_count AS REAL) / total_battles * 100, 1) AS win_rate
		FROM aggregated
		ORDER BY total_battles DESC, win_rate DESC
		LIMIT ? OFFSET ?`

	var results []TeamWinRate
	if err := model.Conn.Raw(dataQuery, append(args, pageSize, offset)...).Scan(&results).Error; err != nil {
		return global.Response{Message: "查询失败: " + err.Error()}.Error()
	}

	log.Printf("查询队伍胜率: name=%s, union=%s, idu=%s, page=%d, total=%d, 结果: %d条", name, uname, idu, page, total, len(results))
	return global.Response{Data: map[string]interface{}{
		"list":     results,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}}.Success()
}

func (a *App) GetTeamWinRateByTeam(name string, uname string, idu string, page int, pageSize int, minLevel int, minHp int) string {
	type TeamWinRateByTeam struct {
		Hero1Id      int64   `json:"hero1_id"`
		Hero2Id      int64   `json:"hero2_id"`
		Hero3Id      int64   `json:"hero3_id"`
		Hero1Level   int64   `json:"hero1_level"`
		Hero2Level   int64   `json:"hero2_level"`
		Hero3Level   int64   `json:"hero3_level"`
		Hero1Star    int64   `json:"hero1_star"`
		Hero2Star    int64   `json:"hero2_star"`
		Hero3Star    int64   `json:"hero3_star"`
		TotalStar    int64   `json:"total_star"`
		TotalBattles int64   `json:"total_battles"`
		WinCount     int64   `json:"win_count"`
		LossCount    int64   `json:"loss_count"`
		DrawCount    int64   `json:"draw_count"`
		WinRate      float64 `json:"win_rate"`
		LastTime     int64   `json:"last_time"`
		AllSkillInfo string  `json:"all_skill_info"`
		Role         string  `json:"role"`
		Players      string  `json:"players"`
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	namePattern := "%" + name + "%"
	unamePattern := "%" + uname + "%"
	iduPattern := "%" + idu + "%"

	baseQuery := `WITH battle_stats AS (
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
			time,
			all_skill_info,
			'attack' AS role,
			CASE WHEN result = 0 THEN 1 ELSE 0 END AS loss,
			CASE WHEN result IN (6,7,8,13) THEN 1 ELSE 0 END AS draw,
			CASE WHEN result IN (1,2,3,4,10,18,19) THEN 1 ELSE 0 END AS win
		FROM battle_report
		WHERE attack_hero1_id != 0 AND attack_hero2_id != 0 AND attack_hero3_id != 0
			AND attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?
			AND attack_hp >= ?
			AND defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?
			AND defend_hp >= ?
			AND LENGTH(all_skill_info) - LENGTH(REPLACE(all_skill_info, ';', '')) = 6
			AND LENGTH(REPLACE(all_skill_info, ',0,', ',')) = LENGTH(all_skill_info)
			AND attack_name LIKE ? AND attack_union_name LIKE ? AND attack_idu LIKE ?
			AND npc = 0 AND result IN (0,1,2,3,4,6,7,8,10,13,18,19)
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
			time,
			all_skill_info,
			'defend' AS role,
			CASE WHEN result IN (1,2,3,4,10,18,19) THEN 1 ELSE 0 END AS loss,
			CASE WHEN result IN (6,7,8,13) THEN 1 ELSE 0 END AS draw,
			CASE WHEN result = 0 THEN 1 ELSE 0 END AS win
		FROM battle_report
		WHERE defend_hero1_id != 0 AND defend_hero2_id != 0 AND defend_hero3_id != 0
			AND defend_hero1_level >= ? AND defend_hero2_level >= ? AND defend_hero3_level >= ?
			AND defend_hp >= ?
			AND attack_hero1_level >= ? AND attack_hero2_level >= ? AND attack_hero3_level >= ?
			AND attack_hp >= ?
			AND LENGTH(all_skill_info) - LENGTH(REPLACE(all_skill_info, ';', '')) = 6
			AND LENGTH(REPLACE(all_skill_info, ',0,', ',')) = LENGTH(all_skill_info)
			AND defend_name LIKE ? AND defend_union_name LIKE ? AND defend_idu LIKE ?
			AND npc = 0 AND result IN (0,1,2,3,4,6,7,8,10,13,18,19)
	),
	aggregated AS (
		SELECT
			hero1_id, hero2_id, hero3_id,
			GROUP_CONCAT(DISTINCT player_name) AS players,
			MAX(hero1_level) AS hero1_level,
			MAX(hero2_level) AS hero2_level,
			MAX(hero3_level) AS hero3_level,
			MAX(hero1_star) AS hero1_star,
			MAX(hero2_star) AS hero2_star,
			MAX(hero3_star) AS hero3_star,
			MAX(total_star) AS total_star,
			MAX(time) AS last_time,
			SUBSTR(MAX(time || '_' || all_skill_info), INSTR(MAX(time || '_' || all_skill_info), '_') + 1) AS all_skill_info,
			SUBSTR(MAX(time || '_' || role), INSTR(MAX(time || '_' || role), '_') + 1) AS role,
			SUM(win) AS win_count,
			SUM(loss) AS loss_count,
			SUM(draw) AS draw_count,
			COUNT(*) AS total_battles
		FROM battle_stats
		GROUP BY hero1_id, hero2_id, hero3_id
	)`

	args := []interface{}{
		minLevel, minLevel, minLevel, minHp, minLevel, minLevel, minLevel, minHp, namePattern, unamePattern, iduPattern,
		minLevel, minLevel, minLevel, minHp, minLevel, minLevel, minLevel, minHp, namePattern, unamePattern, iduPattern,
	}

	dataQuery := baseQuery + ` SELECT hero1_id, hero2_id, hero3_id,
		hero1_level, hero2_level, hero3_level, hero1_star, hero2_star, hero3_star,
		total_star, last_time, all_skill_info, role, players,
		win_count, loss_count, draw_count, total_battles,
		ROUND(CAST(win_count AS REAL) / total_battles * 100, 1) AS win_rate
		FROM aggregated
		ORDER BY total_battles DESC, win_rate DESC`

	var rawResults []TeamWinRateByTeam
	if err := model.Conn.Raw(dataQuery, args...).Scan(&rawResults).Error; err != nil {
		return global.Response{Message: "查询失败: " + err.Error()}.Error()
	}

	// Go 层归一化战法并合并相同队伍
	type teamAcc struct {
		TeamWinRateByTeam
		playerSet map[string]bool
	}
	merged := make(map[string]*teamAcc)
	for _, r := range rawResults {
		// 生成归一化 key: heroIDs + 排序后的战法
		groups := strings.Split(r.AllSkillInfo, ";")
		var skillParts []string
		for _, g := range groups {
			parts := strings.Split(g, ",")
			if len(parts) < 6 {
				continue
			}
			mainSkill := parts[1]
			sub1 := parts[3]
			sub2 := parts[5]
			if sub1 > sub2 {
				sub1, sub2 = sub2, sub1
			}
			skillParts = append(skillParts, mainSkill+"_"+sub1+"_"+sub2)
		}
		key := fmt.Sprintf("%d_%d_%d|%s", r.Hero1Id, r.Hero2Id, r.Hero3Id, strings.Join(skillParts, "|"))

		if existing, ok := merged[key]; ok {
			existing.TotalBattles += r.TotalBattles
			existing.WinCount += r.WinCount
			existing.LossCount += r.LossCount
			existing.DrawCount += r.DrawCount
			if r.LastTime > existing.LastTime {
				existing.LastTime = r.LastTime
				existing.AllSkillInfo = r.AllSkillInfo
				existing.Role = r.Role
			}
			if r.Hero1Level > existing.Hero1Level {
				existing.Hero1Level = r.Hero1Level
			}
			if r.Hero2Level > existing.Hero2Level {
				existing.Hero2Level = r.Hero2Level
			}
			if r.Hero3Level > existing.Hero3Level {
				existing.Hero3Level = r.Hero3Level
			}
			for _, p := range strings.Split(r.Players, ",") {
				if p != "" {
					existing.playerSet[p] = true
				}
			}
		} else {
			ps := make(map[string]bool)
			for _, p := range strings.Split(r.Players, ",") {
				if p != "" {
					ps[p] = true
				}
			}
			merged[key] = &teamAcc{
				TeamWinRateByTeam: r,
				playerSet:         ps,
			}
		}
	}

	// 转换为切片并计算胜率、玩家列表
	var allResults []TeamWinRateByTeam
	for _, v := range merged {
		v.Players = ""
		first := true
		for p := range v.playerSet {
			if first {
				v.Players = p
				first = false
			} else {
				v.Players += "," + p
			}
		}
		if v.TotalBattles > 0 {
			v.WinRate = float64(int(float64(v.WinCount)/float64(v.TotalBattles)*1000)) / 10
		}
		allResults = append(allResults, v.TeamWinRateByTeam)
	}

	// 排序
	sort.Slice(allResults, func(i, j int) bool {
		if allResults[i].TotalBattles != allResults[j].TotalBattles {
			return allResults[i].TotalBattles > allResults[j].TotalBattles
		}
		return allResults[i].WinRate > allResults[j].WinRate
	})

	total := len(allResults)

	// 分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	pageResults := allResults[start:end]

	log.Printf("查询队伍胜率(按队伍): name=%s, union=%s, idu=%s, page=%d, total=%d, 结果: %d条", name, uname, idu, page, total, len(pageResults))
	return global.Response{Data: map[string]interface{}{
		"list":     pageResults,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}}.Success()
}
