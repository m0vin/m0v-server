DROP TABLE IF EXISTS subpub;
DROP TABLE IF EXISTS sub;
DROP TABLE IF EXISTS csub;
DROP TABLE IF EXISTS packet;
DROP TABLE IF EXISTS confo;
DROP TABLE IF EXISTS config;
DROP TABLE IF EXISTS pub;

CREATE TABLE IF NOT EXISTS pub (
	pub_id SERIAL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL default current_timestamp,
	latitude REAL NOT NULL,
	longitude REAL NOT NULL,
	altitude REAL NOT NULL,
	orientation REAL NOT NULL,
	hash INTEGER NOT NULL UNIQUE,
	creator INTEGER,
	protected BOOLEAN default TRUE
);

CREATE TABLE IF NOT EXISTS pubconfig (
	pub_hash INTEGER NOT NULL REFERENCES pub(hash),
	id SERIAL PRIMARY KEY,
	nickname VARCHAR(16),
	typeref VARCHAR(16),
	kwp REAL,
	kwpmake VARCHAR(16),
	kwr REAL,
	kwrmake VARCHAR(16),
	kw_last REAL,
	kwh_hour REAL,
	kwh_day REAL,
	kwh_life REAL,
	since TIMESTAMPTZ default current_timestamp,
	visits_last TIMESTAMPTZ default current_timestamp,
	visits_life INTEGER,
	notify BOOLEAN default FALSE,
	lastnotified TIMESTAMPTZ default current_timestamp
);

CREATE TABLE IF NOT EXISTS sub (
	sub_id SERIAL PRIMARY KEY,
	email VARCHAR(64),
	phone VARCHAR(32),
	name VARCHAR(32),
	pswd VARCHAR(16),
	created_at TIMESTAMPTZ default current_timestamp,
	verification VARCHAR(16),
	verified BOOLEAN default FALSE
);

CREATE TABLE IF NOT EXISTS csub (
	sub_id SERIAL PRIMARY KEY,
	email VARCHAR(64),
	created_at TIMESTAMPTZ default current_timestamp
);

CREATE TABLE IF NOT EXISTS subpub (
	sub_id INTEGER NOT NULL REFERENCES sub(sub_id),
	pub_id INTEGER NOT NULL REFERENCES pub(pub_id)
);

CREATE TABLE IF NOT EXISTS packet (
	pub_hash INTEGER NOT NULL REFERENCES pub(hash),
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL default current_timestamp,
	saved_at TIMESTAMPTZ NOT NULL default current_timestamp,
	voltage REAL,
	frequency REAL,
	protected boolean NOT NULL,
	active_power REAL,
	apparent_power REAL,
	reactive_power REAL,
	power_factor REAL,
	import_active_energy REAL,
	export_active_energy REAL,
	import_reactive_energy REAL,
	export_reactive_energy REAL,
	total_active_energy REAL,
	total_reactive_energy REAL
);

CREATE TABLE IF NOT EXISTS confo (
	id SERIAL,
	created_at TIMESTAMPTZ NOT NULL default current_timestamp,
	devicename VARCHAR(32),
	ssid VARCHAR(32),
	hash INTEGER NOT NULL,
	PRIMARY KEY (created_at, devicename, ssid)
);

/*CREATE INDEX sub_index ON coordinate(user_id);*/
CREATE INDEX packet_index ON packet(id);

CREATE TABLE IF NOT EXISTS hourly (
	pub_hash INTEGER NOT NULL REFERENCES pub(hash),
	timestamp TIMESTAMPTZ,
	voltage_max REAL,
	voltage_min REAL,
	voltage_ave REAL,
	voltage_exceptions INTEGER,
	frequency_max REAL,
	frequency_min REAL,
	frequency_ave REAL,
	frequency_exceptions INTEGER,
	activepwr_max REAL,
	activepwr_min REAL,
	activepwr_ave REAL,
	import_active_energy REAL,
	export_active_energy REAL,
	import_reactive_energy REAL,
	export_reactive_energy REAL,
	total_active_energy REAL,
	total_reactive_energy REAL
);

CREATE TABLE IF NOT EXISTS daily (
	pub_hash INTEGER NOT NULL REFERENCES pub(hash),
	timestamp TIMESTAMPTZ,
	voltage_max REAL,
	voltage_min REAL,
	voltage_ave REAL,
	voltage_exceptions INTEGER,
	frequency_max REAL,
	frequency_min REAL,
	frequency_ave REAL,
	frequency_exceptions INTEGER,
	import_active_energy REAL,
	export_active_energy REAL,
	import_reactive_energy REAL,
	export_reactive_energy REAL,
	total_active_energy REAL,
	total_reactive_energy REAL
);

CREATE TABLE IF NOT EXISTS monthly (
	pub_hash INTEGER NOT NULL REFERENCES pub(hash),
	timestamp TIMESTAMPTZ,
	voltage_max REAL,
	voltage_min REAL,
	voltage_ave REAL,
	voltage_exceptions INTEGER,
	frequency_max REAL,
	frequency_min REAL,
	frequency_ave REAL,
	frequency_exceptions INTEGER,
	import_active_energy REAL,
	export_active_energy REAL,
	import_reactive_energy REAL,
	export_reactive_energy REAL,
	total_active_energy REAL,
	total_reactive_energy REAL
);
