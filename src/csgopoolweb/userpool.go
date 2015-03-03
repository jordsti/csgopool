package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
	"csgopool"
)

type UserPoolPage struct {
	Page
	Pools template.HTML
}

func UserPoolHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}

	t, err := MakeTemplate("userpool.html")
	if err != nil {
		fmt.Println(err)
	}
	
	m := GetMenu(session)
	
	
	db, _ := csgodb.Db.Open()
	
	pools := csgodb.GetMetaPoolsByUserWithoutPoints(db, session.UserId)
	credit := csgodb.GetCreditByUser(db, session.UserId)
	divs_html := ""
	
	if len(pools) == 0 {
		
		if credit.Amount < csgopool.Pool.Settings.PoolCost {
			divs_html = `<h3>You need more credit to enter the pool</h3>`
			divs_html += fmt.Sprintf(`<p>Price is <strong>%.2f</strong>, and you have <strong>%.2f</strong></p>`, csgopool.Pool.Settings.PoolCost, credit.Amount)
		} else {
			createLink := &Link{Caption:"Create your pool", Url:"/createpool/"}
			createLink.AddParameter("action", "form")
			divs_html = fmt.Sprintf(`<h4>No pool found for you !</h4><p>You need to create your pool, %s <br /> Pool fee is <strong>%.2f</strong></p>`, createLink.GetHTML(), csgopool.Pool.Settings.PoolCost)
		}
	} else {
		
		nb_div := 1
		row_html := `<div class="row">%s</div>`
		inner_html := ""
		for _, pool := range pools {
			
			
			inner_html += fmt.Sprintf(`<div class="col-md-2"><h4>%s</h4><ul>`, pool.Division.Name)
			
			playerLink := &Link{Caption: pool.Player.Name, Url:"/viewplayer/"}
			playerLink.AddInt("id", pool.Player.PlayerId)
			
			inner_html += fmt.Sprintf(`<li>%s</li>`, playerLink.GetHTML())
			
			
			inner_html += `</ul></div>`
			
			if nb_div % 3 == 0 && nb_div > 0 {
				divs_html += fmt.Sprintf(row_html, inner_html)
				inner_html = ""
			}
			
			nb_div++
		}
		
		if len(inner_html) > 0 {
			divs_html += fmt.Sprintf(row_html, inner_html)
		}
	}
	
	db.Close()
	
	p := &UserPoolPage{}
	p.Title = "CS:GO Pool - Home"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)
	p.Pools = template.HTML(divs_html)
	t.Execute(w, p)
	
}