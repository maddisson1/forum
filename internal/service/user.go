package service

import (
	"errors"
	"forum/models"
	"forum/pkg/validator"
	"net/http"
)

func (s *service) UserPosts(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	username, err := s.repo.GetUser(data.AuthenticatedUser)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.UserPosts(data.AuthenticatedUser)
	if err != nil {
		return nil, err
	}
	data.Posts = posts
	data.Username = username

	return data, nil
}

func (s *service) UserSignUp(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	form := &models.UserSignupForm{
		Name:     r.FormValue("username-signup"),
		Email:    r.FormValue("email-signup"),
		Password: r.FormValue("password-signup"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data.Form = form
		return data, err
	}

	// username := r.PostFormValue("username-signup")
	// email := r.PostFormValue("email-signup")
	// password := r.PostFormValue("password-signup")
	passwordAgain := r.PostFormValue("password-again")
	// fmt.Println(username, email, password, passwordAgain)

	if form.Password != passwordAgain {
		form.AddFieldError("password", "Password is incorrect")
		data.Form = form
		return data, err
	}
	err = s.repo.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data.Form = form
			return data, err
		} else {
			return nil, err
		}
	}

	return data, nil
}

func (s *service) UserLogin(data *models.TemplateData, r *http.Request) (*models.TemplateData, int, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, 0, err
	}

	form := models.UserLoginForm{
		Email:    r.FormValue("email-login"),
		Password: r.FormValue("password-login"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data.Form = form
		return nil, 0, err
	}

	id, err := s.repo.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("email", "Email or password is incorrect")

			data.Form = form

			// h.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return nil, 0, err
		} else {
			// h.serverError(w, err)
			return nil, 0, err
		}
	}
	return data, id, nil
}
