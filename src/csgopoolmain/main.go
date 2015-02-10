package main

import (
	"fmt"
	"csgopool"
	"csgopoolweb"
	"os"
)

func main() {
	fmt.Println("CS GO Pool")
	fmt.Println("Setting DataPath as "+os.TempDir())
	
	watcher := csgopool.NewWatcher(os.TempDir())
	watcher.LoadData()
	
	//starting web here atm
	web := csgopoolweb.NewWebServer(&watcher.Data, "C:\\Users\\JordSti\\git\\csgopool\\html\\")
	web.Serve()
}