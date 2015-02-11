package main

import (
	"fmt"
	"csgopool"
	"csgopoolweb"
	"os"
	"flag"
)

var datapath string
var webroot string

func initArgs() {
	flag.StringVar(&datapath, "data", os.TempDir() , "Path of the Stored GameData")
	flag.StringVar(&webroot, "web", os.TempDir(), "Path of the WebRoot, containing the HTML Page Template")
}

func main() {
	initArgs()
	
	flag.Parse()
	
	fmt.Println("CS GO Pool")
	fmt.Println("Setting DataPath as "+datapath)
	fmt.Println("Web Root as " + webroot)
	
	watcher := csgopool.NewWatcher(datapath)
	watcher.LoadData()
	go watcher.StartBot()
	//starting web here atm
	
	fmt.Println("Starting the Web Server")
	
	web := csgopoolweb.NewWebServer(&watcher.Data, webroot)
	web.Serve()
}