package models

type EventKill struct {
	PlayerId       string `json:"player_id,required" vd:"len($)>0"`
	PlayerName     string `json:"player_name,required" vd:"len($)>0"`
	KilledEntityId string `json:"killed_entity_id"`
	IsFriendly     bool   `json:"is_friendly,required"`
}
