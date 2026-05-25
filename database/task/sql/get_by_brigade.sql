select id,
       brigade_id,
       object_id,
       plan_visit_at,
       status,
       comment,
       started_at,
       finished_at,
       created_at,
       updated_at
from tasks
where brigade_id = $1
  and ($2::integer is null or status = $2)
  and ($3::timestamptz is null or plan_visit_at >= $3)
  and ($4::timestamptz is null or plan_visit_at <= $4)
order by id
limit $5 offset $6;
