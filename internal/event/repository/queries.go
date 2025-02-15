package repository

const (
	selectEventQuery          = `select * from events where id = $1 limit 1;`
	selectUserIDsByEventQuery = `select user_id from users_and_events ue where ue.event_id = $1;`
	selectEventsQuery         = `select events.*, count(*) over () as total_count from events limit $1 offset $2;`
	selectEventsByUserQuery   = `
        select e.*, count(*) over () as total_count 
        from events e join users_and_events ue on ue.user_id = $1 and e.id = ue.event_id
        limit $2
        offset $3;
    `
	insertEventQuery      = `insert into events(name, timestamp) values (:name, :timestamp) returning id;`
	insertEventUserQuery  = `insert into users_and_events(user_id, event_id) values (:user_id, :event_id);`
	updateEventQuery      = `update events set name = :name, timestamp = :timestamp where id = :id;`
	deleteEventQuery      = `delete from events where id = $1;`
	deleteEventUsersQuery = `delete from users_and_events where event_id = $1 and user_id in ($2);`
)
