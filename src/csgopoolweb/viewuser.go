package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
	//"csgoscrapper"
)

type ViewUserPage struct {
	Page
	UserName string
	Pools template.HTML
	Points string
}

func ViewUserHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("viewuser.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	db, _ := csgodb.Db.Open()
	
	_m_id := r.FormValue("id")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)

	userId := int(m_id)
	
	user := csgodb.GetUserById(db, userId)
	points := csgodb.GetUserPoint(db, userId)
	divs_html := ""
	if user.Id != 0 {
		
		pools := csgodb.GetMetaPoolsByUser(db, user.Id)

		if len(pools) == 0 {

			divs_html = fmt.Sprintf(`<h4>No pools for this user</h4>`)
		} else {
			
			row_html := `<div class="row">%s</div>`

			for _, pool := range pools {
				
				
				inner_html := fmt.Sprintf(`<div class="col-md-2">%s</div>`, pool.Division.Name)
				
				playerLink := &Link{Caption: pool.Player.Name, Url:"/viewplayer/"}
				playerLink.AddInt("id", pool.Player.PlayerId)
				
				inner_html += fmt.Sprintf(`<div class="col-md-2">%s</div>`, playerLink.GetHTML())
				inner_html += fmt.Sprintf(`<div class="col-md-2">%d</div>`, pool.Points)
				
				
				divs_html += fmt.Sprintf(row_html, inner_html)
			}
		}
		
	
	}
	db.Close()

	p := &ViewUserPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - User : %s", user.Name)
	p.UserName = user.Name
	p.Points = fmt.Sprintf("%d", points.Points)
	p.Pools = template.HTML(divs_html)
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.GenerateRightSide(session)
	t.Execute(w, p)


}


	
	
	