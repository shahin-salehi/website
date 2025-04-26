package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/shahin-salehi/website/internal/db"
	"github.com/shahin-salehi/website/internal/handler"
	"github.com/shahin-salehi/website/internal/middleware"
	"github.com/shahin-salehi/website/internal/session"
	"github.com/shahin-salehi/website/internal/web/pages"

	"github.com/a-h/templ"
)

// use as middleware
func RequireAuth(sess *session.Session, next http.Handler) http.Handler {
	// wrap our function to satsify http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check session
		_, ok := sess.GetUserID(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// cont
		next.ServeHTTP(w, r)
	})
}

func main() {
	// init db
	connectionString := os.Getenv("DATABASE_URL")

	d, err := db.NewDatabase(connectionString)
	if err != nil {
		slog.Error("db init failed, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}
	// init repo
	repo := db.NewRepo(d)

	// init middleware
	mw := middleware.NewMiddleware(repo)

	// init session
	sess, err := session.InitStore()
	if err != nil {
		slog.Error("session init failed, shutting down.")
		os.Exit(1)
	}

	// init handler
	handler := handler.NewHandler(repo, sess)

	// create custom mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// pages
	mux.Handle("/", mw.InjectUser(sess, templ.Handler(pages.Index())))
	mux.Handle("/blog", mw.InjectUser(sess, templ.Handler(pages.Blog())))
	mux.Handle("/contact", mw.InjectUser(sess, mw.RequireAuth(sess, templ.Handler(pages.Contact()))))
	mux.Handle("/login", mw.RedirectAuth(sess, templ.Handler(pages.Login())))
	mux.Handle("/signup", mw.RedirectAuth(sess, templ.Handler(pages.Signup())))
	mux.Handle("/privacy", mw.InjectUser(sess, templ.Handler(pages.Privacy())))
	mux.Handle("/error", mw.InjectUser(sess, templ.Handler(pages.InternalServerError())))
	mux.Handle("/privacy-request", mw.InjectUser(sess, templ.Handler(pages.PrivacyRequest())))
	mux.Handle("/account-delete", mw.InjectUser(sess, mw.RequireAuth(sess, templ.Handler(pages.AccountDeleteConfirm()))))
	mux.Handle("/account-deleted", templ.Handler(pages.AccountDeleted()))
	mux.Handle("/data-delete", mw.InjectUser(sess, mw.RequireAuth(sess, templ.Handler(pages.DataDelete()))))
	mux.Handle("/data-deleted", mw.InjectUser(sess, mw.RequireAuth(sess, templ.Handler(pages.DataDeleted()))))

	// API handlers
	mux.HandleFunc("/api/login", handler.Login)
	mux.HandleFunc("/api/logout", handler.Logout)
	mux.HandleFunc("/api/signup", handler.Signup)
	mux.HandleFunc("/api/account-delete", handler.DeleteAccount)
	mux.HandleFunc("/api/data-request", handler.GetUserData)
	mux.HandleFunc("/api/data-delete", handler.DeleteUserData)
	mux.Handle("/api/message", mw.RequireAuth(sess, handler.Message()))

	// start server
	slog.Info("starting webserver...")
	http.ListenAndServe(":8080", mux)

}
