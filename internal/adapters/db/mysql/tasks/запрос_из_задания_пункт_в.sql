-- запрос сделал, но не придумал, как его использовать)

select
    t.id,
    t.team_id,
    t.title,
    t.description,
    t.status,
    t.created_at,
    t.updated_at,

    owner.id,
    owner.username,
    owner.email,

    assignee.id,
    assignee.username,
    assignee.email
from tasks t
join users owner on owner.id = t.created_by
left join users assignee on assignee.id = t.assignee_id
where t.assignee_id is not null
  and not exists (
      select 1
      from team_members tm
      where tm.team_id = t.team_id
        and tm.user_id = t.assignee_id
  )
order by t.created_at desc;
