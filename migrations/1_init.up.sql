//запросы для создания таблицы

CREATE TABLE IF NOT EXISTS users(
    uid VARCHAR(36) PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL, 
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS email_id ON Users (email);

CREATE TABLE IF NOT EXISTS books(
    bid VARCHAR(36) PRIMARY KEY,
    lable TEXT NOT NULL,
    author TEXT NOT NULL,
    uid VARCHAR(36) NOT NULL, 
);

