CREATE TABLE IF NOT EXISTS configs (
  username TEXT DEFAULT '',
  editor TEXT DEFAULT '',
  lua_exec TEXT DEFAULT '',
  lua_dir TEXT DEFAULT '',
  created_at INTEGER 0,
  version TEXT DEFAULT ''
);
