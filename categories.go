package main

import (
	"github.com/kataras/iris"
)

func listCategories(ctx *iris.Context) {
	var categories []Category
	DB.Select(&categories, `
        SELECT
            id,
            name
        FROM
            category
        `)

	ctx.JSON(200, categories)
}

func addCategory(ctx *iris.Context) {
	name := ctx.PostValue("name")

	DB.Exec(`
        INSERT INTO
            category (name)
            VALUES (?)
    `, name)

	ctx.Redirect("/items")
}

func removeCategory(ctx *iris.Context) {
	id := ctx.PostValue("category")

	if id == "1" {
		ctx.WriteString("Unable to remove default category")
		return
	}

	DB.Exec(`
        DELETE FROM 
            category
        WHERE id = ?
    `, id)

	ctx.Redirect("/items")
}
