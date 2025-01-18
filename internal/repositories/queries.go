package repositories

const (
	CreateTaskQuery = `
	INSERT INTO tasks (name, status, time, due, user_id) 
VALUES (:name, :status, :time, :due, :user_id) 
RETURNING id, name, status, time, due, user_id;`

	GetTaskByIDQuery = `
	SELECT id, name, status, time, due 
	FROM tasks 
	WHERE id = $1;`

	GetAllTasksQuery = `
	SELECT id, name, status, time, due 
	FROM tasks;`

	UpdateTaskQuery = `
	UPDATE tasks 
	SET name = :name, status = :status, time = :time, due = :due 
	WHERE id = :id 
	RETURNING id, name, status, time, due;`

	DeleteTaskQuery = `
	DELETE FROM tasks 
	WHERE id = $1;`
)
