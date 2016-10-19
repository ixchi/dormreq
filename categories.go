package main

import (
	"encoding/json"
	"net/http"
)

func listCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	DB.Select(&categories, `
        SELECT
            id,
            name
        FROM
            category
        `)

	j, _ := json.Marshal(categories)
	w.Write(j)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]

	DB.Exec(`
        INSERT INTO
            category (name)
            VALUES (?)
    `, name)

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func removeCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["category"][0]

	if id == "1" {
		w.Write([]byte("Unable to remove default category"))
		return
	}

	DB.Exec(`
        DELETE FROM
            category
        WHERE id = ?
    `, id)

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
