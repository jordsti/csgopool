package csgopoolweb

import (
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
	"net/http"
	"strconv"
	"fmt"
)

type SessionField struct {
	Name string
	Value string
}

type Session struct {
	Id string
	UserId int
	Fields []*SessionField
}

type SessionContainer struct {
	Sessions []*Session
}

func (sc *SessionContainer) NewSession(userId int) *Session {
	//new session id
	
	sess := &Session{UserId: userId, Id:GenerateSessionKey()}
	
	sc.Sessions = append(sc.Sessions, sess)
	
	return sess
}

func (sc *SessionContainer) GetSession(id string) *Session {
	
	for _, sess := range sc.Sessions {
		if sess.Id == id {
			return sess
		}
	}
	
	return nil
}

func (s *Session) ClearFields() {
	s.Fields = []*SessionField{}
}

func (s *Session) GetField(name string) string {
	
	for _, f := range s.Fields {
		if f.Name == name {
			return f.Value
		}
	}
	
	return ""
}

func (s *Session) GetInt(name string) int {
	field := s.GetFieldPtr(name)
	if field != nil {
		return field.Int()
	} else {
		return -1
	}
	
}

func (f *SessionField) Int() int {
	_i, _ := strconv.ParseInt(f.Value, 10, 32)
	return int(_i)
}

func (s *Session) GetFieldPtr(name string) *SessionField {
	
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	
	return nil
}

func (s *Session) IsFieldExists(name string) bool {
	for _, f := range s.Fields {
		if f.Name == name {
			return true
		}
	}
	
	return false
}

func (s *Session) IsLogged() bool {
	if s.UserId == 0 {
		return false
	}
	
	return true
}

func (s *Session) AddField(name string, value string) {
	
	field := &SessionField{Name: name, Value: value}
	s.Fields = append(s.Fields, field)
}

func (s *Session) SetField(name string, value string) {
	
	field := s.GetFieldPtr(name)
	
	if field == nil {
		s.AddField(name, value)
	} else {
		field.Value = value
	}
	
}

func (s *Session) RemoveField(name string) {
	fields := s.Fields
	
	s.Fields = []*SessionField{}
	
	for _, f := range fields {
		if f.Name != name {
			s.Fields = append(s.Fields, f)
		}
	}
}

func (s *Session) PopField(name string) *SessionField {
	
	field := s.GetFieldPtr(name)
	
	if field != nil {
		s.RemoveField(name)
	}
	
	return field
	
}

func GenerateSessionKey() string {
	
   size := 32 // change the length of the generated random string here

   rb := make([]byte,size)
   _, err := rand.Read(rb)


   if err != nil {
      state.Log.Error(fmt.Sprintf("%s", err))
   }
   
   hasher := sha256.New()
   hasher.Write(rb)

   rs := hex.EncodeToString(hasher.Sum(nil))
   
   return rs
}

func (ws *WebServerState) HandleSession(w http.ResponseWriter, r *http.Request) *Session {
	
	state.Log.Info(fmt.Sprintf("Serve %s for %s", r.RequestURI ,r.RemoteAddr))
  
	for _, cookie := range r.Cookies() {
		if cookie.Name == "csgopool" {
			//fmt.Printf("Session found! [%s]\n", cookie.Value)
			
			sessId := cookie.Value
			
			sess := ws.Sessions.GetSession(sessId)
			
			if sess != nil {
				return sess
			}
		}
	}
	

	//fmt.Println("No cookie, creating a session")
	
	sess := state.Sessions.NewSession(0)
	
	c := &http.Cookie{}
	c.Name = "csgopool"
	c.Value = sess.Id
	c.Path = "/"
	//c.Domain = ws.Domain
	
	http.SetCookie(w, c)
	return sess
	
	
}