package models

import "time"

type EventKillFromGame struct {
	PlayerId       string    `json:"player_id,required" vd:"len($)>0"`   // PlayerId is the UID of the player designated by the server
	PlayerName     string    `json:"player_name,required" vd:"len($)>0"` // PlayerName is the in-game name of the player
	KilledEntityId string    `json:"killed_entity_id,required"`          // KilledEntityId only specifies the killed PLAYER; "" => AI
	IsFriendly     bool      `json:"is_friendly,required"`               // IsFriendly specifies if the kill was a friendly kill
	Time           time.Time `json:"time,required"`                      // Time Recorded of kill event from server
}

type EventKillRecord struct {
	EventKillFromGame
	//ServerGroup string `json:"server_group,required"` // ServerGroup is the logical grouping of events
	ServerName string `json:"server_name,required"` // ServerName is the name of the server the kill occurred
	GameName   string `json:"game_name,required"`   // GameName is the name of the physical game (e.g. "arma_reforger")
}
