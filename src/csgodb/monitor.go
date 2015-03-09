package csgodb

import (
	"time"
	"database/sql"
)

const (
	DoneTask = 2
	OngoingTask = 1
	PendingTask = 0
)

type MonitorTask struct {
	TaskId int
	Type int
	Status int
	Date time.Time
}

func PushTask(db *sql.DB, taskType int, status int) {
	
	now := time.Now()
	
	query := "INSERT INTO (task_type, task_status, task_date) VALUES (?, ?, ?)"
	db.Exec(query, taskType, status, now)
}

func (mt *MonitorTask) Update(db *sql.DB) {
	query := "UPDATE tasks SET task_status = ? WHERE task_id = ?"
	db.Exec(query, mt.Status, mt.TaskId)
}

func (mt *MonitorTask) IsDone() bool {
	return (mt.Status == DoneTask)
}

func (mt *MonitorTask) IsOngoing() bool {
	return (mt.Status == OngoingTask)
}

func (mt *MonitorTask) IsPending() bool {
	return (mt.Status == PendingTask)
}

func GetUndoneTasks(db *sql.DB) []*MonitorTask {
	tasks := []*MonitorTask{}
	query := "SELECT task_id, task_type, task_status, task_date FROM tasks WHERE task_status = 0 ORDER BY task_id DESC"
	
	rows, _ := db.Query(query)
	for rows.Next() {
		
		task := &MonitorTask{}
		rows.Scan(&task.TaskId, &task.Type, &task.Status, &task.Date)
		
		tasks = append(tasks, task)
	}
	
	return tasks
}

func GetAllTasks(db *sql.DB) []*MonitorTask {
	tasks := []*MonitorTask{}
	query := "SELECT task_id, task_type, task_status, task_date FROM tasks ORDER BY task_id DESC"
	
	rows, _ := db.Query(query)
	for rows.Next() {
		
		task := &MonitorTask{}
		rows.Scan(&task.TaskId, &task.Type, &task.Status, &task.Date)
		
		tasks = append(tasks, task)
	}
	
	return tasks
}