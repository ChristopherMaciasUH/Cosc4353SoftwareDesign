package clientModel

type UserEntry struct{
	Username string `json: "username"`
	Password string `json: "password"`
}

type UserEntryInfo struct{
	Username string `json: "username"`
	Token string `json: "token"`
}

type ProfileInfo struct{
	Fullname string `json: "fullname"`
	Address [5]string `json: "address"`
}

type RetrievedProfileInfo struct{
	Fullname string `json: "fullname"`
	Address1 string
	Address2 string
	City string
	State string
	Zipcode string
}

type DeliveryData struct {
	Date string `json: "date"`
	Amount string `json: "amount"`
	SuggestedPrice string `json: "suggested"`
	TotalAmount string `json: "total"`
}

type FullDeliveryData struct {
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
