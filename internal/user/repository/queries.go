package repository

const (
	selectUserQuery         = `select * from users where id = $1 limit 1;`
	selectUsersQuery        = `select *, count(*) over () as total_count from users limit $1 offset $2;`
	selectUsersByEventQuery = `
        select u.*, count(*) over () as total_count 
        from users u join users_and_events ue on ue.event_id = $1 and u.id = ue.user_id
        limit $2
        offset $3;
    `
	insertUserQuery = `insert into users(name) values (:name) returning id;`
	updateUserQuery = `update users set name = :name where id = :id;`
	deleteUserQuery = `delete from users where id = $1;`
)
