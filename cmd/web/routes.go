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
  mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
  mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
  mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
  mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
  mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
  mux.Handle("GET /howTo", dynamic.ThenFunc(app.howTo))
  mux.Handle("POST /habit/log/{id}/{day}", dynamic.ThenFunc(app.habitLogPost))
  // mux.Handle("POST /habit/log/{id}/{day}", dynamic.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "Received POST request")}))

  protected := dynamic.Append(app.requireAuthentication)

  mux.Handle("GET /habit/create", protected.ThenFunc(app.habitCreate))
  mux.Handle("POST /habit/create", protected.ThenFunc(app.habitCreatePost))

  mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  return standard.Then(mux)
}
