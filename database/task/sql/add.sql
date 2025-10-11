insert into tasks (brigade_id, object_id, plan_visit_at, status, comment, started_at, finished_at)
values (:brigade_id, :object_id, :plan_visit_at, :status, :comment, :started_at, :finished_at)
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
