create extension if not exists "pgcrypto";

create table coin (
	id serial primary key,
	uuid uuid unique not null,
	symbol varchar(50) unique not null,
	watching boolean not null,
	interval int not null,
	check (interval > 0)
);

create table coin_price_log (
	id serial primary key,
	uuid uuid unique not null,
	price_usd text not null,
	coin_uuid uuid references coin(uuid) not null,
	collected_at timestamp not null
);
