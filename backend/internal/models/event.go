package models

import (
	"time"
)

type EventFromGame struct {
	// ApiKey Reforger does not permit setting of headers
	// ApiKey ApiKeyV1    `json:"api_key,required" vd:"len($)>0"`

	// ServerName Name of the server event(s) were observed from
	ServerName string `json:"server_name,required" vd:"len($)>0"`

	// GameName Name of the game used for aggregation by player and game
	GameName string `json:"game_name,required" vd:"len($)>0"`

	// ServerGroup Used for logical grouping for further analysis
	//ServerGroup string              `json:"server_group,required" vd:"len($)>0"`

	// Kills Recorded kill events on the server
	Kills []EventKillFromGame `json:"kills,required" vd:"len($)>0"`

	// Players is a list of players currently in the server at given time event
	Players PlayerIdMap `json:"players"`

	// Time is the timestamp for this update from the server
	Time time.Time `json:"time"`
}

func (e *EventFromGame) ExtractKillRecords() []*EventKillRecord {

	var allKills []*EventKillRecord

	for _, k := range e.Kills {

		kill := &EventKillRecord{
			EventKillFromGame: k,
			ServerName:        e.ServerName,
			GameName:          e.GameName,
		}

		if kill.Time.IsZero() {
			kill.Time = e.Time
		}

		allKills = append(allKills, kill)
	}

	return allKills
}
