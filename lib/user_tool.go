package lib

import "github.com/joycastle/matching-story-robot-service/model"

func FilterUser(targets []model.User, exists []model.User) []model.User {
	m := make(map[int64]struct{})
	for _, u := range exists {
		m[u.UserID] = struct{}{}
	}
	var result []model.User
	for _, u := range targets {
		if _, ok := m[u.UserID]; !ok {
			result = append(result, u)
		}
	}
	return result
}

func UserIds(targets []model.User) []int64 {
	var result []int64
	for _, u := range targets {
		result = append(result, u.UserID)
	}
	return result
}
