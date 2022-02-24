CREATE TABLE IF NOT EXISTS users (
  username TEXT NOT NULL,
  created_at INTEGER DEFAULT 0,
  config_name TEXT,
  config_version TEXT,

  FOREIGN KEY(config_version, config_name) REFERENCES configs(version, name)
);
