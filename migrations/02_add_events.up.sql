create table if not exists events (
    id integer primary key autoincrement,
    name text not null,
    timestamp timestamp not null
);

create table if not exists users_and_events (
    user_id integer,
    event_id integer,
    foreign key (user_id) references users(id) on update cascade on delete cascade,
    foreign key (event_id) references events(id) on update cascade on delete cascade,
    unique (user_id, event_id)
);
