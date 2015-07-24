create table service (
  uuid     varchar(64) not null,
  app_id   int         not null,
  category varchar(32),
  instance varchar(32),
  hostname varchar(64),
  start_ts bigint,
  stop_ts  bigint
);

create table request (
  uuid           varchar(64) not null,
  service_uuid   varchar(64) not null,
  parent_uuid    varchar(64),
  category       varchar(64),
  begin_ts       bigint,
  end_ts         bigint,
  begin_metadata mediumblob,
  end_metadata   mediumblob
);
