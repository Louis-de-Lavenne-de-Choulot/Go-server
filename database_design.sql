CREATE TABLE IF NOT EXISTS schedules(
	user_id serial PRIMARY KEY,
	schedule VARCHAR ( 255 ) UNIQUE NOT NULL,
	starting_date TIMESTAMP NOT NULL,
	ending_date TIMESTAMP NOT NULL,
	created_on TIMESTAMP NOT NULL,
	additional_key VARCHAR (32)
)

CREATE TABLE IF NOT EXISTS ###hash###(
	id serial PRIMARY KEY,
	tasks json NOT NULL,
	roles json NOT NULL,
	report json NOT NULL,
	notifs text[] NOT NULL
)

CREATE TABLE IF NOT EXISTS ###hash###_accounts(
	user_id serial PRIMARY KEY,
	username VARCHAR ( 50 ) UNIQUE vvNOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
)