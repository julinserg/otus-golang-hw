-- +goose Up
CREATE table events (
    id              text primary key,
    title           text not null,
    time_start      timestamp not null,
    time_stop       timestamp not null,
    description     text,
    user_id         text not null,    
    time_notify     interval
);

-- +goose Down
drop table events;
