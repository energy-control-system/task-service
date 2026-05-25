delete from tasks
where id = $1
returning id,
          brigade_id,
          object_id,
          plan_visit_at,
          status,
          comment,
          created_at,
          updated_at;
