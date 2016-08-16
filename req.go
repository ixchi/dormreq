package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection.
var DB *sqlx.DB

func main() {
	app := iris.New()

	db := "info.db"
	if envDb := os.Getenv("DB"); envDb != "" {
		db = envDb
	}

	DB = sqlx.MustOpen("sqlite3", db)
	createDatabase()

	app.Get("/", homepage)

	app.Get("/items", listItems)
	app.Post("/items/add", addItem)
	app.Post("/items/remove", removeItem)

	app.Post("/items/purchase", purchaseItem)
	app.Post("/items/unpurchase", unpurchaseItem)

	app.Get("/categories", listCategories)
	app.Post("/categories/add", addCategory)
	app.Post("/categories/remove", removeCategory)

	host := ":8080"
	if envHost := os.Getenv("HOST"); envHost != "" {
		host = envHost
	}

	app.Listen(host)
}

func homepage(ctx *iris.Context) {
	ctx.Render("home.html", nil)
}
