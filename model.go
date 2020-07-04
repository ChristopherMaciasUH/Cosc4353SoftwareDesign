package model
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	Username     string             `json: "username"`
	Password     string             `json: "password"`
	Token        string             `json: "token"`
	PersonalInfo primitive.ObjectID `json: info`
	Deliveries   primitive.ObjectID `json: "deliveries"`
}

type UserInfo struct{
	FullName string `json: "fullname"`
	Address []string `json: "address"`
}

type DeliveryData struct {
	FullName []string `json: "name"`
	Address []string `json :"address"`
	Date []string `json: "date"`
	Amount []string `json: "amount"`
	SuggestedPrice []string `json: "suggested"`
	TotalAmount []string `json: "total"`
}
type ResponseResult struct {
	Error string `json: "error"`
	Result string `json: "result"`
}

