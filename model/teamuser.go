package model

type TeamUser struct {
	Id              int    `json:"id" gorm:"column:id"`
	Name            string `json:"name" gorm:"column:name"`
	ContributeTotal int    `json:"contribute_total" gorm:"column:contribute_total"`
	ContributeWeek  int    `json:"contribute_week" gorm:"column:contribute_week"`
	Pos             int    `json:"pos" gorm:"column:pos"`
	Power           int    `json:"power" gorm:"column:power"`
	Wu              int    `json:"wu" gorm:"column:wu"`
	Group           string `json:"group" gorm:"column:group"`
	JoinTime        int    `json:"join_time" gorm:"column:join_time"`
}

func (TeamUser) TableName() string {
	return "teamuser"
}

/*
		同盟成员信息索引
		[
			//0 id?
	        //1 名称
			//2 总贡献
			//6 坐标
			//7 本周贡献
			//8 势力
			//10 武勋
			//13 分组
			//30 应该是加入时间
		]
*/

func ToTeamUser(data []any) TeamUser {
	if len(data) < 31 {
		return TeamUser{}
	}

	if data[13].(string) == "" {
		data[13] = "未分组"
	}

	teamUser := TeamUser{
		Id:              int(data[0].(float64)),
		Name:            data[1].(string),
		ContributeTotal: int(data[2].(float64)),
		ContributeWeek:  int(data[7].(float64)),
		Pos:             int(data[6].(float64)),
		Power:           int(data[8].(float64)),
		Wu:              int(data[10].(float64)),
		Group:           data[13].(string),
		JoinTime:        int(data[30].(float64)),
	}

	return teamUser
}
