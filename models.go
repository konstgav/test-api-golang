package main

type Product struct {
	ID    int    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Sku   string `json:"sku" bson:"sku,omitempty"`
	Type  string `json:"type" bson:"type,omitempty"`
	Price int    `json:"price" bson:"price,omitempty"`
}
