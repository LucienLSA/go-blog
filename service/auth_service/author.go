package auth_service

import "blog-service/models"

type Author struct {
	Username string
	Password string
}

func (a *Author) Check() (bool, error) {
	return models.CheckAuthor(a.Username, a.Password)
}
