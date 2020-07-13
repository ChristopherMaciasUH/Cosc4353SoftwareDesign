package requests

import (
	"bytes"
	"client/clientModel"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UserRegistration(info clientModel.UserEntry) {
	//info := clientModel.UserEntry{
	//	Username: "clientTest",
	//	Password: "password",
	//}
	fmt.Println("Does this run?")
	jsonValue, _ := json.Marshal(info)
	response, err := http.Post("http://localhost:8000/register", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func UserLogin(info clientModel.UserEntry) (clientModel.UserEntryInfo, clientModel.ResponseResult) {
	//info := clientModel.UserEntry{
	//	Username: "clientTest",
	//	Password: "password",
	//}
	fmt.Println("Does this run too?")
	jsonValue, _ := json.Marshal(info)
	response, err := http.Post("http://localhost:8000/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		var error clientModel.ResponseResult
		fmt.Println(err)
		fmt.Println(response)
		fmt.Println("login failed my dude")
		json.NewDecoder(response.Body).Decode(&error)
		return clientModel.UserEntryInfo{}, error
	} else {
		var userEntry clientModel.UserEntryInfo
		json.NewDecoder(response.Body).Decode(&userEntry)
		//data, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(string(data))
		fmt.Println(userEntry)
		return userEntry, clientModel.ResponseResult{}
	}
}

func UserProfileSetter(token string, userInfo clientModel.ProfileInfo) {
	jsonValue, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println(err)
	}
	//var bearer = "Bearer " + token
	request, _ := http.NewRequest("PUT", "http://localhost:8000/profileSetter", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func UserProfileGetter(token string) clientModel.ProfileInfo {
	request, _ := http.NewRequest("GET", "http://localhost:8000/profileInfo", nil)
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return clientModel.ProfileInfo{}
	} else {
		var returnInfo clientModel.ProfileInfo
		json.NewDecoder(response.Body).Decode(&returnInfo)
		fmt.Println(returnInfo)
		return returnInfo
	}
}

func FuelQuoteForm(token string, deliveryInfo clientModel.DeliveryData) {
	jsonValue, err := json.Marshal(deliveryInfo)
	if err != nil {
		fmt.Println(err)
	}
	request, _ := http.NewRequest("PUT", "http://localhost:8000/fuelQuoteForm", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
func FuelQuoteInfo(token string) clientModel.FullDeliveryData {
	request, _ := http.NewRequest("GET", "http://localhost:8000/fuelQuoteHistory", nil)
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return clientModel.FullDeliveryData{}
	} else {
		var returnInfo clientModel.FullDeliveryData
		data, _ := ioutil.ReadAll(response.Body)
		json.NewDecoder(response.Body).Decode(&returnInfo)
		fmt.Println(string(data))
		return returnInfo
	}
}
