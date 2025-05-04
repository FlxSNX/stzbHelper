package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"stzbHelper/model"
)

func ParseData(cmdId int, data []byte) {
	if isDebug == true {
		log.Println("收到[" + strconv.Itoa(cmdId) + "]消息:" + string(parseZlibData(data)))
	}

	if cmdId == 103 {
		parseTeamUser(data)
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

		/*teamUsersJson, err := json.Marshal(teamUsers)
		if err != nil {
			return
		}
		fmt.Println(string(teamUsersJson))*/
		log.Println("同盟成员消息解析成功！共" + strconv.Itoa(len(teamUsers)) + "人")
		model.Conn.Save(teamUsers)
		model.Conn.Not("id", ids).Delete(model.TeamUser{})
		//var records []model.TeamUser
		//model.Conn.Not("id", ids).Find(&records)
		//fmt.Println(records)
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
