package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"sync"
)

type User struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int    `json:"birth_date"`
}

type UsersMap struct {
	users map[uint]*User
	sync.RWMutex
}

func (u *UsersMap) Get(id uint) *User {
	u.RLock()
	defer u.RUnlock()

	return u.users[id]
}

func (u *UsersMap) Update(user User) {
	u.Lock()
	u.users[user.Id] = &user
	u.Unlock()
}

func getUserRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if user := usersMap.Get(entityId); user != nil {
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

		go func() {
			usersMap.Update(*user)
		}()
		return
	}

	ctx.Error("{}", 400)
}

func updateUserRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if user := usersMap.Get(entityId); user != nil {
		if updatedUser, err := updateUser(ctx.PostBody(), user); err == nil {
			ctx.SetConnectionClose()
			ctx.Success("application/json", []byte("{}"))

			go func() {
				usersMap.Update(*updatedUser)
			}()
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
	if user := usersMap.Get(user.Id); user != nil {
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
