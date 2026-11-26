ALTER TABLE users ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN must_change_password INTEGER NOT NULL DEFAULT 0 CHECK (must_change_password IN (0,1));
ALTER TABLE users ADD COLUMN temp_password_expires_at TEXT;
ALTER TABLE users ADD COLUMN created_by_user_id TEXT REFERENCES users(id);
ALTER TABLE users ADD COLUMN auth_version INTEGER NOT NULL DEFAULT 1 CHECK (auth_version >= 1);

INSERT OR IGNORE INTO roles (id, name, description, created_at) VALUES
  ('role-admin', 'admin', 'Administrador do sistema', datetime('now')),
  ('role-gerente', 'gerente', 'Gerente operacional', datetime('now'));