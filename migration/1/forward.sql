create table request (
  uuid               varchar(64) not null,
  parent_uuid        varchar(64),
  service            varchar(32) not null,
  category           varchar(32) not null,
  sync_option        varchar(64),
  begin_ts           bigint,
  end_ts             bigint,
  begin_metadata     mediumblob,
  end_metadata       mediumblob,
  request_group_uuid varchar(64),
  create_ts          bigint,
  update_ts          bigint
);

create table request_group (
  uuid               varchar(64) not null,
  request_seq        mediumblob,
  request_parent_seq mediumblob,
  invoke_chain_id    bigint,
  in_process         tinyint
);

create table request_invoke_chain (
  id                 bigint       not null auto_increment primary key,
  header             varchar(128) not null,
  request_seq        mediumblob   not null,
  request_parent_seq mediumblob   not null
);
