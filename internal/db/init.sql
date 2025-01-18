CREATE DATABASE postgres;

\c postgres;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       key VARCHAR(255) NOT NULL
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       status VARCHAR(50) NOT NULL,
                       time TIMESTAMP NOT NULL,
                       due TIMESTAMP NOT NULL,
                       user_id INT NOT NULL,
                       CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_id ON tasks (user_id);
CREATE INDEX idx_task_status ON tasks (status);

INSERT INTO users (name, key) VALUES
                                  ('John Doe', 'abc123'),
                                  ('Jane Smith', 'def456');

INSERT INTO tasks (name, status, time, due, user_id) VALUES
                                                         ('Task 1', 'In Progress', NOW(), NOW() + INTERVAL '2 days', 1),
                                                         ('Task 2', 'Completed', NOW(), NOW() + INTERVAL '5 days', 1),
                                                         ('Task 3', 'Pending', NOW(), NOW() + INTERVAL '3 days', 2);
