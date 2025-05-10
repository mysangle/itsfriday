-- migration_history
CREATE TABLE IF NOT EXISTS migration_history (
  version TEXT NOT NULL PRIMARY KEY,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- user
CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',
  username TEXT NOT NULL UNIQUE,
  role TEXT NOT NULL CHECK (role IN ('HOST', 'ADMIN', 'USER')) DEFAULT 'USER',
  email TEXT NOT NULL,
  nickname TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  avatar_url TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_user_username ON user (username);

-- user_setting
CREATE TABLE IF NOT EXISTS user_setting (
  user_id INTEGER NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  UNIQUE(user_id, key)
);

-- libro service --

-- book
CREATE TABLE IF NOT EXISTS book (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  translator TEXT NOT NULL DEFAULT '',
  pages INTEGER NOT NULL,
  pub_year INTEGER NOT NULL,
  genre TEXT NOT NULL DEFAULT '',
  UNIQUE(title, author)
);

-- book_review
CREATE TABLE IF NOT EXISTS book_review (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  user_id INTEGER NOT NULL,
  book_id INTEGER NOT NULL,
  date_read TEXT NOT NULL CHECK (length(date_read) = 10 AND substr(date_read, 5, 1) = '-' AND substr(date_read, 8, 1) = '-'), -- YYYY-MM-DD
  rating REAL NOT NULL CHECK (rating >= 0 AND rating <= 5),
  review TEXT NOT NULL DEFAULT '',
  public INTEGER NOT NULL DEFAULT 0,
  UNIQUE(user_id, book_id, date_read)
);

CREATE INDEX IF NOT EXISTS idx_book_review_user_id_date_read ON book_review (user_id, date_read);
CREATE INDEX IF NOT EXISTS idx_book_review_book_id_date_read ON book_review (book_id, date_read);

-- monero service --

-- category
CREATE TABLE IF NOT EXISTS expense_category (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  priority INTEGER NOT NULL DEFAULT 1,
  UNIQUE(user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_expense_category_user_id ON expense_category (user_id);

-- expense
CREATE TABLE IF NOT EXISTS expense (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  user_id INTEGER NOT NULL,
  category_id INTEGER NOT NULL,
  date_used TEXT NOT NULL CHECK (length(date_used) = 10 AND substr(date_used, 5, 1) = '-' AND substr(date_used, 8, 1) = '-'), -- YYYY-MM-DD
  item TEXT NOT NULL,
  price INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_expense_user_id ON expense (user_id);
CREATE INDEX IF NOT EXISTS idx_expense_user_id_category_id_date_used ON expense (user_id, category_id, date_used);
