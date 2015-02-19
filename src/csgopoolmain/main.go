package main

import (
	"fmt"
	"csgopool"
	"csgopoolweb"
	"csgoscrapper"
	"os"
	"flag"
	"time"
)

var datapath string
var webroot string
var webport int
var importSnapshot bool
var snapshot bool
var snapshotUrl string
var refreshTime string
var minYear int

func initArgs() {
	flag.StringVar(&datapath, "data", os.TempDir() , "Path of the Stored Configuration and Data")
	flag.StringVar(&webroot, "web", os.TempDir(), "Path of the WebRoot, containing the HTML Page Template")
	flag.IntVar(&webport, "port", 8000, "Listening port on the web service")
	flag.BoolVar(&importSnapshot, "import", true, "Import from a snapshot")
	flag.StringVar(&snapshotUrl, "snapurl", "https://raw.githubusercontent.com/jordsti/csgopool/master/snapshots/snapshot-current.json", "Snapshot Url")
	flag.BoolVar(&snapshot, "snapshot", false, "Generate a snapshot")
	flag.StringVar(&refreshTime, "refresh", "30m", "Time between each HLTV update")
	flag.IntVar(&minYear, "minyear", time.Now().Year(), "Minimum year to parse, before this year, matches will be ignored")
}

func main() {
	initArgs()
	
	flag.Parse()
	
	fmt.Println("[CSGOPOOLMAIN] CS GO Pool")
	fmt.Println("[CSGOPOOLMAIN] Setting DataPath as "+datapath)
	fmt.Println("[CSGOPOOLMAIN] Web Root as " + webroot)
	fmt.Printf("[CSGOPOOLMAIN] Web Service Port : %d\n", webport)
	fmt.Printf("[CSGOPOOLMAIN] Snapshot Url : %s\n", snapshotUrl)
	fmt.Printf("[CSGOPOOLMAIN] Minimum Year : %d\n", minYear)
	
	csgoscrapper.NewScrapperState(minYear)
	
	csgoscrapper.NewLogger(datapath+"/scrapper.log", 3)
	
	watcher := csgopool.NewWatcher(datapath, snapshotUrl, importSnapshot, snapshot)
	watcher.RefreshTime = refreshTime
	watcher.LoadData()
	go watcher.StartBot()
	//starting web here atm
	
	fmt.Println("[CSGOPOOLMAIN] Starting the Web Server")
	
	web_log := datapath + "/csgopoolweb.log"
	
	web := csgopoolweb.NewWebServer(&watcher.Data, &watcher.Users, webport, webroot, web_log)
	web.Serve()
}