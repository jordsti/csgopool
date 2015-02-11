package csgopoolweb

type SessionField struct {
	Name string
	Value string
}

type Session struct {
	Id string
	UserId int
	Fields []SessionField
}

type SessionContainer struct {
	Sessions []Session
}

func (sc *SessionContainer) NewSession(userId int) *Session {
	//new session id
	
	sess := Session{UserId: userId}
	
	sc.Sessions = append(sc.Sessions, sess)
	
	return &sess
}

func (s *Session) GetField(name string) string {
	
	for _, f := range s.Fields {
		if f.Name == name {
			return f.Value
		}
	}
	
	return ""
}

func (s *Session) AddField(name string, value string) {
	
	field := SessionField{Name: name, Value: value}
	s.Fields = append(s.Fields, field)
}