CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT(8) NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT(128)
);

CREATE INDEX idx_scheduler_date on scheduler(date);