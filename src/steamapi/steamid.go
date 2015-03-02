package steamapi

import (
	"fmt"
)

type SteamID struct {
	AccountId int64
	Type int
}

func (sid *SteamID) Equals(other *SteamID) bool {
	return sid.AccountId == other.AccountId
}

func (sid *SteamID) String() string {
	return fmt.Sprintf(`STEAM_0:%d:%d`, sid.Type, sid.AccountId)
}
