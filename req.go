package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var DB *sqlx.DB

type Purchase struct {
	ID         int64   `db:"id"`
	Name       string  `db:"name"`
	Purchased  bool    `db:"purchased"`
	Price      float32 `db:"price"`
	Notes      *string `db:"notes"`
	CategoryID int64   `db:"category_id"`
}

type Category struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func main() {
	app := iris.New()

	DB = sqlx.MustOpen("sqlite3", "info.db")

	createDatabase()
	//testData()

	app.Get("/", homepage)

	app.Get("/items", listItems)
	app.Post("/items/add", addItem)
	app.Post("/items/remove", removeItem)

	app.Post("/items/purchase", purchaseItem)
	app.Post("/items/unpurchase", unpurchaseItem)

	app.Get("/categories", listCategories)
	app.Post("/categories/add", addCategory)
	app.Post("/categories/remove", removeCategory)

	app.Listen(":8080")
}

func homepage(ctx *iris.Context) {
	ctx.Render("home.html", nil)
}

type PurchaseRender struct {
	ID        int64
	Name      string
	Price     string
	Notes     string
	Category  string
	Purchased bool
}

func ParsePrice(price float32) string {
	num := float32(price) / 100
	str := "$" + fmt.Sprintf("%.2f", num)

	return str
}

func NewPurchaseRender(purchase Purchase) PurchaseRender {
	var notes string
	if purchase.Notes == nil {
		notes = ""
	} else {
		notes = *purchase.Notes
	}

	var category Category
	DB.Get(&category, `select * from category where id = ?`, purchase.CategoryID)

	return PurchaseRender{
		ID:        purchase.ID,
		Name:      purchase.Name,
		Price:     ParsePrice(purchase.Price),
		Notes:     notes,
		Category:  category.Name,
		Purchased: purchase.Purchased,
	}
}

func listItems(ctx *iris.Context) {
	var purchases []Purchase
	err := DB.Select(&purchases, "select * from purchase order by purchased asc, id desc")
	if err != nil {
		fmt.Println(err.Error())
	}

	var cost float32 = 0

	var rendered []PurchaseRender
	for _, purchase := range purchases {
		cost += purchase.Price
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
	DB.Select(&categories, `select * from category`)

	var removable_categories []Category
	DB.Select(&removable_categories, `select * from category where id not in (select category_id from purchase group by category_id)`)

	ctx.Render("list.html", struct {
		Purchases           map[string][]PurchaseRender
		Categories          []Category
		RemovableCategories []Category
		TotalCost           string
	}{sorted, categories, removable_categories, ParsePrice(cost)})
}

func addItem(ctx *iris.Context) {
	name := ctx.PostValue("name")
	price := ctx.PostValue("price")
	parsed_price, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Println(err.Error())
		ctx.WriteString("invalid price")
		return
	}
	notes := ctx.PostValue("notes")
	category_id := ctx.PostValue("category")
	var category Category
	err = DB.Get(&category, "select * from category where id = ?", category_id)
	if err == sql.ErrNoRows {
		ctx.WriteString("unknown category")
		return
	}

	_, err = DB.Exec(`insert into purchase
		(name, price, notes, category_id)
		values (?, ?, ?, ?)`, name, parsed_price*100, notes, category_id)
	if err != nil {
		ctx.WriteString("error storing data")
		fmt.Println(err.Error())
		return
	}

	ctx.Redirect("/items")
}

func removeItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`delete from purchase where id = ?`, id)

	ctx.Redirect("/items")
}

func purchaseItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`update purchase set purchased = 1 where id = ?`, id)

	ctx.Redirect("/items")
}

func unpurchaseItem(ctx *iris.Context) {
	id := ctx.PostValue("item")

	DB.Exec(`update purchase set purchased = 0 where id = ?`, id)

	ctx.Redirect("/items")
}

func purgePurchasedItems(ctx *iris.Context) {
	DB.Exec(`delete from purchase where purchased = 1`)

	ctx.Redirect("/items")
}

func listCategories(ctx *iris.Context) {
	var categories []Category
	DB.Select(&categories, `select * from category`)

	ctx.JSON(200, categories)
}

func addCategory(ctx *iris.Context) {
	name := ctx.PostValue("name")

	DB.Exec(`insert into category (name) values (?)`, name)

	ctx.Redirect("/items")
}

func removeCategory(ctx *iris.Context) {
	id := ctx.PostValue("category")

	if id == "1" {
		ctx.WriteString("Unable to remove default category")
		return
	}

	DB.Exec(`delete from category where id = ?`, id)

	ctx.Redirect("/items")
}
