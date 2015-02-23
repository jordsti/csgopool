package main

import (
	"fmt"
	"net/http"
	"csgodb"
	"os"
	"os/exec"
	"flag"

)

type Stacks struct {
	Instances []*StackInfo
}

var passKey string
var rootFolder string
var db *csgodb.Database
var stacks *Stacks

var currentId int


func initArgs() {
	flag.StringVar(&rootFolder, "root", os.TempDir(), "Root stacker folder")
	flag.StringVar(&passKey, "passkey", csgodb.RandomString(48), "Pass Key")
}

func main() {
	currentId = 8000
	initArgs()
	flag.Parse()
	
	stacks = &Stacks{}
	
	fmt.Println("CS:GO Pool Stacker")
	stackerPort := 14000
	
	fmt.Printf("Current PassKey : %s\n", passKey)
	
	//loading db info
	db_config := rootFolder + "/db.json"
	db = &csgodb.Database{}
	db.LoadConfig(db_config)
	
	if len(db.Username) == 0 {
		fmt.Printf("Edit your database configuration at %s\n", db_config)
		db.SaveConfig(db_config)
		os.Exit(-1)
	}
	
	http.HandleFunc("/", StackerHomeHandler)
	
	http.ListenAndServe(fmt.Sprintf(":%d", stackerPort), nil)
}

func StackerHomeHandler(w http.ResponseWriter, r *http.Request) {
	rpassKey := r.FormValue("passkey")
	
	if rpassKey == passKey {
		action := r.FormValue("action")
		
		if action == "launch" {
			
			//launching new stack
			port := currentId
			
			currentId++
			fmt.Printf("Launching instance ! on port %d\n", port)
			stack := DefaultPoolStack()
			stack.GenerateId()
			stack.PrepareDatabase()
			//creating folder of the instance
			stack_path := rootFolder + "/" + stack.Id + "/"
			stack.Port = port
			stack.DataPath = stack_path
			stack.WebRoot = stack_path + "/csgopool/html/"
			
			os.Mkdir(stack.DataPath, 0644)
			
			stacks.Instances = append(stacks.Instances, stack)
			
			//cloning repo
			
			c := exec.Command("git", "clone", stack.GitUrl, stack_path + "/csgopool")
			fmt.Println("Cloning CSGO Pool...")
			
			err := c.Run()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			
			//creating db config file
			
			stack_db := &csgodb.Database{}
			stack_db.Name = "s_" + stack.Id
			stack_db.Username = "s_" + stack.Id
			stack_db.Password = stack.DbPassword
			stack_db.Location = db.Location
			stack_db.Address = db.Address
			
			stack_db.SaveConfig(stack.DataPath + "/db.json")
			
			//launch watcher and web service
			go LaunchStack(stack)
 		}
		
	} else {
		fmt.Fprintf(w, "Bad Passkey")
	}
}

func LaunchStack(s *StackInfo) {
	
	
	/*csgoscrapper.NewScrapperState(time.Now().Year())
	
	csgoscrapper.NewLogger(s.DataPath+"/scrapper.log", 3)
	
	watcher := csgopool.NewWatcher(s.DataPath, s.SnapshotUrl, true, false)
	watcher.RefreshTime = "20m"
	watcher.LoadData()
	watcher.NoUpdate = false
	go watcher.StartBot()
	//starting web here atm
	
	fmt.Println("[CSGOPOOLMAIN] Starting the Web Server")
	
	web_log := s.DataPath + "/csgopoolweb.log"
	
	web := csgopoolweb.NewWebServer(&watcher.Data, &watcher.Users, s.Port, s.WebRoot, web_log)
	web.Serve()*/
}