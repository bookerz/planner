CREATE SEQUENCE employee_seq start 100;

CREATE TABLE employee (
	id bigint default nextval('employee_seq'),
	first_name varchar(100),
	last_name varchar(100),
	PRIMARY KEY(id)
);