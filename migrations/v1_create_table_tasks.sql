
CREATE TYPE taskstatus AS ENUM('NotStarted','InProgress','Completed');

CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) UNIQUE,
    description TEXT,
    status taskstatus,
    createdat  TIMESTAMP DEFAULT current_timestamp,
    updatedat TIMESTAMP DEFAULT current_timestamp
);