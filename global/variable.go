package global

import (
	"context"
	"log"
	"os"
	"path/filepath"
)

type WebExVar struct {
	NeededReportPos        int  //需要获取战报的坐标
	NeedGetReport          bool //是否需要获取战报
	NeedSyncTeamUser       bool //是否需要同步同盟成员信息
	BindIpInfo             bool //是否绑定IP信息 开启后将过滤掉其他IP的数据包(特殊情况使用)
	NeedGetBattleData      bool //是否开启抓取详细战报数据 用于抓取队伍
	NeedPushBookData       bool //是否推送主公簿数据到前端
	NeedPushBattleCallData bool //是否推送战役叫阵数据到前端
}

var ExVar = WebExVar{
	0, false, false, false, false, false, false,
}

var IsDebug bool = false
var Version string = "0.0.4Beta202605030300"
var OnlySrcIp = ""
var OnlyDstIp = ""
var PacketLoss = false
var LossBytes []byte
var LossCmdId = 0
var NeedBufSize = 0
var AppCtx context.Context
var AppDir string

func InitAppDir() {
	if dir := os.Getenv("STZB_APP_DIR"); dir != "" {
		AppDir = dir
	} else {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal("获取程序路径失败:", err)
		}
		exeDir := filepath.Dir(exePath)
		if _, err := os.Stat(filepath.Join(exeDir, "wails.json")); err == nil {
			AppDir = exeDir
		} else {
			wd, _ := os.Getwd()
			if _, err := os.Stat(filepath.Join(wd, "wails.json")); err == nil {
				AppDir = wd
			} else {
				AppDir = exeDir
			}
		}
	}
	AppDir = filepath.ToSlash(AppDir)
	log.Println("应用数据目录:", AppDir)
}
