package hltvscrapper

import (

)


type TeamStat struct {
	MapsPlayed int
	Wins int
	Draws int
	Losses int
	Frags int
	Deaths int
	RoundsPlayed int
	KDRatio float32
}


type Team struct {
	Name string
	TeamId int
	Stats TeamStat
	Players []Player
}
