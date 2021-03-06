//Package csvauto GENERATED BY CSV AUTO; DO NOT EDIT
package csvauto

//RobotTeamConfig auto
type RobotTeamConfig struct {
	ID                     int   //#
	LevelRange             []int //机器人生成时的等级加减范围值
	InitialHelp            int   //机器人生成时的初始帮助数
	InitialLike            []int //机器人生成时的初始点赞数
	InitialTimeRange       []int //机器人生成检测时间间隔（秒）（X1-X2）
	RobotNumMaxlimit       int   //机器人加入条件：机器人数量上限（A）
	TeammemberNumLimit     int   //机器人加入条件：俱乐部真实用户数小于该值（B）
	GenerateRobotNum       []int //单次生成机器人的数量(Y1-Y2)
	RobotNumMinlimit       int   //机器人清退条件：机器人数量下限（C）
	RobotKickTimerange     []int //机器人清退条件：清退检测时间（M1|M2）秒
	RobotKickLevel         []int //机器人清退条件：前24小时俱乐部活跃用户推关数（Z）
	RobotKickActiveusernum int   //机器人清退条件：前24小时俱乐部满足关卡条件的活跃用户数（N）
	RobotKickRobotnum      []int //机器人清退条件：单次清退机器人数量
	LifeRequestTimegap     []int //机器人生命请求延迟（秒）
	HelpTimegap            []int //机器人帮助点体力延迟（秒）
	JoinTalkTimegap        []int //机器人初次入会发送问候语延迟（秒）语言表类型1
	HelpTalkTimegap        []int //机器人请求体力后发送语言延迟（秒）语言表类型2
}

//IRobotTeamConfig auto
type IRobotTeamConfig interface {
	GetID() int
	GetLevelRangeLen() int
	GetLevelRangeByIndex(index int) int
	GetInitialHelp() int
	GetInitialLikeLen() int
	GetInitialLikeByIndex(index int) int
	GetInitialTimeRangeLen() int
	GetInitialTimeRangeByIndex(index int) int
	GetRobotNumMaxlimit() int
	GetTeammemberNumLimit() int
	GetGenerateRobotNumLen() int
	GetGenerateRobotNumByIndex(index int) int
	GetRobotNumMinlimit() int
	GetRobotKickTimerangeLen() int
	GetRobotKickTimerangeByIndex(index int) int
	GetRobotKickLevelLen() int
	GetRobotKickLevelByIndex(index int) int
	GetRobotKickActiveusernum() int
	GetRobotKickRobotnumLen() int
	GetRobotKickRobotnumByIndex(index int) int
	GetLifeRequestTimegapLen() int
	GetLifeRequestTimegapByIndex(index int) int
	GetHelpTimegapLen() int
	GetHelpTimegapByIndex(index int) int
	GetJoinTalkTimegapLen() int
	GetJoinTalkTimegapByIndex(index int) int
	GetHelpTalkTimegapLen() int
	GetHelpTalkTimegapByIndex(index int) int
}

//GetID auto
func (r *RobotTeamConfig) GetID() int {
	return r.ID
}

//GetLevelRangeLen auto
func (r *RobotTeamConfig) GetLevelRangeLen() int {
	return len(r.LevelRange)
}

//GetLevelRangeByIndex auto
func (r *RobotTeamConfig) GetLevelRangeByIndex(index int) int {
	return r.LevelRange[index]
}

//GetInitialHelp auto
func (r *RobotTeamConfig) GetInitialHelp() int {
	return r.InitialHelp
}

//GetInitialLikeLen auto
func (r *RobotTeamConfig) GetInitialLikeLen() int {
	return len(r.InitialLike)
}

//GetInitialLikeByIndex auto
func (r *RobotTeamConfig) GetInitialLikeByIndex(index int) int {
	return r.InitialLike[index]
}

//GetInitialTimeRangeLen auto
func (r *RobotTeamConfig) GetInitialTimeRangeLen() int {
	return len(r.InitialTimeRange)
}

//GetInitialTimeRangeByIndex auto
func (r *RobotTeamConfig) GetInitialTimeRangeByIndex(index int) int {
	return r.InitialTimeRange[index]
}

//GetRobotNumMaxlimit auto
func (r *RobotTeamConfig) GetRobotNumMaxlimit() int {
	return r.RobotNumMaxlimit
}

//GetTeammemberNumLimit auto
func (r *RobotTeamConfig) GetTeammemberNumLimit() int {
	return r.TeammemberNumLimit
}

//GetGenerateRobotNumLen auto
func (r *RobotTeamConfig) GetGenerateRobotNumLen() int {
	return len(r.GenerateRobotNum)
}

//GetGenerateRobotNumByIndex auto
func (r *RobotTeamConfig) GetGenerateRobotNumByIndex(index int) int {
	return r.GenerateRobotNum[index]
}

//GetRobotNumMinlimit auto
func (r *RobotTeamConfig) GetRobotNumMinlimit() int {
	return r.RobotNumMinlimit
}

//GetRobotKickTimerangeLen auto
func (r *RobotTeamConfig) GetRobotKickTimerangeLen() int {
	return len(r.RobotKickTimerange)
}

//GetRobotKickTimerangeByIndex auto
func (r *RobotTeamConfig) GetRobotKickTimerangeByIndex(index int) int {
	return r.RobotKickTimerange[index]
}

//GetRobotKickLevelLen auto
func (r *RobotTeamConfig) GetRobotKickLevelLen() int {
	return len(r.RobotKickLevel)
}

//GetRobotKickLevelByIndex auto
func (r *RobotTeamConfig) GetRobotKickLevelByIndex(index int) int {
	return r.RobotKickLevel[index]
}

//GetRobotKickActiveusernum auto
func (r *RobotTeamConfig) GetRobotKickActiveusernum() int {
	return r.RobotKickActiveusernum
}

//GetRobotKickRobotnumLen auto
func (r *RobotTeamConfig) GetRobotKickRobotnumLen() int {
	return len(r.RobotKickRobotnum)
}

//GetRobotKickRobotnumByIndex auto
func (r *RobotTeamConfig) GetRobotKickRobotnumByIndex(index int) int {
	return r.RobotKickRobotnum[index]
}

//GetLifeRequestTimegapLen auto
func (r *RobotTeamConfig) GetLifeRequestTimegapLen() int {
	return len(r.LifeRequestTimegap)
}

//GetLifeRequestTimegapByIndex auto
func (r *RobotTeamConfig) GetLifeRequestTimegapByIndex(index int) int {
	return r.LifeRequestTimegap[index]
}

//GetHelpTimegapLen auto
func (r *RobotTeamConfig) GetHelpTimegapLen() int {
	return len(r.HelpTimegap)
}

//GetHelpTimegapByIndex auto
func (r *RobotTeamConfig) GetHelpTimegapByIndex(index int) int {
	return r.HelpTimegap[index]
}

//GetJoinTalkTimegapLen auto
func (r *RobotTeamConfig) GetJoinTalkTimegapLen() int {
	return len(r.JoinTalkTimegap)
}

//GetJoinTalkTimegapByIndex auto
func (r *RobotTeamConfig) GetJoinTalkTimegapByIndex(index int) int {
	return r.JoinTalkTimegap[index]
}

//GetHelpTalkTimegapLen auto
func (r *RobotTeamConfig) GetHelpTalkTimegapLen() int {
	return len(r.HelpTalkTimegap)
}

//GetHelpTalkTimegapByIndex auto
func (r *RobotTeamConfig) GetHelpTalkTimegapByIndex(index int) int {
	return r.HelpTalkTimegap[index]
}
