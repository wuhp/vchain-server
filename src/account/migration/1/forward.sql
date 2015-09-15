CREATE TABLE user (
  id            BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name          VARCHAR(32), 
  email         VARCHAR(64),
  password      VARCHAR(64),
  create_ts     BIGINT,
  last_login_ts BIGINT
);

CREATE TABLE repo (
  id        BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id   BIGINT NOT NULL,
  name      VARCHAR(32) NOT NULL,
  hash      VARCHAR(64),
  create_ts BIGINT
);

ALTER TABLE user ADD UNIQUE(name);
ALTER TABLE user ADD UNIQUE(email);

ALTER TABLE repo ADD UNIQUE(hash);
ALTER TABLE repo ADD CONSTRAINT fk_repo_user_id FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
