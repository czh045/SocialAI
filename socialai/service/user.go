package service

import (
	"fmt"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	"github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("username", username)).
		Must(elastic.NewTermQuery("password", password))

	result, err := backend.ESBackend.ReadFromES(query, constants.UserIndex)
	if err != nil {
		return false, err
	}
	return result.TotalHits() > 0, nil
}

func AddUser(user *model.User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
	result, err := backend.ESBackend.ReadFromES(query, constants.UserIndex)
	if err != nil {
		return false, err
	}
	if result.TotalHits() > 0 {
		return false, nil
	}

	if err := backend.ESBackend.SaveToES(user, constants.UserIndex, user.Username); err != nil {
		return false, err
	}

	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
