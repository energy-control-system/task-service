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
order by id
limit $1 offset $2;
