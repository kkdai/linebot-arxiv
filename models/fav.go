package models

type UserFavorite struct {
	Id        int64    `bson:"_id"`
	UserId    string   `json:"user_id" bson:"user_id"`
	Favorites []string `json:"favorites" bson:"favorites"`
}

type UserFavData interface {
	Add(user UserFavorite)
	Get(uid string) (result *UserFavorite, err error)
	ShowAll() (result []UserFavorite, err error)
	Update(user *UserFavorite) (err error)
}
