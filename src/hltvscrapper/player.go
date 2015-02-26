package hltvscrapper

import (

)

type PlayerStat struct {
	Frags int
	Headshots float32
	Deaths int
	KDRatio float32
	MapsPlayed int
	RoundsPlayed int
	AvgFragsPerRound float32
	AvgAssistsPerRound float32
	AvgDeathsPerRound float32
	Rating float32
}

type Player struct {
	Name string
	PlayerId int
	Stats PlayerStat
}

