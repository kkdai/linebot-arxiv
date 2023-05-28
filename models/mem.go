package models

import (
	"errors"
	"log"
)

type MemStorage map[string]UserFavorite

type MemDB struct {
	Db MemStorage
}

func NewMemDB() *MemDB {
	return &MemDB{
		Db: make(MemStorage),
	}
}

func (u *MemDB) Add(user UserFavorite) {
	u.Db[user.UserId] = user
}

func (u *MemDB) Get(uid string) (result *UserFavorite, err error) {
	ret, exist := u.Db[uid]
	log.Println("***Get Fav uUID=", uid)
	if !exist {
		log.Println("No result on uid:", uid)
		return nil, errors.New("No result on uid:" + uid)
	}
	log.Println("UserFavorite DB result= ", ret)
	return &ret, nil
}

// ShowAll: Print all result.
func (u *MemDB) ShowAll() (result []UserFavorite, err error) {
	log.Println("***Get All DB")
	for _, v := range u.Db {
		result = append(result, v)
	}
	log.Println("***Start server all users =", u.Db)
	return result, nil
}

func (u *MemDB) Update(user *UserFavorite) (err error) {
	log.Println("***Update Fav User=", u)
	if _, exist := u.Db[user.UserId]; !exist {
		return errors.New("No result on uid:" + user.UserId)
	} else {
		u.Db[user.UserId] = *user
	}
	return nil
}
