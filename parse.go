package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"stzbHelper/global"
	"stzbHelper/model"
)

func ParseData(cmdId int, data []byte) {
	if isDebug == true {
		log.Println("收到[" + strconv.Itoa(cmdId) + "]消息:" + string(parseZlibData(data)))
	}

	if cmdId == 103 {
		parseTeamUser(data)
	} else if cmdId == 92 {
		parseReport(data)
	}
}

func parseReport(data []byte) {
	log.Println("收到同盟战报消息")
	if global.ExVar.NeedGetReport == false {
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

		//fmt.Println(jsondata[0])
		//var ids []int
		//var teamUsers []model.TeamUser
		//for _, item := range jsondata {
		//	teamUsers = append(teamUsers, model.ToTeamUser(item))
		//	ids = append(ids, int(item[0].(float64)))
		//}
		//
		//log.Println("同盟成员消息解析成功！共" + strconv.Itoa(len(teamUsers)) + "人")
		//model.Conn.Save(teamUsers)
		//model.Conn.Not("id", ids).Delete(model.TeamUser{})
	} else {
		log.Println("解析同盟战报消息失败")
	}
}

func parseTeamUser(data []byte) {
	log.Println("收到同盟成员消息")
	if isDebug == true {
		log.Println(string(parseZlibData(data)))
	}

	msgdata := parseZlibData(data)

	if len(msgdata) > 0 {
		var jsondata [][]any

		json.Unmarshal(msgdata, &jsondata)

		//fmt.Println(jsondata)
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
	if data[0] == 120 && data[1] == 156 {
		compressedReader := bytes.NewReader(data)

		zlibReader, err := zlib.NewReader(compressedReader)
		if err != nil {
			fmt.Println("Error creating zlib reader:", err)
			return []byte{}
		}
		defer zlibReader.Close()

		// Read the uncompressed data
		uncompressedData, err := io.ReadAll(zlibReader)
		if err != nil {
			fmt.Println("Error reading uncompressed data:", err)
			return []byte{}
		}

		return uncompressedData
	}
	return []byte{}
}
