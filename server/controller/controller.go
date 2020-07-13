package controller

import (
	"context"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"quickstart/config/db"
	"quickstart/model"
)



func HashPassword(password string)(string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes), err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("running...")
	w.Header().Set("Content-Type","application/json")
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var res model.ResponseResult
	var deliveryModel model.DeliveryData
	collection, deliveryCollection,infoCollection, err := db.GetDBcollection()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username",user.Username}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := HashPassword(user.Password)
			if err != nil {
				res.Error = "Error while hashing password"
				json.NewEncoder(w).Encode(res)
				fmt.Println(res.Error)
				return
			}
			user.Password = hash
			deliveryModel.FullName = append(deliveryModel.FullName, "")
			deliveryModel.Address = append(deliveryModel.Address, "")
			deliveryModel.Date = append(deliveryModel.Date, "")
			deliveryModel.Amount = append(deliveryModel.Amount, "")
			deliveryModel.SuggestedPrice = append(deliveryModel.SuggestedPrice, "")
			deliveryModel.TotalAmount = append(deliveryModel.TotalAmount, "")
			deliveryId, err := deliveryCollection.InsertOne(context.TODO(), deliveryModel)
			var infoModel model.UserInfo
			infoId, err := infoCollection.InsertOne(context.TODO(), infoModel)
			user.PersonalInfo = infoId.InsertedID.(primitive.ObjectID)
			user.Deliveries = deliveryId.InsertedID.(primitive.ObjectID)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error while creating new User, try again"
				json.NewEncoder(w).Encode(res)
				fmt.Println(res.Error)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			fmt.Println(res.Result)
			return
		}
		res.Result = "User already exists"
		json.NewEncoder(w).Encode(res)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r* http.Request){
	w.Header().Set("Content-Type","application/json")
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	collection, _,_,err := db.GetDBcollection()
	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult
	err = collection.FindOne(context.TODO(), bson.D{{"username",user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		res.Error = "invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
		res.Error = "Error while generating token"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""
	json.NewEncoder(w).Encode(result)
}

func GetProfileInfo(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	tokenString  := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username =  claims["username"].(string)
		userCollection,_,infoCollection,err := db.GetDBcollection()
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		var user model.User
		var information model.UserInfo
		err = userCollection.FindOne(context.TODO(),bson.D{{"username",result.Username}}).Decode(&user)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
		}
		filter := bson.M{"_id":bson.M{"$eq":user.PersonalInfo}}
		err = infoCollection.FindOne(context.TODO(),filter).Decode(&information)
		json.NewEncoder(w).Encode(information)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func InsertProfileInfo(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	tokenString  := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var user model.User
	var res model.ResponseResult
	var information model.UserInfo
	_ = json.NewDecoder(r.Body).Decode(&information)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		collection, _,infoCollection, err := db.GetDBcollection()
		if err != nil {
			log.Fatal(err)
		}
		user.Username =  claims["username"].(string)
		fmt.Println(user)
		err = collection.FindOne(context.TODO(), bson.D{{"username",user.Username}}).Decode(&user)
		update := bson.D{{"$set",bson.M{"fullname":information.FullName, "address":information.Address}}}
		filter := bson.M{"_id": bson.M{"$eq":user.PersonalInfo}}
		result, err := infoCollection.UpdateOne(context.TODO(),filter, update)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		fmt.Println(result)
		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func DeliveryRequestHandler(w http.ResponseWriter, r* http.Request){
	w.Header().Set("Content-Type", "application/json")
	tokenString  := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var user model.User
	var res model.ResponseResult
	var delivery model.DeliveryData
	var userInfo model.UserInfo
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		collection, deliveryCollection,information, err := db.GetDBcollection()
		if err != nil {
			log.Fatal(err)
		}
		user.Username =  claims["username"].(string)
		err = collection.FindOne(context.TODO(), bson.D{{"username",user.Username}}).Decode(&user)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		err = deliveryCollection.FindOne(context.TODO(), bson.D{{"_id",user.Deliveries}}).Decode(&delivery)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		err = information.FindOne(context.TODO(), bson.D{{"_id",user.PersonalInfo}}).Decode(&userInfo)
		//TODO: create an imported JSON delivery model
		delivery.FullName = append(delivery.FullName,userInfo.FullName)
		delivery.Address = append(delivery.Address,userInfo.Address[0])
		delivery.Date = append(delivery.Date,user.Username)
		delivery.Amount = append(delivery.Amount,user.Username)
		delivery.SuggestedPrice = append(delivery.SuggestedPrice,user.Username)
		delivery.TotalAmount = append(delivery.TotalAmount,user.Username)
		update := bson.M{"$set":bson.M{
			"fullname": delivery.FullName,
			"address": delivery.Address,
			"date": delivery.Date,
			"amount": delivery.Amount,
			"suggestedprice": delivery.SuggestedPrice,
			"totalamount": delivery.TotalAmount,
		}}
		filter := bson.M{"_id":bson.M{"$eq": user.Deliveries}}
		result, err := deliveryCollection.UpdateOne(context.Background(),filter, update)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		fmt.Println(result)
		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func GetDeliveryRequests(w http.ResponseWriter, r* http.Request){
	w.Header().Set("Content-Type", "application/json")
	tokenString  := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var user model.User
	var deliveryData model.DeliveryData
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		collection, deliveryCollection,_,err := db.GetDBcollection()
		if err != nil {
			log.Fatal(err)
		}
		user.Username =  claims["username"].(string)
		err = collection.FindOne(context.TODO(), bson.D{{"username",user.Username}}).Decode(&user)
		filter := bson.M{"_id":bson.M{"$eq":user.Deliveries}}
		err = deliveryCollection.FindOne(context.TODO(),filter).Decode(&deliveryData)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		fmt.Println(deliveryData)
		json.NewEncoder(w).Encode(deliveryData)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}
