create table if not exists users (
    id char(36) primary key,
    username varchar(100) not null unique,
    email varchar(100) not null unique,
    password_hash text not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);

create table if not exists sessions (
    id char(36) primary key,
    user_id char(36) not null,
    refresh_hash text unique not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    expires_at timestamp not null,
    foreign key (user_id) references users(id) on delete cascade
);

create table if not exists teams (
    id char(36) primary key,
    name varchar(100) not null unique,
    created_by char(36) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    foreign key (created_by) references users(id) on delete cascade
);

create table if not exists team_members (
    id char(36) primary key,
    team_id char(36) not null,
    user_id char(36) not null,
    role enum('admin','member', 'owner') not null,
    joined_at timestamp default current_timestamp,
    unique(team_id, user_id),
    foreign key (team_id) references teams(id) on delete cascade,
    foreign key (user_id) references users(id) on delete cascade
);

create table if not exists tasks (
    id char(36) primary key,
    team_id char(36) not null,
    title varchar(255) not null,
    description text,
    status enum('todo','in_progress','done') default 'todo',
    assignee_id char(36),
    created_by char(36) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    foreign key (team_id) references teams(id) on delete cascade,
    foreign key (assignee_id) references users(id) on delete set null,
    foreign key (created_by) references users(id) on delete cascade
);

create table if not exists task_history (
    id char(36) primary key,
    task_id char(36) not null,
    changed_by char(36) not null,
    snapshot json not null,
    changed_at timestamp default current_timestamp,
    foreign key (task_id) references tasks(id) on delete cascade,
    foreign key (changed_by) references users(id) on delete cascade
);

create table if not exists task_comments (
    id char(36) primary key,
    task_id char(36) not null,
    user_id char(36) not null,
    comment text not null,
    created_at timestamp default current_timestamp,
    foreign key (task_id) references tasks(id) on delete cascade,
    foreign key (user_id) references users(id) on delete cascade
);

create index if not exists idx_tasks_team_status on tasks(team_id, status);
create index if not exists idx_tasks_assignee on tasks(assignee_id);
create index if not exists idx_task_history_task_changed on task_history(task_id, changed_at);
create index if not exists idx_task_comments_task on task_comments(task_id);
create index if not exists idx_team_members_team_user on team_members(team_id, user_id);