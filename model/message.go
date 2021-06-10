package model

type MessageCode struct {
	Text string `json:"text"`
}

type CodeMsg struct {
	Code    string `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}
