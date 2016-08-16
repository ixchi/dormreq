package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/kataras/iris"
)

func listItems(ctx *iris.Context) {
	var purchases []Purchase
	err := DB.Select(&purchases, `
        SELECT
            *
        FROM
            purchase
        ORDER BY
            purchased asc,
            id desc
    `)
	if err != nil {
		fmt.Println(err.Error())
	}

	var cost float32

	var rendered []PurchaseRender
	for _, purchase := range purchases {
		if !purchase.Purchased {
			cost += purchase.Price
		}
		rendered = append(rendered, NewPurchaseRender(purchase))
	}

	sorted := make(map[string][]PurchaseRender)

	for _, purchase := range rendered {
		if _, ok := sorted[purchase.Category]; !ok {
			sorted[purchase.Category] = []PurchaseRender{}
		}

		sorted[purchase.Category] = append(sorted[purchase.Category], purchase)
	}

	var categories []Category
	DB.Select(&categories, `
        SELECT
            id,
            name
        FROM
            category
    `)

	var removableCategories []Category
	DB.Select(&removableCategories, `
        SELECT
            id,
            name
        FROM
            category
        WHERE
            id NOT IN (
                SELECT
                    category_id 
                FROM
                    purchase
                GROUP BY
                    category_id
            )
    `)

	ctx.Render("list.html", struct {
		Purchases           map[string][]PurchaseRender
		Categories          []Category
		RemovableCategories []Category
		TotalCost           string
	}{sorted, categories, removableCategories, ParsePrice(cost)})
}

func addItem(ctx *iris.Context) {
	name := ctx.PostValue("name")
	price := ctx.PostValue("price")
	parsedPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Println(err.Error())
		ctx.WriteString("invalid price")
		return
	}
	notes := ctx.PostValue("notes")
	categoryID := ctx.PostValue("category")
	var category Category
	err = DB.Get(&category, `
        SELECT
            id,
            name
        FROM
            category
        WHERE
            id = ?
    `, categoryID)
	if err == sql.ErrNoRows {
		ctx.WriteString("unknown category")
		return
	}

	_, err = DB.Exec(`
        INSERT INTO
            purchase (name, price, notes, category_id)
		    VALUES (?, ?, ?, ?)
    `, name, parsedPrice*100, notes, categoryID)
	if err != nil {
		ctx.WriteString("error storing data")
		fmt.Println(err.Error())
		return
	}

	ctx.Redirect("/items")
}

func removeItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`
        DELETE FROM
            purchase
        WHERE
            id = ?
    `, id)

	ctx.Redirect("/items")
}

func purchaseItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`
        UPDATE
            purchase
        SET
            purchased = 1
        WHERE
            id = ?
    `, id)

	ctx.Redirect("/items")
}

func unpurchaseItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`
        UPDATE
            purchase
        SET
            purchased = 0
        WHERE
            id = ?
    `, id)

	ctx.Redirect("/items")
}
