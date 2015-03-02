package csgodb

import (
//	"database/sql"
	"time"
)

type OpenIDEntry struct {
	UserId int
	NonceId string
	EndPoint string
	Nonce string
	Time time.Time
}

