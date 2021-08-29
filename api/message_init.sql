CREATE TABLE IF NOT EXISTS message(
id INTEGER PRIMARY KEY AUTOINCREMENT,
post_id INTEGER REFERENCES post(id) ON DELETE CASCADE,
user TEXT,
email TEXT,
message TEXT,
inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX index_message_post_id_inserted ON message(post_id,inserted_at DESC);
