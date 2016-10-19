package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func listItems(w http.ResponseWriter, r *http.Request) {
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

	html, _ := templatesListHtmlBytes()
	t, _ := template.New("items").Parse(string(html))
	t.Execute(w, struct {
		Purchases           map[string][]PurchaseRender
		Categories          []Category
		RemovableCategories []Category
		TotalCost           string
	}{sorted, categories, removableCategories, ParsePrice(cost)})
}

func addItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]
	price := r.Form["price"][0]
	parsedPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("invalid price"))
		return
	}
	notes := r.Form["notes"][0]
	categoryID := r.Form["category"][0]
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
		w.Write([]byte("unknown category"))
		return
	}

	_, err = DB.Exec(`
        INSERT INTO
            purchase (name, price, notes, category_id)
		    VALUES (?, ?, ?, ?)
    `, name, parsedPrice*100, notes, categoryID)
	if err != nil {
		w.Write([]byte("error storing data"))
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func removeItem(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value("item").(Purchase)

	DB.Exec(`
        DELETE FROM
            purchase
        WHERE
            id = ?
    `, item.ID)

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func purchaseItem(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value("item").(Purchase)

	DB.Exec(`
        UPDATE
            purchase
        SET
            purchased = 1
        WHERE
            id = ?
    `, item.ID)

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func unpurchaseItem(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value("item").(Purchase)

	DB.Exec(`
        UPDATE
            purchase
        SET
            purchased = 0
        WHERE
            id = ?
    `, item.ID)

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
