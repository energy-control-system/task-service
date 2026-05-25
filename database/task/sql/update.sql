update tasks
set plan_visit_at = :plan_visit_at,
    comment       = :comment
where id = :id
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
