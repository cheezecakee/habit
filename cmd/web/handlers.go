package main

import (
	"errors"
	"log"
	"net/http"

	"habits.cheezecake.net/internal/models"
	"habits.cheezecake.net/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
  data := app.newTemplateData(r)

  if ok {
  
    habits, err := app.habits.List(userID)
    if err != nil {
      app.serveError(w, err)
    return
    }

    data.Habits = habits

    data.IsAuthenticated = true
  } else {
    data.IsAuthenticated = false
  }

  app.render(w, http.StatusOK, "home.html", data)
} 

type habitCreateForm struct {
  Title               string `form:"title"`
  validator.Validator `form:"-"`
}

func (app *application) habitCreate(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  data.Form = habitCreateForm{}

  app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) habitCreatePost(w http.ResponseWriter, r *http.Request) {
  var form habitCreateForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
  form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")

  // If there are any errors, dump them in a plain text HTTP response and
  // return from the handler.
  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "create.html", data)
    return
  } 

  userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
  if !ok {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }
  
  id, err := app.habits.Insert(form.Title, userID)
  if err != nil {
    app.serveError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Habit successfully created!")

  log.Printf("New habit created with ID: %d", id)
  log.Printf("User id: %d", userID)
  http.Redirect(w, r, "/", http.StatusSeeOther)
}

type userSignupForm struct {
  Name                string `form:"name"`
  Email               string `form:"email"`
  Password            string `form:"password"`
  validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  data.Form = userSignupForm{}
  app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
  var form userSignupForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
  form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
  form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
  form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
  form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
    return
  }

  err = app.users.Insert(form.Name, form.Email, form.Password)
  if err != nil {
    if errors.Is(err, models.ErrDuplicateEmail) {
      form.AddFieldError("email", "Email address is already in use")
      data := app.newTemplateData(r)
      data.Form = form
      app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
    } else {
      app.serveError(w, err)
    }
    
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

  http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
  Email               string `form:"email"` 
  Password            string `form:"password"` 
  validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  data.Form = userLoginForm{}
  app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
  var form userLoginForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
  form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
  form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "login.html", data)
    return
  }

  id, err := app.users.Authenticate(form.Email, form.Password)
  if err != nil {
    if errors.Is(err, models.ErrInvalidCredentials) {
      form.AddNonFieldError("Email or password is incorrect")

      data := app.newTemplateData(r)
      data.Form = form
      app.render(w, http.StatusUnprocessableEntity, "login.html", data)} else {
      app.serveError(w, err)
    }
    return
  }

  err = app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serveError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

  http.Redirect(w, r, "/habit/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
  err := app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serveError(w, err)
    return
  }

  app.sessionManager.Remove(r.Context(), "authenticatedUserID")

  app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) howTo(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  app.render(w, http.StatusOK, "howto.html", data)
}
