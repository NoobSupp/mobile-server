-- SQLite
UPDATE users SET is_admin = 1 WHERE lower(trim(username)) = 'admin';
