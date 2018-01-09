package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/utrack/woofer/interface/ihttp"
	"github.com/utrack/woofer/lib/session/inmemsessions"
	"github.com/utrack/woofer/service"
)

func main() {
	logrus.Info("Wiring up services...")
	svc, err := service.Bootstrap(
		service.Config{
			SQLiteConnString: "./db.sqlite",
			SQLiteMigrations: "../../migrations",
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	sess := inmemsessions.New(time.Hour * 24)
	hdl := ihttp.NewHandler(svc, sess)

	r := chi.NewRouter()
	r.Use(middleware.Timeout(time.Second * 10))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(ihttp.UserAuthCtx(sess))

	r.Get("/", func(http.ResponseWriter, *http.Request) {
		// 200 healthcheck
	})

	r.Post("/user/create", hdl.UserCreate)
	r.Post("/auth", hdl.Login)
	r.Route("/", func(r chi.Router) {
		r.Use(ihttp.RequireAuth)
		r.Post("/tweet", hdl.Tweet)
		r.Get("/posts", hdl.GetTweetPage)
		r.Route("/u/{nickname}", func(r chi.Router) {
			r.Get("/", hdl.GetUser)
			r.Get("/tweets", hdl.GetProfileTweets)
			r.Get("/subscribe", hdl.Subscribe)
			r.Get("/unsubscribe", hdl.Unsubscribe)
		})
		r.Get("/subscriptions", hdl.Subscriptions)
		r.Get("/subscribers", hdl.Subscribers)
	})
	logrus.Info("Listening on :3333")
	http.ListenAndServe(":3333", r)

}
