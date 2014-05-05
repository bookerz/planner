CREATE SEQUENCE customer_seq start 100;

CREATE TABLE customer (
	id bigint default nextval('customer_seq'),
	first_name varchar(100),
	last_name varchar(100),
	PRIMARY KEY(id)
);