//Package csvauto GENERATED BY CSV AUTO; DO NOT EDIT
package csvauto

//Guildrecommend auto
type Guildrecommend struct {
	ID           int //ID
	GuildCnt     int //公会数量
	MinMemberCnt int //最少人数
	MaxMemberCnt int //最大人数
}

//IGuildrecommend auto
type IGuildrecommend interface {
	GetID() int
	GetGuildCnt() int
	GetMinMemberCnt() int
	GetMaxMemberCnt() int
}

//GetID auto
func (g *Guildrecommend) GetID() int {
	return g.ID
}

//GetGuildCnt auto
func (g *Guildrecommend) GetGuildCnt() int {
	return g.GuildCnt
}

//GetMinMemberCnt auto
func (g *Guildrecommend) GetMinMemberCnt() int {
	return g.MinMemberCnt
}

//GetMaxMemberCnt auto
func (g *Guildrecommend) GetMaxMemberCnt() int {
	return g.MaxMemberCnt
}
