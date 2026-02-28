package sessions

const create = "insert into sessions (id, user_id, refresh_hash, expires_at) values (?, ?, ?, ?);"

const findOne = "select id, user_id, refresh_hash, created_at, updated_at, expires_at from sessions where refresh_hash = ?;"

const update = "update sessions set refresh_hash = ?, expires_at = ? where id = ?;"
