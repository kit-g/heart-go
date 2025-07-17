package models

type ImageDescription struct {
	Link   *string `json:"link" example:"https://example.com/image.jpg"`
	Width  *int    `json:"width" example:"100"`
	Height *int    `json:"height" example:"100"`
} // @name ImageDescription
