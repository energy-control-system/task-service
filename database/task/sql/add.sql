insert into tasks (brigade_id, object_id, plan_visit_at, status, comment)
values (:brigade_id, :object_id, :plan_visit_at, 1, :comment)
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
