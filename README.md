## csgopool
A CS:GO Pro Players Pool

### Purpose

Our goal is to make a website with a CS:GO Pro players pool. Like for many sports, we want to create a game by choosing CS:GO pro players and with their
 performance, participant will gain some points. All the participants will be ranked by their points. We want to make small season, like 4-6 months pool.
 We want to add a bet system when the pool system will be completed and tested. Maybe a SteamMarket Bet like we can do in CSGO Lounge.

### How to run

The projet is coded in Go langage, so you need to get golang package.

  - You need to get our dependencies first : 
	
	- golang.org/x/crypto/bcrypt
	- github.com/go-sql-driver/mysql
	- github.com/moovweb/gokogiri
	
You start the CS:GO Pool by running csgopoolmain. You need to pass two arguments for CSGO Pool to work

  -data=/path/to/desired/data/folder
  -web=/path/to/html/file
 
If it's the first run time, csgopool will fetch the current default Snapshot that is located in the repository. But you can specified the Snapshot URL if you want a custom Snapshot.

### csgopoolmain
This is the launcher of the csgopool.
You can specify many settings with application switch

  - data=[path/to/data] : Working folder of the application, contains some configuration file
  - web=[path/to/html/file] : Where the HTML Template are located
  - port=[##] : Listening port of the Web Server
  - import=[true|false] : Allow import from a Snapshot, true by default
  - snapurl=[http://url.json] : Url of the Snapshot
  - snapshot=[true|false] : If its true, generate a Snapshot of the current stats into data path
  - refresh=[##m] : Time between HLTV stats update
  - minyear=2015 : Ignore matches before this year, default value is the current year
  

### Modules summary
  
  - csgopool : This is the Watcher and the module will attribute points to users pool when new match will be added.
  - csgodb : This is the Database persistance module
  - csgopoolweb : This module handle web request
  - hltvscrapper : This module fetch Stats from hltv.org
  - eseascrapper : Fetch stats from ESEA.net
  - poolstacker : csgopool instance manager
  
### What is done

  - HLTV Stats parsing
  - ESEA Stats parsing
  - Events parsing (will be removed)
  - Matches parsing
  - Teams parsing
  - News
  - Auto Pool Creation
  
### To Do

This is the TODO list for the near future

  - ESEA parser (nearly done, some errors handling todo)
  - HLTV parser refactor (done, removing events, errors handling)
  - Match revoke, add match to pool
  - Database configuration interface if none found could be nice
  - User constraint serialization and implementation (password min char, username min and max)
  - Maybe a settings file could be nice
  - User Pool Submission
  - User Main page (Dashboard with last matches stats and points attribution)

  
