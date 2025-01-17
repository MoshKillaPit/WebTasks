package models

type User struct {
	ID    int
	Name  string
	Tasks []Task // Слайс из структуры задачи. Куча задач будут в виде слайсов
}

type Task struct {
	ID     int
	Name   string
	Status string
}
