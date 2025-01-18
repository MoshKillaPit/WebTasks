package models

import "time"

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Key   string `db:"key"`
	Tasks []Task // Слайс из структуры задачи. Куча задач будут в виде слайсов для одного пользователя
}

type Task struct {
	ID     int       `db:"id" json:"id"`
	Name   string    `db:"name" json:"name"`
	Status string    `db:"status" json:"status"`
	Time   time.Time `db:"time" json:"time"`
	Due    time.Time `db:"due" json:"due"`
	UserID int       `db:"user_id" json:"user_id"`
}
