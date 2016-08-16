package main

func createDatabase() {
	DB.MustExec(`
		create table if not exists category (
			id integer primary key autoincrement,
			name string not null unique
		);

		create table if not exists purchase (
			id integer primary key autoincrement,
			name string not null,
			purchased boolean not null check(purchased in (0, 1)) default 0,
			price integer not null,
			notes string,
			category_id int default 1
		);

		insert into category (name)
			select 'Uncategorized'
			where not exists (select 1 from category where id = 1);
	`)
}

func testData() {
	DB.MustExec(`
		insert into purchase (name, purchased, price) values
		('A test item', 0, 145),
		('A thing', 1, 354);
	`)
}
