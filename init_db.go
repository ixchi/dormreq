package main

func createDatabase() {
	DB.MustExec(`
		CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name STRING NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS purchase (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name STRING NOT NULL,
			purchased boolean NOT NULL CHECK (purchased in (0, 1)) DEFAULT 0,
			price INTEGER NOT NULL,
			notes STRING,
			category_id int DEFAULT 1
		);

		INSERT INTO category (name)
			SELECT 'Uncategorized'
			WHERE NOT EXISTS (select 1 from category where id = 1);
	`)
}
