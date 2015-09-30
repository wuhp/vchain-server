CREATE TABLE request (
  uuid        VARCHAR(64) NOT NULL,
  parent_uuid VARCHAR(64),
  service     VARCHAR(32) NOT NULL,
  category    VARCHAR(32) NOT NULL,
  sync        TINYINT,
  begin_ts    BIGINT,
  end_ts      BIGINT,
  group_uuid  VARCHAR(64),
  create_ts   BIGINT
);

CREATE TABLE request_group (
  uuid            VARCHAR(64) NOT NULL,
  request_uuids   MEDIUMBLOB,
  parents_index   MEDIUMBLOB,
  invoke_chain_id BIGINT
);

CREATE TABLE invoke_chain (
  id            BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
  header        VARCHAR(128) NOT NULL,
  request_types MEDIUMBLOB   NOT NULL,
  parents_index MEDIUMBLOB   NOT NULL
);

CREATE TABLE request_log (
  uuid      VARCHAR(64) NOT NULL,
  timestamp BIGINT      NOT NULL,
  log       MEDIUMBLOB  NOT NULL
);
