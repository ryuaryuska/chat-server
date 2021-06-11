package model

type CodeMsg struct {
	Code    string `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}
