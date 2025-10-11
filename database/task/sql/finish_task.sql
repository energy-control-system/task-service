update tasks
set status      = 3,
    finished_at = now()
where id = $1
returning id,
    brigade_id,
    object_id,
    plan_visit_at,
    status,
    comment,
    started_at,
    finished_at,
    created_at,
    updated_at;
