# csgopool
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

You start the CS:GO Pool by running csgopoolmain. You need to pass two arguments for CSGO Pool to work

  -data /path/to/desired/data/folder
  -web /path/to/html/file
 
If it's the first run time, all the stats from HLTV will be fetched. This can take about 2 minutes.
  
### What is done

  - HLTV Stats parsing
  - Events parsing
  - Matches parsing
  - Teams parsing
  
  - HLTV Watcher, to get new matches and events

  - Web Interface login
  - Information page
  - Users creation
  
### To Do

This is the TODO list for the near future
  - User password modification
  - User constraint serialization and implementation (password min char, username min and max)
  - Maybe a settings file could be nice
  - Pool
    - Player selection (How player will be selected and constraint)
    - Points attribution per game performance
  - User space
    - Dashboard, My Pool
  - JSON Data or in a Database ? -> MySQL migration in progress
  - Handle Orphan player
  - Pool Master and pool creation page
