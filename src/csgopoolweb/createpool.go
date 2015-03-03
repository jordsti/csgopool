package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
	"csgopool"
	"strconv"
)

type CreatePoolPage struct {
	Page
	Form template.HTML
}

func CreatePoolHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	if !session.IsLogged() {
		http.Redirect(w, r, "/", 301)
	}

	t, err := MakeTemplate("usercreatepool.html")
	if err != nil {
		fmt.Println(err)
	}
	
	m := GetMenu(session)
	
	db, _ := csgodb.Db.Open()
	
	action := r.FormValue("action")
	
	pools := csgodb.GetMetaPoolsByUser(db, session.UserId)
	
	form_html := ""
	divisions := csgodb.GetAllDivisionsWithPlayer(db)
	credit := csgodb.GetCreditByUser(db, session.UserId)
	if len(action) == 0 || action == "form" {
		
		if len(pools) == 0 {
			
			form_html = `<form method="POST" action="/createpool/?action=submit">`
			
			inner_html := ""
			it := 1
			for _, div := range divisions {
				
				inner_html += `<div class="col-sm-2">`
				inner_html += fmt.Sprintf(`<h4>%s</h4>`, div.Name)
				
				id_div := fmt.Sprintf(`division_%d`, div.DivisionId)
				
				for i, pl := range div.Players {
					checked := ""
					if i == 0 {
						checked = "checked"
					}
					
					playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/", Target:"_blank"}
					playerLink.AddInt("id", pl.PlayerId)
					
					inner_html += fmt.Sprintf(`<div class="radio"><label><input type="radio" name="%s" id="%s" value="%d" %s/>%s</label></div>`, id_div, id_div, pl.PlayerId, checked, playerLink.GetHTML())
				}
				
				inner_html += `</div>`
				
				if it != 0 && it % 3 == 0 {
					form_html += fmt.Sprintf(`<div class="row">%s</div>`, inner_html)
					inner_html = ""
				}
				
				it++
			}
			
			if len(inner_html) > 0 {
				form_html += fmt.Sprintf(`<div class="row">%s</div>`, inner_html)
			}
			
			form_html += `<button type="submit" class="btn btn-default">Create my pool</button>`
			form_html += `</form>`
			
		} else {
			form_html = `<h4>You already got a pool !</h4>`
		}
	
	} else if action == "submit" {
		if len(pools) == 0 {
			choices := []*csgodb.UserPool{}
			//todo
			//pool creation with form data
			
			for _, div := range divisions {
				d_name := fmt.Sprintf(`division_%d`, div.DivisionId)
				_d_value := r.FormValue(d_name)
				
				tmp, _ := strconv.ParseInt(_d_value, 10, 32)
				d_value := int(tmp)
				
				if div.IsPlayerIn(d_value) {
					
					choice := &csgodb.UserPool{}
					choice.UserId = session.UserId
					choice.DivisionId = div.DivisionId
					choice.PlayerId = d_value
					choices = append(choices, choice)
					
				} else {
					state.Log.Info(fmt.Sprintf("User [%d] try to insert a player [%d] that doesn't belong to this division [%d]", session.UserId, d_value, div.DivisionId))
				}
				
				if credit.Amount >= csgopool.Pool.Settings.PoolCost && credit.UserId != 0 {
					
					credit.Substract(csgopool.Pool.Settings.PoolCost)
					credit.UpdateCredit(db)
					
					if len(choices) == len(divisions) {
					csgodb.InsertPoolChoices(db, choices)
					form_html = `<h4>Pool submitted with success!</h4>`
					} else {
						form_html = `<h4>Incorrect pool choice</h4>`
					}
				} else {
					form_html = `<h3>You need more credit to enter the pool</h3>`
					form_html += fmt.Sprintf(`<p>Price is <strong>%.2f</strong>, and you have <strong></strong>%.2f</p>`, csgopool.Pool.Settings.PoolCost, credit.Amount)
				}
				
			}
			
		}
	}
	
	db.Close()
	
	p := &CreatePoolPage{}
	p.Title = "CS:GO Pool - Home"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)
	p.Form = template.HTML(form_html)
	t.Execute(w, p)
	
}

