package csgopool

import (
	"fmt"
	"eseascrapper"
	"csgodb"
	"time"
	"os"
	"logger"
)


var watcher *WatcherState


type WatcherState struct {
	
	DataPath string
	RefreshTime string
	Running bool
	ImportSnapshot bool
	GenerateSnapshot bool
	SnapshotUrl string
	NoUpdate bool
	Log *logger.LoggerState
}



func NewWatcher(dataPath string, snapshotUrl string, importSnapshot bool, generateSnapshot bool) *WatcherState {
	
	state := &WatcherState{DataPath: dataPath}
	state.Running = false
	state.ImportSnapshot = importSnapshot
	state.GenerateSnapshot = generateSnapshot
	state.SnapshotUrl = snapshotUrl
	state.RefreshTime = "30m"
	state.Log = &logger.LoggerState{LogPath: dataPath+"/watcher.log", Level:3}
	state.NoUpdate = false
	watcher = state
	
	NewPoolState(dataPath + "/settings.json")
	
	return state
}


func (w *WatcherState) LoadData() {
	
	w.Running = true
	//csgodb.Db.Location = "US%2FEastern"
	//init database
	//todo
	//config file

	db_path := w.DataPath + "/db.json"
	csgodb.Db.LoadConfig(db_path)
	
	if len(csgodb.Db.Username) == 0 {
		w.Log.Error("Invalid Database config, please edit db.json!")
		csgodb.Db.SaveConfig(db_path)
		os.Exit(-2)
	}
	
	db, err := csgodb.Db.Open()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	csgodb.InitTables(db)
	
	teams := csgodb.GetAllTeams(db)
	
	if w.GenerateSnapshot {
		w.Log.Info("Generating a Snapshot...")
		snapshot := csgodb.GenerateSnapshot(db)
		snap_path := w.DataPath + "/snapshot.json"
		snapshot.Save(snap_path)
	}
	
	if len(teams) == 0 {
		w.Log.Info("No team found !")
		w.Log.Info("Starting initial import...")
		
		if !w.ImportSnapshot {
			dayDelta := 25
			w.Log.Info(fmt.Sprintf("Importing from ESEA.net (North America), from %d days ago...", dayDelta))
			w.UpdateESEA(dayDelta, "invite", eseascrapper.NorthAmerica)
			w.UpdateESEA(dayDelta, "premier", eseascrapper.NorthAmerica)
			
			w.Log.Info(fmt.Sprintf("Importing from ESEA.net (Europe), from %d days ago...", dayDelta))
			w.UpdateESEA(dayDelta, "invite", eseascrapper.Europe)
			w.UpdateESEA(dayDelta, "premier", eseascrapper.Europe)
			//w.UpdateESEA(dayDelta, "main")
			
			hltvpages := 4
			
			w.Log.Info(fmt.Sprintf("Importing %d matches pages from Hltv.org...", hltvpages))
			w.UpdateHltv(hltvpages)
		} else {
			w.Log.Info(fmt.Sprintf("Import from a snapshot : %s", w.SnapshotUrl))
			
			snapshot := &csgodb.Snapshot{}
			snapshot.ImportFromURL(db, w.SnapshotUrl)
		}
	}
	
	sources := csgodb.GetAllSources(db)
	if len(sources) == 0 {
		csgodb.InsertSource(db, csgodb.HltvSource, "HLTV.org")
		csgodb.InsertSource(db, csgodb.EseaSource, "ESEA.net")
	}

	if csgodb.UsersCount(db) == 0 {
		w.Log.Info("No user found, Creating default PoolMaster Account")
		
		passwd := csgodb.RandomString(12)
		
		csgodb.CreateUser(db, "poolmaster", passwd, "poolmaster@localhost", PoolMaster)
		
		w.Log.Info("PoolMaster Account Created !")
		w.Log.Info(fmt.Sprintf("poolmaster:%s", passwd))
	}
	
	
	db.Close()
	w.Running = false
}

func (w *WatcherState) StartBot()  {
	d, _ := time.ParseDuration(w.RefreshTime)

	w.Log.Info("Starting watcher Bot")
	for {
		if !w.NoUpdate {
				
			w.Running = true
			db, _ := csgodb.Db.Open()
			//updating last events
			
			//hltv
			w.UpdateHltv(0)
			
			//esea
			//parsing today and yesterday
			dayDelta := 2
			w.UpdateESEA(dayDelta, "invite", eseascrapper.NorthAmerica)
			w.UpdateESEA(dayDelta, "premier", eseascrapper.NorthAmerica)
			
			w.UpdateESEA(dayDelta, "invite", eseascrapper.Europe)
			w.UpdateESEA(dayDelta, "premier", eseascrapper.Europe)
			
			//points attribution !
			
			csgodb.InsertWatcherUpdate(db)
			db.Close()
			w.Running = false
		}
		w.Log.Info(fmt.Sprintf("Sleeping %f minutes...", d.Minutes()))
		time.Sleep(d)
	}
	
}