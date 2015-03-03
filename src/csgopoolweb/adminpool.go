package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"csgodb"
	"csgopool"
	//"github.com/akavel/go-openid"
	
)

type AdminPoolPage struct {
	Page
	Content template.HTML
}

func ReadFile(relpath string) string {
	
	b, err := ioutil.ReadFile(state.RootPath + relpath)
	if err != nil {
		state.Log.Error(fmt.Sprintf("Error while loading %s", relpath))
	}
	
	return string(b)
}

func AdminPoolHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	if !session.IsLogged() {
		http.Redirect(w, r, "/", 301)
	} else {
		if !session.User.IsPoolMaster() {
			state.Log.Debug("Not PoolMaster!")
			http.Redirect(w, r, "/", 301)
		}
	}
	
	msgHtml := ""
	if session.IsFieldExists("message") {
		field := session.PopField("message")
		msgHtml = fmt.Sprintf(`<div>%s</div>`, field.Value)
	}
	
	t, err := MakeTemplate("adminpool.html")
	if err != nil {
		state.Log.Error(fmt.Sprintf("%s", err))
	}
	
	action := r.FormValue("action")
	p := &AdminPoolPage{}
	
	
	if len(action) == 0 || action == "menu" {
		//home page
		p.Content = template.HTML(ReadFile("adminpoolmenu.html"))
	} else if action == "createpool" {
		p.Content = template.HTML(ReadFile("createpool.html"))
	} else if action == "createdivision" {
		
		db, _ := csgodb.Db.Open()
		
		dcount := csgodb.DivisionCount(db)
		
		db.Close()
		
		if dcount == 0 {
			_nb_div := r.FormValue("nbdiv")
			_tmp, _ := strconv.ParseInt(_nb_div, 10, 32)
			nb_div := int(_tmp)
			it := 0
			div_html := `<h4>Divisions</h4><form method="POST" action="/adminpool/?action=submitdiv">`
			
			for ; it < nb_div ; it++ {
	
				div_html += fmt.Sprintf(`<div class="form-group"><h5>Division %d</h5>`, it+1)
				div_html += fmt.Sprintf(`<div class="form-group"><label for="div_%d_name">Division Name</label><input class="form-control" type="text" name="div_%d_name" id="div_%d_name" value="Division %d"/></div>`, it, it, it, it + 1 )
				div_html += fmt.Sprintf(`<div class="form-group"><label for="div_%d_players">Players</label><input class="form-control" type="text" name="div_%d_players" id="div_%d_players" /></div>`, it, it, it )
				div_html += `</div>`
				
			}
			
			div_html += fmt.Sprintf(`<input type="hidden" name="nbdiv" id="nbdiv»" value="%d" />`, nb_div)
			div_html += `<button class="btn btn-default" type="submit">Create Pool</button>`
			div_html += `</form>`
			
			p.Content = template.HTML(div_html)
		} else {
			adminLink := &Link{Caption: "Back", Url:"/adminpool/"}
			adminLink.AddParameter("action", "menu")
			
			div_html := fmt.Sprintf(`<h4>Divions</h4><p>A Pool is already existing, clear it before <br />%s</p>`, adminLink.GetHTML())
			p.Content = template.HTML(div_html)
		}
	} else if action == "submitdiv" {
		_nb_div := r.FormValue("nbdiv")
		_tmp, _ := strconv.ParseInt(_nb_div, 10, 32)
		nb_div := int(_tmp)
		
		div_html := "<h4>Divisions</h4>"
		
		db, _ := csgodb.Db.Open()
		dcount := csgodb.DivisionCount(db)
		
		if dcount == 0 {
			div_html += `<div>`
			for it := 0; it < nb_div; it++ {
				div_name_f := fmt.Sprintf("div_%d_name", it)
				div_players_f := fmt.Sprintf("div_%d_players", it)
				
				div_name := r.FormValue(div_name_f)
				players := r.FormValue(div_players_f)
				pl_id := strings.Split(players, ";")
				
				div_html += fmt.Sprintf(`<div><h5>%s</h5>Players<br /><ul>`, div_name)
				division := csgodb.AddDivision(db, div_name)
				for _, p_id := range pl_id {
					_p_id, _ := strconv.ParseInt(p_id, 10, 32)
					ip_id := int(_p_id)
					
					if ip_id != 0 {
						player := csgodb.GetPlayerById(db, ip_id)
						
						if player != nil {
						division.AddPlayer(db, player.PlayerId)
						pLink := &Link{Caption: player.Name, Url:"/viewplayer/"}
						pLink.AddInt("id", ip_id)
						div_html += fmt.Sprintf("<li>%s</li>", pLink.GetHTML())
						
						} else {
							div_html += fmt.Sprintf("<li>%d</li>", ip_id)
						}
					}
				}
				div_html += `</ul></div>`
	
			}
			
			div_html += `</div>`
		} else {
			
			adminLink := &Link{Caption: "Back", Url:"/adminpool/"}
			adminLink.AddParameter("action", "menu")
			
			div_html += fmt.Sprintf(`<p>A Pool is already existing, clear it before <br />%s</p>`, adminLink.GetHTML())
		}
		db.Close()
		p.Content = template.HTML(div_html)
		
		
	} else if action == "clearpool" {
		
		clear_html := `<h4>Clear current Pool</h4>`
		sure := r.FormValue("sure")
		if sure == "yes" {
			db, _ := csgodb.Db.Open()
			csgodb.ClearPool(db)
			db.Close()
			
			clear_html += `<h4>Pool Cleared !</h4>`
			
		} else {
			clearLink := &Link{Caption: "Yes", Url:"/adminpool/"}
			clearLink.AddParameter("action", "clearpool")
			clearLink.AddParameter("sure", "yes")
			
			adminLink := &Link{Caption: "No", Url:"/adminpool/"}
			adminLink.AddParameter("action", "menu")
			
			clear_html += fmt.Sprintf(`<p>Are you sure you want to clear the current pool ? <br /> %s - %s</p>`, clearLink.GetHTML(), adminLink.GetHTML())
		}
		
		p.Content = template.HTML(clear_html)
	} else if action == "autogenerateform" {
		p.Content = template.HTML(ReadFile("autocreatepool.html"))
	} else if action == "autocreatedivision" {
		db, _ := csgodb.Db.Open()
		dcount := csgodb.DivisionCount(db)
		
		_nb_div := r.FormValue("nbdiv")
		_tmp, _ := strconv.ParseInt(_nb_div, 10, 32)
		nb_div := int(_tmp)
		
		_nb_player := r.FormValue("nb_player")
		_tmp, _ = strconv.ParseInt(_nb_player, 10, 32)
		nb_player := int(_tmp)
		
		//custom query sometime, or radio button
		//todo, to be analyze
		
		pid := 0
		players_id := []int{}
		query := "SELECT ms.player_id "
		query += "FROM matches_stats ms "
		query += "JOIN players p ON p.player_id = ms.player_id "
		query += "GROUP BY player_id "
		query += "ORDER BY SUM(ms.frags) DESC, AVG(ms.kdratio) DESC "
		
		rows, _ := db.Query(query)
		for rows.Next() {
			cur_pid := 0
			rows.Scan(&cur_pid)
			players_id = append(players_id, cur_pid)
		}
		
		//not enough players error
		//todo
		
		if dcount == 0  {
			for di := 0; di < nb_div; di++ {
				d_name := fmt.Sprintf("Division %d", di+1)
				div := csgodb.AddDivision(db, d_name)
				for pi := 0; pi < nb_player; pi++ {
					div.AddPlayer(db, players_id[pid])
					pid++
				}
			}
		} 
		
		db.Close()
		
		p.Content = "<p>Pool generated with success !</p>"
	} else if action == "settings" {
		save := r.FormValue("save")
		if save == "yes" {
			//persist setting and parse form here
			csgopool.Pool.Settings.PoolOn = ParseBool(r.FormValue("poolon"))
			csgopool.Pool.Settings.AutoAddMatches = ParseBool(r.FormValue("autoadd"))
			csgopool.Pool.Settings.SteamBot = ParseBool(r.FormValue("autoadd"))
			csgopool.Pool.Settings.PoolCost = ParseFloat(r.FormValue("poolcost"))
			csgopool.Pool.Settings.SteamKey = r.FormValue("steamkey")
			csgopool.Pool.Settings.MailVerification = ParseBool(r.FormValue("mailverification"))
			
			csgopool.Pool.Settings.Mail.Address = r.FormValue("address")
			csgopool.Pool.Settings.Mail.Port = ParseInt(r.FormValue("mailport"))
			csgopool.Pool.Settings.Mail.Username = r.FormValue("mailuser")
			csgopool.Pool.Settings.Mail.Password = r.FormValue("mailpassword")
			csgopool.Pool.Settings.Mail.Email = r.FormValue("email")
			
			csgopool.Pool.SaveSetting(csgopool.Pool.Path)
			
			p.Content = template.HTML(`<h4>Settings saved!</h4>`)
			
		} else {
			content := ReadFile("adminpoolsettings.html")
			content = strings.Replace(content, "{{.PoolOn}}", BoolToString(csgopool.Pool.Settings.PoolOn), 1)
			content = strings.Replace(content, "{{.SteamKey}}", csgopool.Pool.Settings.SteamKey, 1)
			content = strings.Replace(content, "{{.SteamBot}}", BoolToString(csgopool.Pool.Settings.SteamBot), 1)
			content = strings.Replace(content, "{{.AutoAdd}}", BoolToString(csgopool.Pool.Settings.AutoAddMatches), 1)
			content = strings.Replace(content, "{{.PoolCost}}", FloatToString(csgopool.Pool.Settings.PoolCost), 1)
			content = strings.Replace(content, "{{.MailVerification}}", BoolToString(csgopool.Pool.Settings.MailVerification), 1)
			content = strings.Replace(content, "{{.MailAddress}}", csgopool.Pool.Settings.Mail.Address, 1)
			content = strings.Replace(content, "{{.MailPort}}", fmt.Sprintf("%d", csgopool.Pool.Settings.Mail.Port), 1)
			content = strings.Replace(content, "{{.MailUser}}", csgopool.Pool.Settings.Mail.Username, 1)
			content = strings.Replace(content, "{{.MailPassword}}", csgopool.Pool.Settings.Mail.Password, 1)
			content = strings.Replace(content, "{{.Email}}", csgopool.Pool.Settings.Mail.Email, 1)
			p.Content = template.HTML(content)
		}
		
		
	} else if action == "matches" {
		db, _ := csgodb.Db.Open()
		content := `<div class="row">
		<table class="table table-striped">
			<thead>
				<tr>
					<th>Id</th>
					<th>Date</th>
					<th>Team 1</th>
					<th>Team 2</th>
					<th>Event Id</th>
					<th>Pool Status</th>
				</tr>
			</thead>
			<tbody>
				%s
			</tbody>
		</table>
		</div>
		`
		
		matches := csgodb.GetAllMatches(db)
		db.Close()
		
		rows_html := ""
		
		for _, m := range matches {
			
			matchDate := fmt.Sprintf("%d-%02d-%02d", m.Date.Year(), m.Date.Month(), m.Date.Day())
			matchLink := &Link{Caption: matchDate, Url: "/viewmatch/"}
			matchLink.AddInt("id", m.MatchId)
			
			t1Link := &Link{Caption: fmt.Sprintf("%s (%d)", m.Team1.Name, m.Team1.Score), Url: "/viewteam/"}
			t1Link.AddInt("id", m.Team1.TeamId)
			
			t2Link := &Link{Caption: fmt.Sprintf("%s (%d)",m.Team2.Name, m.Team2.Score), Url: "/viewteam/"}
			t2Link.AddInt("id", m.Team2.TeamId)
			
			status := "none"
			
			if m.PoolStatus == 0 {
				pLink := &Link{Caption: "Add to pool", Url:"/adminpool/"}
				pLink.AddParameter("action", "poolmatch")
				pLink.AddInt("id", m.MatchId)
				
				status = fmt.Sprintf("Not Pooled <br /> %s", pLink.GetHTML())
			} else {
				pLink := &Link{Caption: "Revoke from pool", Url:"/adminpool/"}
				pLink.AddParameter("action", "revokematch")
				pLink.AddInt("id", m.MatchId)
				
				status = fmt.Sprintf("Pooled <br /> %s", pLink.GetHTML())
			}
			
			rows_html += fmt.Sprintf(`
				<tr>
					<td>%d</td>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
					<td>%d</td>
					<td>%s</td>
				</tr>
			`, 
			m.MatchId, matchLink.GetHTML(), t1Link.GetHTML(), t2Link.GetHTML(), m.EventId, status)
		}
		
		content = fmt.Sprintf(content, rows_html)
		p.Content = template.HTML(content)
	} else if action == "mergeplayer" {
		p.Content = template.HTML(ReadFile("adminmergeplayer.html"))
	} else if action == "merge" {
		
		str_playerId := r.FormValue("playerid")
		str_mergerId := r.FormValue("mergerid")
		
		_playerId, _ := strconv.ParseInt(str_playerId, 10, 32)
		_mergerId, _ := strconv.ParseInt(str_mergerId, 10, 32)
		
		db, _ := csgodb.Db.Open()
		
		rs := csgodb.MergePlayer(db, int(_playerId), int(_mergerId))
		
		db.Close()
		
		if rs {
			p.Content = template.HTML(`<h4>Merge completed</h4>`)
		} else {
			p.Content = template.HTML(`<h4>Merge Error !</h4>`)
		}
	} else if action == "addnews" {
		p.Content = template.HTML(ReadFile("adminaddnews.html"))
	} else if action == "postnews" {
		title := r.FormValue("title")
		text := r.FormValue("text")
		authorId := session.UserId
		
		news := &csgodb.News{}
		news.Title = title
		news.Text = text
		news.Author.AuthorId = authorId
		
		db, _ := csgodb.Db.Open()
		
		news.Insert(db)
		
		db.Close()
		
		p.Content = template.HTML(`<h4>News posted !</h4>`)
	} else if action == "steamlogin" {
		
		
	} else if action == "logincheck" {
		
	}

	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Pool Administration"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Message = template.HTML(msgHtml)
	t.Execute(w, p)
}
