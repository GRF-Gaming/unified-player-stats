package models

type Event struct {
	// ApiKey ApiKeyV1    `json:"api_key,required" vd:"len($)>0"`
	ServerName  string      `json:"server_name,required" vd:"len($)>0"`
	ServerGroup string      `json:"server_group,required" vd:"len($)>0"`
	Kills       []EventKill `json:"kills,required" vd:"len($)>0"`
}
