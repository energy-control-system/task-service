-- +goose Up
create table if not exists task_statuses
(
    id   int primary key generated always as identity,
    name text not null
);

insert into task_statuses (name)
values ('Planned'),
       ('InWork'),
       ('Done');

create table if not exists tasks
(
    id            int primary key generated always as identity,
    brigade_id    int,
    object_id     int         not null,
    plan_visit_at timestamptz, -- Запланированное время визита
    status        int         not null references task_statuses (id) on delete restrict,
    comment       text,
    started_at    timestamptz,
    finished_at   timestamptz,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),
    check ( started_at is null or finished_at is null or started_at <= finished_at )
);

create index if not exists idx_tasks_brigade on tasks (brigade_id);
create index if not exists idx_tasks_object on tasks (object_id);
create index if not exists idx_tasks_plan_visit_at on tasks (plan_visit_at);
create index if not exists idx_tasks_status on tasks (status);

-- +goose StatementBegin
create or replace function update_updated_at_column()
    returns trigger as
$$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

create trigger trg_tasks_updated_at
    before update
    on tasks
    for each row
execute function update_updated_at_column();

-- +goose Down
drop trigger if exists trg_tasks_updated_at on tasks;
drop function if exists update_updated_at_column();
drop table if exists tasks;
drop table if exists task_statuses;
