package model

type BattleReport struct {
	ID                    int64  `json:"id" gorm:"id"`
	BattleId              int64  `json:"battle_id" gorm:"column:battle_id;uniqueIndex"`
	AttackHelpId          string `json:"attack_help_id" gorm:"attack_help_id"`
	Time                  int64  `json:"time" gorm:"time"`
	Wid                   string `json:"wid" gorm:"wid"`
	WidName               string `json:"wid_name" gorm:"wid_name"`                                 // 战斗地点名字
	AttackName            string `json:"attack_name" gorm:"attack_name"`                           // 进攻方名字
	AttackUnionName       string `json:"attack_union_name" gorm:"attack_union_name"`               // 进攻方同盟名字
	AttackClanName        string `json:"attack_clan_name" gorm:"attack_clan_name"`                 // 未知
	DefendClanName        string `json:"defend_clan_name" gorm:"defend_clan_name"`                 // 未知
	DefendName            string `json:"defend_name" gorm:"defend_name"`                           // 防守方名字
	DefendUnionName       string `json:"defend_union_name" gorm:"defend_union_name"`               // 防守方同盟名字
	AttackAdvance         string `json:"attack_advance" gorm:"attack_advance"`                     // 进攻方武将进阶信息
	AttackAllHeroInfo     string `json:"attack_all_hero_info" gorm:"attack_all_hero_info"`         // 进攻方武将信息
	AttackerGearInfo      string `json:"attacker_gear_info" gorm:"attacker_gear_info"`             // 进攻方宝物信息
	DefendAdvance         string `json:"defend_advance" gorm:"defend_advance"`                     // 防守方武将进阶信息
	DefendAllHeroInfo     string `json:"defend_all_hero_info" gorm:"defend_all_hero_info"`         // 防守方武将信息
	DefenderGearInfo      string `json:"defender_gear_info" gorm:"defender_gear_info"`             // 防守方宝物信息
	AttackHeroType        string `json:"attack_hero_type" gorm:"attack_hero_type"`                 // 进攻方武将兵种信息
	AttackHeroTypeAdvance string `json:"attack_hero_type_advance" gorm:"attack_hero_type_advance"` // 进攻方武将兵种进阶信息
	DefendHeroType        string `json:"defend_hero_type" gorm:"defend_hero_type"`                 // 防守方武将兵种信息
	DefendHeroTypeAdvance string `json:"defend_hero_type_advance" gorm:"defend_hero_type_advance"` // 防守方武将兵种进阶信息
	AttackHero1Id         int64  `json:"attack_hero1_id" gorm:"attack_hero1_id"`                   // 进攻方大营武将ID
	AttackHero2Id         int64  `json:"attack_hero2_id" gorm:"attack_hero2_id"`                   // 进攻方中军武将ID
	AttackHero3Id         int64  `json:"attack_hero3_id" gorm:"attack_hero3_id"`                   // 进攻方前锋武将ID
	AttackHero1Level      int64  `json:"attack_hero1_level" gorm:"attack_hero1_level"`             // 进攻方大营武将等级
	AttackHero2Level      int64  `json:"attack_hero2_level" gorm:"attack_hero2_level"`             // 进攻方中军武将等级
	AttackHero3Level      int64  `json:"attack_hero3_level" gorm:"attack_hero3_level"`             // 进攻方前锋武将等级
	AttackHero1Star       int64  `json:"attack_hero1_star" gorm:"attack_hero1_star"`               // 进攻方大营武将红度
	AttackHero2Star       int64  `json:"attack_hero2_star" gorm:"attack_hero2_star"`               // 进攻方中军武将红度
	AttackHero3Star       int64  `json:"attack_hero3_star" gorm:"attack_hero3_star"`               // 进攻方前锋武将红度
	AttackTotalStar       int64  `json:"attack_total_star" gorm:"attack_total_star"`               // 进攻方总红度
	DefendHero1Id         int64  `json:"defend_hero1_id" gorm:"defend_hero1_id"`                   // 防守方大营武将ID
	DefendHero2Id         int64  `json:"defend_hero2_id" gorm:"defend_hero2_id"`                   // 防守方中军武将ID
	DefendHero3Id         int64  `json:"defend_hero3_id" gorm:"defend_hero3_id"`                   // 防守方前锋武将ID
	DefendHero1Level      int64  `json:"defend_hero1_level" gorm:"defend_hero1_level"`             // 防守方大营武将等级
	DefendHero2Level      int64  `json:"defend_hero2_level" gorm:"defend_hero2_level"`             // 防守方中军武将等级
	DefendHero3Level      int64  `json:"defend_hero3_level" gorm:"defend_hero3_level"`             // 防守方前锋武将等级
	DefendHero1Star       int64  `json:"defend_hero1_star" gorm:"defend_hero1_star"`               // 防守方大营武将红度
	DefendHero2Star       int64  `json:"defend_hero2_star" gorm:"defend_hero2_star"`               // 防守方中军武将红度
	DefendHero3Star       int64  `json:"defend_hero3_star" gorm:"defend_hero3_star"`               // 防守方前锋武将红度
	DefendTotalStar       int64  `json:"defend_total_star" gorm:"defend_total_star"`               // 防守方总红度
	AttackHp              int64  `json:"attack_hp" gorm:"attack_hp"`                               // 进攻方总兵力
	DefendHp              int64  `json:"defend_hp" gorm:"defend_hp"`                               // 防守方总兵力
	Npc                   int64  `json:"npc" gorm:"npc"`                                           // 是否为与npc战斗
	AllSkillInfo          string `json:"all_skill_info" gorm:"all_skill_info"`                     // 技能信息
	Result                int64  `json:"result" gorm:"result"`                                     // 战斗结果
}

// TableName 表名称
func (*BattleReport) TableName() string {
	return "battle_report"
}
