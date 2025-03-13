CREATE TABLE userss (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(255)
);

CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    users_id INT NOT NULL,
    FOREIGN KEY (users_id) REFERENCES userss(id)
);
