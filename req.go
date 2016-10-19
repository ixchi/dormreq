package main

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

//go:generate go-bindata templates/

// DB is the database connection.
var DB *sqlx.DB

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db := "info.db"
	if envDb := os.Getenv("DB"); envDb != "" {
		db = envDb
	}

	DB = sqlx.MustOpen("sqlite3", db)
	createDatabase()

	r.Get("/", homepage)

	r.Route("/items", func(r chi.Router) {
		r.Get("/", listItems)

		r.Post("/add", addItem)

		r.Route("/:itemID", func(r chi.Router) {
			r.Use(ItemCtx)

			r.Post("/remove", removeItem)
			r.Post("/purchase", purchaseItem)
			r.Post("/unpurchase", unpurchaseItem)
		})
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", listCategories)

		r.Post("/add", addCategory)
		r.Post("/remove", removeCategory)
	})

	host := ":8080"
	if envHost := os.Getenv("HOST"); envHost != "" {
		host = envHost
	}

	http.ListenAndServe(host, r)
}

func ItemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemID := chi.URLParam(r, "itemID")

		var item Purchase
		err := DB.Get(&item, `SELECT * FROM purchase WHERE id = ?`, itemID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), "item", item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func homepage(w http.ResponseWriter, r *http.Request) {
	html, _ := templatesHomeHtmlBytes()
	t, _ := template.New("homepage").Parse(string(html))
	t.Execute(w, nil)
}
