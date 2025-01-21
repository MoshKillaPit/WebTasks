package repositories

const (
	CreateTaskQuery = `
	INSERT INTO public.tasks (name, status, time, due, user_id) 
VALUES (:name, :status, :time, :due, :user_id) 
RETURNING id, name, status, time, due, user_id;`

	GetTaskByIDQuery = `
	SELECT id, name, status, time, due 
	FROM public.tasks 
	WHERE id = $1;`

	GetAllTasksQuery = `
	SELECT id, name, status, time, due 
	FROM public.tasks;`

	UpdateTaskQuery = `
	UPDATE public.tasks 
	SET name = :name, status = :status, time = :time, due = :due 
	WHERE id = :id 
	RETURNING id, name, status, time, due;`

	DeleteTaskQuery = `
	DELETE FROM public.tasks 
	WHERE id = $1;`
)
