package logger

import (
  "fmt"
  "os"
)

type LoggerState struct {
  Level int
  LogPath string
}

//level 0 -> 
//level 1 -> info
//level 2 -> debug
//level 3 -> paranoid


func (l *LoggerState) Error(message string) {

   if l.Level >= 0 {
    msg := fmt.Sprintf("[ERROR] : %s\n", message)
    fmt.Printf(msg)
    l.AppendFile(msg)
   }
}

func (l *LoggerState) Info(message string) {

   if l.Level >= 1 {
    msg := fmt.Sprintf("[INFO] : %s\n", message)
    fmt.Printf(msg)
    l.AppendFile(msg)
   }
}

func (l *LoggerState) Debug(message string) {
  if l.Level >= 2 {
   msg := fmt.Sprintf("[DEBUG] : %s\n", message)
   fmt.Printf(msg)
   l.AppendFile(msg)
  }  
}

func (l *LoggerState) Log(message string) {
    if l.Level >= 3 {
   msg := fmt.Sprintf("[OUT] : %s\n", message)
   fmt.Printf(msg)
   l.AppendFile(msg)
  }  
}

func (l *LoggerState) AppendFile(message string) {
 file, err := os.OpenFile(l.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
 if err != nil {
    panic(err)
  }
  
 defer file.Close()
 
 if _, err = file.WriteString(message); err != nil {
    panic(err)
  }
}