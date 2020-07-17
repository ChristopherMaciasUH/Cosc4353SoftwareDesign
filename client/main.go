package main

import (
	"client/clientModel"
	"client/requests"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)



var tpl *template.Template
var entryInfo clientModel.UserEntryInfo
var currentProfileInfo clientModel.ProfileInfo

func init() {
	tpl = template.Must(template.ParseGlob("templates/*_.gohtml"))
}

var userInfoForRightNow clientModel.UserEntryInfo
//var userProfileInfoForRightNow clientModel.ProfileInfo

func addressToString(addressArr [5]string)(address string){
	address = ""
	for i:= 0; i < len(addressArr); i++{
		if addressArr[i] != "" {
			address += addressArr[i]
		}
		if i < len(addressArr)-1 && addressArr[i] != "" {
			address += ","
		}
	}
	return address
}

func main(){
	router := mux.NewRouter()
	//Pages
	router.HandleFunc("/profile",profile)
	router.HandleFunc("/fuelQuote",fuelQuote)
	router.HandleFunc("/fuelHistory",fuelHistory)
	router.HandleFunc("/",login)
	router.HandleFunc("/register",register)
	//functions
	router.HandleFunc("/login",UserLoginHandler)
	router.HandleFunc("/registration",UserRegistrationHandler)
	router.HandleFunc("/profileInfo",UserProfileManagementHandler)
	router.HandleFunc("/quoteForm",FuelQuoteHandler)
	router.HandleFunc("/logout",logout)
	http.ListenAndServe(":9000",router)
}

func profile(w http.ResponseWriter, r *http.Request){
	currentProfileInfo = requests.UserProfileGetter(entryInfo.Token)
	tpl.ExecuteTemplate(w, "profile.gohtml", currentProfileInfo)
}

func fuelQuote(w http.ResponseWriter, r *http.Request){
	type AddressStruct struct{
		Address string
	}
	AddressRn := AddressStruct{
		Address: addressToString(currentProfileInfo.Address),
	}
	tpl.ExecuteTemplate(w, "fuelQuote.gohtml", AddressRn)
}

func fuelHistory(w http.ResponseWriter, r *http.Request){
	fuelInfo := requests.FuelQuoteInfo(entryInfo.Token)
	tpl.ExecuteTemplate(w, "fuelHistory.gohtml", fuelInfo)
}

func login(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "login_.gohtml", nil)
}

func register(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "register_.gohtml", nil)
}

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request){
	newUserRegistration := clientModel.UserEntry{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	requests.UserRegistration(newUserRegistration)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request){
	newUserLogin := clientModel.UserEntry{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	fmt.Println("we out here")
	returnedUserInfo, err := requests.UserLogin(newUserLogin)
	if err.Error != ""{
		fmt.Println("Login Failed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		fmt.Println(returnedUserInfo)
		entryInfo = returnedUserInfo
		tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func UserProfileManagementHandler(w http.ResponseWriter, r *http.Request){
	newUserProfileInfo := clientModel.RetrievedProfileInfo{
		Fullname: r.FormValue("name"),
		Address1: r.FormValue("address1"),
		Address2: r.FormValue("address2"),
		City: r.FormValue("city"),
		State: r.FormValue("state"),
		Zipcode: r.FormValue("zipcode"),
	}
	fmt.Println(newUserProfileInfo)
	compatibleUserProfileInfo := clientModel.ProfileInfo{
		Fullname: newUserProfileInfo.Fullname,
		Address: [5]string{newUserProfileInfo.Address1,newUserProfileInfo.Address2,newUserProfileInfo.City,newUserProfileInfo.State,newUserProfileInfo.Zipcode},
	}
	fmt.Println(compatibleUserProfileInfo)
	requests.UserProfileSetter(entryInfo.Token, compatibleUserProfileInfo)
	//currentProfileInfo = requests.UserProfileGetter(entryInfo.Token)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func FuelQuoteHandler(w http.ResponseWriter, r *http.Request){
	newDeliveryData := clientModel.DeliveryData{
		Date: r.FormValue("dateInput"),
		Amount: r.FormValue("amount"),
		SuggestedPrice: r.FormValue("suggested"),
		TotalAmount: r.FormValue("total"),
	}
	requests.FuelQuoteForm(entryInfo.Token, newDeliveryData)
	http.Redirect(w, r, "/fuelQuote", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request){
	entryInfo = clientModel.UserEntryInfo{}
	currentProfileInfo = clientModel.ProfileInfo{}
	tpl = template.Must(template.ParseGlob("templates/*_.gohtml"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
