CREATE TABLE instances (
  id             BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  host           VARCHAR(32),
  port           INT,
  admin_user     VARCHAR(16),
  admin_password VARCHAR(16),
  active         TINYINT,
  create_ts      BIGINT
);

CREATE TABLE mappings (
  project_id     BIGINT,
  mysql_host     VARCHAR(32),
  mysql_port     INT,
  mysql_user     VARCHAR(16),
  mysql_password VARCHAR(16),
  mysql_db       VARCHAR(16),
  mysql_active   TINYINT,
  create_ts      BIGINT
);

ALTER TABLE instances ADD UNIQUE(host, port);
ALTER TABLE mappings ADD UNIQUE(project_id);
