package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"csgodb"
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
	} else if action == "submitdiv" {
		_nb_div := r.FormValue("nbdiv")
		_tmp, _ := strconv.ParseInt(_nb_div, 10, 32)
		nb_div := int(_tmp)
		
		div_html := "<h4>Divisions</h4>"
		div_html += `<div>`
		
		db, _ := csgodb.Db.Open()
		csgodb.ClearPool(db)
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
		
		db.Close()
		p.Content = template.HTML(div_html)
		
		
	}

	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Pool Administration"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Message = template.HTML(msgHtml)
	t.Execute(w, p)
}