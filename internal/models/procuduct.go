package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}
