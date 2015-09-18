CREATE TABLE users (
  id             BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name           VARCHAR(32), 
  email          VARCHAR(64),
  email_verified TINYINT,
  password       VARCHAR(64),
  create_ts      BIGINT,
  last_login_ts  BIGINT
);

CREATE TABLE repos (
  id        BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id   BIGINT NOT NULL,
  name      VARCHAR(32) NOT NULL,
  hash      VARCHAR(64),
  create_ts BIGINT
);

ALTER TABLE users ADD UNIQUE(name);
ALTER TABLE users ADD UNIQUE(email);

ALTER TABLE repos ADD UNIQUE(hash);
ALTER TABLE repos ADD CONSTRAINT fk_repos_users_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
