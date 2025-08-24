package global

type WebExVar struct {
	NeededReportPos  int  //需要获取战报的坐标
	NeedGetReport    bool //是否需要获取战报
	NeedSyncTeamUser bool //是否需要同步同盟成员信息
	BindIpInfo       bool //是否绑定IP信息 开启后将过滤掉其他IP的数据包(特殊情况使用)
}

var ExVar = WebExVar{
	0, false, false, false,
}

var IsDebug bool = false
var Version string = "0.0.3"
var OnlySrcIp = ""
var OnlyDstIp = ""
