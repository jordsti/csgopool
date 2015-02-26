package hltvscrapper

import (

)

type ScrapperState struct {
	MinYear int
}

var State *ScrapperState

func NewScrapperState(minYear int) {
	State = &ScrapperState{MinYear: minYear}
}