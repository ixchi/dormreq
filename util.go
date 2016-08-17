package main

import (
	"fmt"
)

// ParsePrice turns an integer version of the price into a human formatted version.
func ParsePrice(price float32) string {
	num := float32(price) / 100
	str := "$" + fmt.Sprintf("%.2f", num)

	return str
}

// NewPurchaseRender converts a Purchase from the database into a PurchaseRender.
func NewPurchaseRender(purchase Purchase) PurchaseRender {
	var notes string
	if purchase.Notes == nil {
		notes = ""
	} else {
		notes = *purchase.Notes
	}

	var category Category
	DB.Get(&category, `
        SELECT
            id,
            name
        FROM
            category
        WHERE
            id = ?
    `, purchase.CategoryID)

	return PurchaseRender{
		ID:        purchase.ID,
		Name:      purchase.Name,
		Price:     ParsePrice(purchase.Price),
		Notes:     notes,
		Category:  category.Name,
		Purchased: purchase.Purchased,
	}
}
