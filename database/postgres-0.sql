DROP TABLE IF EXISTS change_tenancy_log;
DROP TABLE IF EXISTS change_tenancy_code;
DROP TYPE IF EXISTS cot_status;

CREATE TYPE cot_status as enum ('PENDING', 'ACCEPTED', 'VERIFYING', 'ACTIONED', 'REJECTED');

CREATE TABLE IF NOT EXISTS change_tenancy_log  (
request_id SERIAL PRIMARY KEY,
client_id INTEGER NOT NULL,
mac VARCHAR(64) NOT NULL,
created_at TIMESTAMPTZ,
updated_at TIMESTAMPTZ,
status cot_status default 'PENDING',
impl_time TIMESTAMPTZ NOT NULL,
last_fail_time TIMESTAMPTZ,
new_hes_account1 INTEGER,
new_hes_account2 INTEGER,
new_hes_account3 INTEGER,
new_hes_account4 INTEGER,
new_hes_account5 INTEGER,
old_hes_account1 INTEGER,
old_hes_account2 INTEGER,
old_hes_account3 INTEGER,
old_hes_account4 INTEGER,
old_hes_account5 INTEGER
);

CREATE INDEX ctl_client_index ON change_tenancy_log(client_id);
CREATE INDEX ctl_mac_index ON change_tenancy_log(mac);

CREATE UNIQUE INDEX nodups_index ON change_tenancy_log (mac) where ( status != 'ACTIONED' and status != 'REJECTED' ) -- do not allow multiple cots for a single hub, only actioned and rejected. perhaps even delete complete ones.


CREATE TABLE change_tenancy_code (
status cot_status,
short_name VARCHAR(32),
description VARCHAR(255)
);
