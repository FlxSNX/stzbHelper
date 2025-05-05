package global

type WebExVar struct {
	NeededReportPos  int  //需要获取战报的坐标
	NeedGetReport    bool //是否需要获取战报
	NeedSyncTeamUser bool //是否需要同步同盟成员信息
}

var ExVar = WebExVar{
	NeededReportPos:  0,
	NeedGetReport:    false,
	NeedSyncTeamUser: false,
}
