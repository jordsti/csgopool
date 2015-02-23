package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os/exec"
)

type EnvVar struct {
	Name string
	Value string
}

type Environment struct {
	Vars []*EnvVar
}


func (ev *EnvVar) String() string {
	return fmt.Sprintf("%s=%s", ev.Name, ev.Value)
}

func (e *Environment) Push(name string, value string) {
	
	ev := &EnvVar{Name: name, Value:value}
	e.Vars = append(e.Vars, ev)
}

func (e *Environment) ApplyDefaults() {
	e.Push("GOROOT", "goroot")
}

func (e *Environment) ApplyToCommand(c *exec.Cmd) {
	
	for _, v := range e.Vars {
		str_env := v.String()
		c.Env = append(c.Env, str_env)
	}
	
}

func (e *Environment) Save(path string) {
	
	b, err := json.MarshalIndent(e, "", "	")
	
	if err != nil {
		fmt.Println("Error while encondig envs file [1]")
	}
	
	err = ioutil.WriteFile(path, []byte(b), 0644)
	
	if err != nil {
		fmt.Println("Error while writing envs file [1]")
	}
}

func (e *Environment) Load(path string) {
	
	b, err := ioutil.ReadFile(path)
	
	if err != nil {
		fmt.Println("Error while reading envs file [1]")
		fmt.Println("Saving a default envs config")
		e.ApplyDefaults()
		e.Save(path)
	}
	
	err = json.Unmarshal(b, e)
	
	if err != nil {
		fmt.Println("Error while parsing envs  file [1]")
	}
	
}


