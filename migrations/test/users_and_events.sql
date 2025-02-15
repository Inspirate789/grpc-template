insert into users(name) values ('aboba1');
insert into users(name) values ('aboba2');
insert into users(name) values ('aboba3');
insert into users(name) values ('aboba4');
insert into users(name) values ('aboba5');

insert into events(name, timestamp) values ('event1', current_timestamp);
insert into events(name, timestamp) values ('event2', current_timestamp);
insert into events(name, timestamp) values ('event3', current_timestamp);
insert into events(name, timestamp) values ('event4', current_timestamp);
insert into events(name, timestamp) values ('event5', current_timestamp);

insert into users_and_events(user_id, event_id) values (1, 1);
insert into users_and_events(user_id, event_id) values (1, 2);
insert into users_and_events(user_id, event_id) values (2, 2);
insert into users_and_events(user_id, event_id) values (3, 2);
insert into users_and_events(user_id, event_id) values (5, 4);
insert into users_and_events(user_id, event_id) values (2, 4);
insert into users_and_events(user_id, event_id) values (4, 3);
