package main

// Purchase is a purchase item stored in the database.
type Purchase struct {
	ID         int64   `db:"id"`
	Name       string  `db:"name"`
	Purchased  bool    `db:"purchased"`
	Price      float32 `db:"price"`
	Notes      *string `db:"notes"`
	CategoryID int64   `db:"category_id"`
}

// Category is a categorization of a Purchase.
type Category struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// PurchaseRender is formatted data for displaying.
type PurchaseRender struct {
	ID        int64
	Name      string
	Price     string
	Notes     string
	Category  string
	Purchased bool
}
