create table if not exists users (
  id char(36) primary key default (uuid()),
  email varchar(255) not null unique
);