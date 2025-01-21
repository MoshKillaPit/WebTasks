package repositories

const (
	CreateUserQuery = `
	INSERT INTO public.users (name, key) 
	VALUES (:name, :key) 
	RETURNING id, name, key;`

	GetUserByIDQuery = `
	SELECT id, name, key
	FROM public.users
	WHERE id = $1;`

	GetAllUsersQuery = `
	SELECT id, name, key
	FROM public.users;`

	UpdateUserQuery = `
	UPDATE public.users
	SET name = :name, key = :key
	WHERE id = :id
	RETURNING id, name, key;`

	DeleteUserQuery = `
	DELETE FROM public.users
	WHERE id = $1;`
)
