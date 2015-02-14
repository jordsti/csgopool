package csgopool

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
)

type UserConstraint struct {
	NameMax int
	NameMin int
	PasswordMin int
}

type User struct {
	Id int
	Name string
	Password string
	Email string
}

type Users struct {
	Users []User
	CurrentId int
	Constraints UserConstraint
}

func (u *Users) Login(username string, password string) (*User, error) {
	
	user := u.GetUserByName(username)
	
	if user == nil {
		return nil, errors.New("Invalid username and/or password")
	}
	
	b, _ := hex.DecodeString(user.Password)
	
	err := bcrypt.CompareHashAndPassword(b, []byte(password))
	
	if err == nil {
		//success
		return user, nil
	} else {
		return nil, errors.New("Invalid username and/or password")
	}
	
}

func (u *Users) CreateUser(username string, password string, email string) (*User, error) {
	
	//username unique check
	user := u.GetUserByName(username)
	
	if user != nil {
		return nil, errors.New("Username already exists")
	}
	
	user = u.GetUserByEmail(email)
	
	if user != nil {
		return nil, errors.New("Email already exists")
	}
	
	u.CurrentId = u.CurrentId + 1
	
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	pw_str := hex.EncodeToString(pwd)
	
	nuser := User{Name: username, Password: pw_str, Email:email, Id:u.CurrentId}
	
	u.Users = append(u.Users, nuser)
	
	users_file := watcher.DataPath + "users.json"
	
	u.SaveUsers(users_file)
	
	return &nuser, nil
}

func (u *Users) GetUserByName(username string) *User {
	for _, user := range u.Users {
		if user.Name == username {
			return &user
		}
	}
	
	return nil
}

func (u *Users) GetUserByEmail(email string) *User {
		for _, user := range u.Users {
		if user.Email == email {
			return &user
		}
	}
	
	return nil
}

func (u *Users) SaveUsers(path string) {
	
	data, err := json.MarshalIndent(u.Users, "", "	")
	
	if err != nil {
		fmt.Println("Error while saving users [1]")
	}
	
	err = ioutil.WriteFile(path, data, 0644)
	
	if err != nil {
		fmt.Println("Error while saving users [2]")
	}
}

func (u *Users) LoadUsers(path string) {
	
	u.Users = []User{}
	u.CurrentId = 0
	
	//default constraint
	u.Constraints.PasswordMin = 6
	u.Constraints.NameMin = 4
	u.Constraints.NameMax = 12
	
	b, err := ioutil.ReadFile(path)
	
	if err != nil {
		fmt.Println("Error while reading users [1]")
	}
	
	err = json.Unmarshal(b, &u.Users)
	
	if err != nil {
		fmt.Println("Error while reading users [2]")
	}
	
	fmt.Printf("%d users loaded\n", len(u.Users))
	
	for _, user := range u.Users {
		if user.Id > u.CurrentId {
			u.CurrentId = user.Id
		}
	}
}