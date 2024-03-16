package main

import (
  "net/http"

  "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
  mux := http.NewServeMux()

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })

  mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))

  dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

  mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
  mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
  mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
  mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

  protected := dynamic.Append(app.requireAuthentication)

  mux.Handle("GET /{$}", protected.ThenFunc(app.home))
  mux.Handle("GET /habit/create", protected.ThenFunc(app.habitCreate))
  mux.Handle("POST /habit/create", protected.ThenFunc(app.habitCreatePost))
  mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  return standard.Then(mux)
}

