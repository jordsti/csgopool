package main

import (
	"fmt"
	"csgopool"
	"csgopoolweb"
	"csgoscrapper"
	"os"
	"flag"
)

var datapath string
var webroot string
var webport int

func initArgs() {
	flag.StringVar(&datapath, "data", os.TempDir() , "Path of the Stored GameData")
	flag.StringVar(&webroot, "web", os.TempDir(), "Path of the WebRoot, containing the HTML Page Template")
	flag.IntVar(&webport, "port", 8080, "Listening port on the web service")
}

func main() {
	initArgs()
	
	flag.Parse()
	
	fmt.Println("[CSGOPOOLMAIN] CS GO Pool")
	fmt.Println("[CSGOPOOLMAIN] Setting DataPath as "+datapath)
	fmt.Println("[CSGOPOOLMAIN] Web Root as " + webroot)
	fmt.Printf("[CSGOPOOLMAIN] Web Service Port : %d\n", webport)
	
	csgoscrapper.NewLogger(datapath+"/scrapper.log", 3)
	
	watcher := csgopool.NewWatcher(datapath)
	watcher.LoadData()
	go watcher.StartBot()
	//starting web here atm
	
	fmt.Println("[CSGOPOOLMAIN] Starting the Web Server")
	
	web_log := datapath + "/csgopoolweb.log"
	
	web := csgopoolweb.NewWebServer(&watcher.Data, &watcher.Users, webport, webroot, web_log)
	web.Serve()
}