CREATE TABLE IF NOT EXISTS account (id string not null primary key);
CREATE TABLE IF NOT EXISTS user (id string not null primary key, account_id string not null,
    email string not null, password string not null, name string,
    CONSTRAINT fk_account
        FOREIGN KEY (account_id) REFERENCES account (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE);
CREATE TABLE IF NOT EXISTS project (id string not null primary key, account_id string not null,
    name string not null, description string not null,
    CONSTRAINT fk_account
        FOREIGN KEY (account_id) REFERENCES account (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE);
CREATE TABLE IF NOT EXISTS sprint (id string not null primary key, project_id string not null,
    sprint_number int not null, goal string not null,
    start_date sqlite3_int64 not null, end_date sqlite3_int64 not null,
    CONSTRAINT fk_project
    FOREIGN KEY (project_id) REFERENCES project (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE);
CREATE TABLE IF NOT EXISTS story (id string not null primary key, sprint_id string not null,
    estimation int not null, description string not null,
    status string not null, user_id string,
    CONSTRAINT fk_sprint
        FOREIGN KEY (sprint_id) REFERENCES sprint (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES user (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE);