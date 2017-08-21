package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
)

type User struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int    `json:"birth_date"`
}

func getUserRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if user, ok := usersMap[entityId]; ok {
		response, _ := json.Marshal(user)
		ctx.Success("application/json", response)
		return
	}
	ctx.NotFound()
}

func createUserRequestHandler(ctx *fasthttp.RequestCtx) {
	if user, err := createUser(ctx.PostBody()); err == nil {
		ctx.SetConnectionClose()
		ctx.Success("application/json", []byte("{}"))

		usersMap[user.Id] = user
		return
	}

	ctx.Error("{}", 400)
}

func updateUserRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if user, ok := usersMap[entityId]; ok {
		if updatedUser, err := updateUser(ctx.PostBody(), user); err == nil {
			ctx.SetConnectionClose()
			ctx.Success("application/json", []byte("{}"))

			usersMap[user.Id] = updatedUser
			return
		}
		ctx.Error("{}", 400)
		return
	}
	ctx.NotFound()
}

func createUser(postData []byte) (*User, error) {
	user := User{}
	if err := json.Unmarshal(postData, &user); err != nil {
		return nil, err
	}

	if user.Id == 0 || len(user.First_name) == 0 || len(user.Last_name) == 0 ||
		(user.Gender != "m" && user.Gender != "f") {
		return nil, errors.New("Validation error")
	}
	if _, ok := usersMap[user.Id]; ok {
		return nil, errors.New("User already exists")
	}

	return &user, nil
}

func updateUser(postBody []byte, user *User) (*User, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(postBody, &data); err != nil {
		return nil, err
	}

	updatedUser := *user

	if email, ok := data["email"]; ok {
		if email != nil {
			updatedUser.Email = email.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if firstName, ok := data["first_name"]; ok {
		if firstName != nil {
			updatedUser.First_name = firstName.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if lastName, ok := data["last_name"]; ok {
		if lastName != nil {
			updatedUser.Last_name = lastName.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if gender, ok := data["gender"]; ok {
		if gender != nil && (gender.(string) == "m" || gender.(string) == "f") {
			updatedUser.Gender = gender.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if birthDate, ok := data["birth_date"]; ok {
		if birthDate != nil {
			updatedUser.Birth_date = int(birthDate.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}

	return &updatedUser, nil
}

func updateUsersMaps(user User) {
	usersMap[user.Id] = &user
}
