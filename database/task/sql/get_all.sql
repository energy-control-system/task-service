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
where ($1::integer is null or status = $1)
  and ($2::timestamptz is null or plan_visit_at >= $2)
  and ($3::timestamptz is null or plan_visit_at <= $3)
order by
    case when $4 = 'asc' then plan_visit_at end asc nulls last,
    case when $4 = 'desc' then plan_visit_at end desc nulls last,
    id
limit $5 offset $6;
