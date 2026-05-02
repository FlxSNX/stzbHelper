package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"stzbHelper/global"
	"stzbHelper/model"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func parseBookData(data []byte) {
	var raw []interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		log.Println("解析主公簿数据失败:", err)
		return
	}
	if len(raw) < 2 {
		return
	}
	dataMap, ok := raw[1].(map[string]interface{})
	if !ok {
		return
	}

	result := map[string]interface{}{}

	// log.role_name, log.liked
	if logData, ok := dataMap["log"].(map[string]interface{}); ok {
		if v, ok := logData["role_name"].(string); ok {
			result["role_name"] = v
		}
		if v, ok := logData["liked"]; ok {
			result["likes"] = v
		}
	}

	// server[0]
	if server, ok := dataMap["server"].([]interface{}); ok && len(server) > 0 {
		result["server"] = server[0]
	}

	// personal 数组
	if personal, ok := dataMap["personal"].([]interface{}); ok {
		if len(personal) > 14 {
			result["power"] = personal[14]
		}
		if len(personal) > 41 {
			result["main_city_pos"] = personal[41]
		}
	}

	// union [0,"",id,"同盟名","分组名",...]
	if union, ok := dataMap["union"].([]interface{}); ok {
		if len(union) > 3 {
			result["alliance_name"] = union[3]
		}
		if len(union) > 4 {
			result["group_name"] = union[4]
		}
	}

	// history [登录天数, 最高灭敌, 最高武勋, 最高势力, 赛季参与数, 最高攻城, [灭敌数,武将名,武将ID], ...]
	if history, ok := dataMap["history"].([]interface{}); ok {
		if len(history) > 0 {
			result["login_days"] = history[0]
		}
		if len(history) > 1 {
			result["max_season_kills"] = history[1]
		}
		if len(history) > 2 {
			result["max_merit"] = history[2]
		}
		if len(history) > 3 {
			result["max_power"] = history[3]
		}
		if len(history) > 4 {
			result["season_count"] = history[4]
		}
		if len(history) > 5 {
			result["max_season_siege"] = history[5]
		}
		if len(history) > 6 {
			if team, ok := history[6].([]interface{}); ok && len(team) > 2 {
				result["best_team_kills"] = team[0]
				result["best_team_hero_name"] = team[1]
				result["best_team_hero_id"] = team[2]
			}
		}
		if len(history) > 7 {
			result["max_season_demolish"] = history[7]
		}
	}

	// zanAndvistor [访客数, 点赞数, []]
	if zv, ok := dataMap["zanAndvistor"].([]interface{}); ok {
		if len(zv) > 0 {
			result["visitors"] = zv[0]
		}
	}

	// city_card ["[\"ip\",\"code\",\"location\"]", ...]
	if cc, ok := dataMap["city_card"].([]interface{}); ok && len(cc) > 0 {
		if ccStr, ok := cc[0].(string); ok && ccStr != "" {
			var ccArr []interface{}
			if json.Unmarshal([]byte(ccStr), &ccArr) == nil {
				if len(ccArr) > 0 {
					result["ip"] = ccArr[0]
				}
				if len(ccArr) > 2 {
					result["location"] = ccArr[2]
				}
			}
		}
	}

	result["raw"] = dataMap

	if global.AppCtx != nil {
		runtime.EventsEmit(global.AppCtx, "bookData", result)
	}
}

// func parseBattleCallData(data []byte) {
// 	var raw []interface{}
// 	err := json.Unmarshal(data, &raw)
// 	if err != nil {
// 		log.Println("解析战役叫阵数据失败:", err)
// 		return
// 	}

// 	var messages []map[string]interface{}
// 	for _, item := range raw {
// 		entry, ok := item.([]interface{})
// 		if !ok || len(entry) < 2 {
// 			continue
// 		}
// 		fields, ok := entry[1].([]interface{})
// 		if !ok {
// 			continue
// 		}
// 		msg := map[string]interface{}{}
// 		if len(fields) > 4 {
// 			msg["content"] = fields[4]
// 		}
// 		if len(fields) > 5 {
// 			msg["timestamp"] = fields[5]
// 		}
// 		if len(fields) > 7 {
// 			msg["alliance_name"] = fields[7]
// 		}
// 		if len(fields) > 44 {
// 			msg["player_name"] = fields[44]
// 		}
// 		messages = append(messages, msg)
// 	}

// 	if global.AppCtx != nil {
// 		runtime.EventsEmit(global.AppCtx, "battleCallData", messages, string(data))
// 	}
// }

func ParseData(cmdId int, data []byte) {
	if global.IsDebug {
		log.Println("收到[" + strconv.Itoa(cmdId) + "]消息:" + string(parseZlibData(data)))
	}

	if cmdId == 724 {
		// if global.ExVar.NeedPushBattleCallData {
		// 	go parseBattleCallData(parseZlibData(data))
		// }
	} else if cmdId == 103 {
		parseTeamUser(data)
	} else if cmdId == 92 {
		if global.ExVar.NeedGetBattleData {
			log.Println("已开启获取详细战报,目前会暂停考勤战报的获取")
			parseBattleData(data)
		} else {
			parseReport(data)
		}
	}
}

func DecodeType5(data []byte) string {
	if data[0] == 5 {
		result := make([]byte, len(data)-1)
		for index, value := range data[1:] {
			result[index] = value ^ 152
		}
		return string(result)
	}
	return ""
}

// 原始数据结构
type RawData []interface{}

type BattleData struct {
	BattleId              int64       `json:"battle_id"`
	AttackHelpId          string      `json:"attack_help_id"`
	Time                  int64       `json:"time"`
	Wid                   interface{} `json:"wid"`
	WidName               string      `json:"wid_name"`
	AttackName            string      `json:"attack_name"`
	AttackUnionName       string      `json:"attack_union_name"`
	AttackClanName        string      `json:"attack_clan_name"`
	DefendClanName        string      `json:"defend_clan_name"`
	DefendName            string      `json:"defend_name"`
	DefendUnionName       string      `json:"defend_union_name"`
	AttackAdvance         string      `json:"attack_advance"`
	AttackAllHeroInfo     string      `json:"attack_all_hero_info"`
	AttackerGearInfo      string      `json:"attacker_gear_info"`
	DefendAdvance         string      `json:"defend_advance"`
	DefendAllHeroInfo     string      `json:"defend_all_hero_info"`
	DefenderGearInfo      string      `json:"defender_gear_info"`
	AttackHeroType        string      `json:"attack_hero_type"`
	AttackHeroTypeAdvance string      `json:"attack_hero_type_advance"`
	DefendHeroType        string      `json:"defend_hero_type"`
	DefendHeroTypeAdvance string      `json:"defend_hero_type_advance"`
	AttackHp              int64       `json:"attack_hp"`
	DefendHp              int64       `json:"defend_hp"`
	Npc                   int64       `json:"npc"`
	AllSkillInfo          string      `json:"all_skill_info"`
	Result                int64       `json:"result"`
	AttackIdu             string      `json:"attack_idu"` //进攻方队伍ID
	DefendIdu             string      `json:"defend_idu"` //防守方队伍ID
}

func parseBattleData(data []byte) {
	msgdata := parseZlibData(data)
	fmt.Println("原始数据:", string(msgdata))

	if len(msgdata) > 0 {
		var rawData RawData
		err := json.Unmarshal(msgdata, &rawData)
		if err != nil {
			log.Printf("解析JSON失败: %v", err)
			return
		}

		fmt.Printf("数据长度: %d\n", len(rawData))

		// 遍历所有战斗记录
		battleCount := 0
		for _, item := range rawData {
			// 每个item是一个数组 [战斗数据, 其他数据...]
			battleArray, ok := item.([]interface{})
			if !ok || len(battleArray) == 0 {
				continue
			}

			// 第一个元素是战斗数据
			battleMap, ok := battleArray[0].(map[string]interface{})
			if !ok {
				continue
			}

			// 转换为结构体
			var battleData BattleData
			jsonData, err := json.Marshal(battleMap)
			if err != nil {
				log.Printf("转换战斗数据失败: %v", err)
				continue
			}

			if err := json.Unmarshal(jsonData, &battleData); err != nil {
				log.Printf("解析战斗数据失败: %v", err)
				continue
			}

			fmt.Printf("处理战斗ID: %d\n", battleData.BattleId)

			widStr := ""
			switch v := battleData.Wid.(type) {
			case string:
				widStr = v
			case float64:
				widStr = strconv.FormatInt(int64(v), 10)
			case int64:
				widStr = strconv.FormatInt(v, 10)
			case int:
				widStr = strconv.Itoa(v)
			default:
				// 尝试转换为字符串
				widStr = fmt.Sprintf("%v", v)
			}

			// 创建战斗报告
			report := model.BattleReport{
				BattleId:              battleData.BattleId,
				AttackHelpId:          battleData.AttackHelpId,
				Time:                  battleData.Time,
				Wid:                   widStr,
				WidName:               battleData.WidName,
				AttackName:            battleData.AttackName,
				AttackUnionName:       battleData.AttackUnionName,
				AttackClanName:        battleData.AttackClanName,
				DefendClanName:        battleData.DefendClanName,
				DefendName:            battleData.DefendName,
				DefendUnionName:       battleData.DefendUnionName,
				AttackAdvance:         battleData.AttackAdvance,
				AttackAllHeroInfo:     battleData.AttackAllHeroInfo,
				AttackerGearInfo:      battleData.AttackerGearInfo,
				DefendAdvance:         battleData.DefendAdvance,
				DefendAllHeroInfo:     battleData.DefendAllHeroInfo,
				DefenderGearInfo:      battleData.DefenderGearInfo,
				AttackHeroType:        battleData.AttackHeroType,
				AttackHeroTypeAdvance: battleData.AttackHeroTypeAdvance,
				DefendHeroType:        battleData.DefendHeroType,
				DefendHeroTypeAdvance: battleData.DefendHeroTypeAdvance,
				AttackHp:              battleData.AttackHp,
				DefendHp:              battleData.DefendHp,
				Npc:                   battleData.Npc,
				AllSkillInfo:          battleData.AllSkillInfo,
				Result:                battleData.Result,
				AttackIdu:             battleData.AttackIdu,
				DefendIdu:             battleData.DefendIdu,
			}

			// 解析进阶信息和武将信息
			report = parseHeroInfo(report)

			fmt.Printf("保存战斗报告: %+v\n", report)

			// 保存到数据库
			result := model.Conn.Save(&report)
			if result.Error != nil {
				log.Printf("保存战斗报告失败: %v", result.Error)
			} else {
				battleCount++
				fmt.Printf("成功保存战斗报告, ID: %d, 影响行数: %d\n", report.BattleId, result.RowsAffected)
			}
		}

		log.Printf("共处理 %d 条战斗记录", battleCount)
	}
}

// 解析武将信息
func parseHeroInfo(report model.BattleReport) model.BattleReport {
	// 解析进攻方进阶信息
	attackAdvance := splitAndFilter(report.AttackAdvance, ";")
	fmt.Printf("进攻方进阶信息: %v\n", attackAdvance)

	attackTotal := int64(0)
	for i, advance := range attackAdvance {
		if i == 0 { // 跳过第一个元素
			continue
		}
		if len(advance) > 0 {
			star, err := strconv.ParseInt(advance[0], 10, 64)
			if err == nil {
				switch i {
				case 1:
					report.AttackHero1Star = star
				case 2:
					report.AttackHero2Star = star
				case 3:
					report.AttackHero3Star = star
				}
				attackTotal += star
			}
		}
	}
	report.AttackTotalStar = attackTotal

	// 解析防守方进阶信息
	defendAdvance := splitAndFilter(report.DefendAdvance, ";")
	fmt.Printf("防守方进阶信息: %v\n", defendAdvance)

	defendTotal := int64(0)
	for i, advance := range defendAdvance {
		if i == 3 { // 跳过第三个元素
			continue
		}
		if len(advance) > 0 {
			star, err := strconv.ParseInt(advance[0], 10, 64)
			if err == nil {
				switch i {
				case 0:
					report.DefendHero3Star = star
				case 1:
					report.DefendHero2Star = star
				case 2:
					report.DefendHero1Star = star
				}
				defendTotal += star
			}
		}
	}
	report.DefendTotalStar = defendTotal

	// 解析进攻方武将信息
	attackHeroInfo := splitAndFilter(report.AttackAllHeroInfo, ";")
	fmt.Printf("进攻方武将信息: %v\n", attackHeroInfo)

	for i, hero := range attackHeroInfo {
		if len(hero) >= 2 {
			heroID, _ := strconv.ParseInt(hero[0], 10, 64)
			heroLevel, _ := strconv.ParseInt(hero[1], 10, 64)

			switch i {
			case 0:
				report.AttackHero1Id = heroID
				report.AttackHero1Level = heroLevel
			case 1:
				report.AttackHero2Id = heroID
				report.AttackHero2Level = heroLevel
			case 2:
				report.AttackHero3Id = heroID
				report.AttackHero3Level = heroLevel
			}
		}
	}

	// 解析防守方武将信息
	defendHeroInfo := splitAndFilter(report.DefendAllHeroInfo, ";")
	fmt.Printf("防守方武将信息: %v\n", defendHeroInfo)

	for i, hero := range defendHeroInfo {
		if len(hero) >= 2 {
			heroID, _ := strconv.ParseInt(hero[0], 10, 64)
			heroLevel, _ := strconv.ParseInt(hero[1], 10, 64)

			switch i {
			case 0:
				report.DefendHero1Id = heroID
				report.DefendHero1Level = heroLevel
			case 1:
				report.DefendHero2Id = heroID
				report.DefendHero2Level = heroLevel
			case 2:
				report.DefendHero3Id = heroID
				report.DefendHero3Level = heroLevel
			}
		}
	}

	return report
}

// 分割和过滤字符串
func splitAndFilter(input string, delimiter string) [][]string {
	if input == "" {
		return [][]string{}
	}

	parts := strings.Split(input, delimiter)
	var result [][]string

	for _, part := range parts {
		if part != "" {
			// 进一步按逗号分割
			subParts := strings.Split(part, ",")
			var filtered []string
			for _, subPart := range subParts {
				if subPart != "" {
					filtered = append(filtered, subPart)
				}
			}
			if len(filtered) > 0 {
				result = append(result, filtered)
			}
		}
	}

	return result
}

func parseReport(data []byte) {
	log.Println("收到同盟战报消息")
	if !global.ExVar.NeedGetReport {
		log.Println("由于未开启获取战报,本次跳过解析")
		return
	}
	msgdata := parseZlibData(data)
	if len(msgdata) > 0 {
		var jsondata [][]any
		json.Unmarshal(msgdata, &jsondata)

		var reports []model.Report
		var neededreports []model.Report

		for _, v := range jsondata {
			reportJSON, err := json.Marshal(v[0])
			if err != nil {
				fmt.Println("Error marshalling report:", err)
				continue
			}

			var report model.Report
			err = json.Unmarshal(reportJSON, &report)
			if err != nil {
				fmt.Println("Error unmarshalling report:", err)
				continue
			}

			reports = append(reports, report)
			if report.Wid == global.ExVar.NeededReportPos {
				neededreports = append(neededreports, report)
			}
		}

		log.Println("解析同盟战报成功,共" + strconv.Itoa(len(reports)) + "条 符合条件的共" + strconv.Itoa(len(neededreports)) + "条")
		if len(neededreports) > 0 {
			action := model.Conn.Save(&neededreports)
			fmt.Println("数据库共新增" + strconv.Itoa(int(action.RowsAffected)) + "条战报")
		}
	} else {
		log.Println("解析同盟战报消息失败")
	}
}

func parseTeamUser(data []byte) {
	log.Println("收到同盟成员消息")
	if global.IsDebug {
		log.Println(string(parseZlibData(data)))
	}

	msgdata := parseZlibData(data)
	if len(msgdata) > 0 {
		var jsondata [][]any
		json.Unmarshal(msgdata, &jsondata)

		var ids []int
		var teamUsers []model.TeamUser
		for _, item := range jsondata {
			teamUsers = append(teamUsers, model.ToTeamUser(item))
			ids = append(ids, int(item[0].(float64)))
		}

		log.Println("同盟成员消息解析成功！共" + strconv.Itoa(len(teamUsers)) + "人")
		model.Conn.Save(teamUsers)
		model.Conn.Not("id", ids).Delete(model.TeamUser{})
	} else {
		log.Println("解析同盟成员消息失败")
	}
}

func parseZlibData(data []byte) []byte {
	if len(data) >= 2 && data[0] == 120 && data[1] == 156 {
		compressedReader := bytes.NewReader(data)
		zlibReader, err := zlib.NewReader(compressedReader)
		if err != nil {
			fmt.Println("Error creating zlib reader:", err)
			return []byte{}
		}
		defer zlibReader.Close()

		uncompressedData, err := io.ReadAll(zlibReader)
		if err != nil {
			fmt.Println("Error reading uncompressed data:", err)
			return []byte{}
		}
		return uncompressedData
	}
	return data
}
