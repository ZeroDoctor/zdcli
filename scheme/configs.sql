CREATE TABLE IF NOT EXISTS configs (
  name TEXT NOT NULL,
  editor TEXT DEFAULT '',
  lua_exec TEXT DEFAULT '',
  lua_dir TEXT DEFAULT '',
  created_at INTEGER DEFAULT 0,
  version TEXT DEFAULT '',
  
  PRIMARY KEY(name, version)
);
