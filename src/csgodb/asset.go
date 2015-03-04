package csgodb

import (
	"database/sql"
	"steamapi"
	"strconv"
	"fmt"
)

type CSGOAsset struct {
	AssetId uint64
	ClassId uint64
	Name string
	MarketName string
	MarketHashName string
	AssetType string
	IconUrl string
}

func (a CSGOAsset) GetIconUrl() string {
	return fmt.Sprintf("http://cdn.steamcommunity.com/economy/image/%s", a.IconUrl)
}

func FromAssetInfo(ai steamapi.AssetInfo) *CSGOAsset {
	asset := &CSGOAsset{}
	
	_classId, _ := strconv.ParseUint(ai.ClassId, 10, 64)
	asset.ClassId = uint64(_classId)
	asset.AssetType = ai.Type
	asset.IconUrl = ai.IconUrl
	asset.MarketHashName = ai.MarketHashName
	asset.MarketName = ai.MarketName
	asset.Name = ai.Name
	
	return asset
}

func ImportAsset(db *sql.DB, ai steamapi.AssetInfo) *CSGOAsset {
	query := "INSERT INTO local_assets (class_id, asset_name, market_name, market_hash_name, asset_type, icon_url) VALUES (?, ?, ?, ?, ?, ?)"
	_classId, _ := strconv.ParseUint(ai.ClassId, 10, 64)
	db.Exec(query, uint64(_classId), ai.Name, ai.MarketName, ai.MarketHashName, ai.Type, ai.IconUrl)
	return GetAssetByClassId(db, _classId)
}

func GetAssetByClassId(db *sql.DB, classId uint64) *CSGOAsset {
	
	asset := &CSGOAsset{}
	
	query := "SELECT asset_id, class_id, asset_name, market_name, market_hash_name, asset_type, icon_url FROM local_assets WHERE class_id = ?"
	rows, _ := db.Query(query, classId)
	
	for rows.Next() {
		rows.Scan(&asset.AssetId, &asset.ClassId, &asset.Name, &asset.MarketName, &asset.MarketHashName, &asset.AssetType, &asset.IconUrl)
	}
	
	return asset
}

func GetAssetById(db *sql.DB, assetId uint64) *CSGOAsset {
	
	asset := &CSGOAsset{}
	
	query := "SELECT asset_id, class_id, asset_name, market_name, market_hash_name, asset_type, icon_url FROM local_assets WHERE asset_id = ?"
	rows, _ := db.Query(query, assetId)
	
	for rows.Next() {
		rows.Scan(&asset.AssetId, &asset.ClassId, &asset.Name, &asset.MarketName, &asset.MarketHashName, &asset.AssetType, &asset.IconUrl)
	}
	
	return asset
}