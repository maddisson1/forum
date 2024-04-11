package handlers

import "net/http"

func (h *handlers) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	mux.HandleFunc("/", h.cookieChecker(h.home))
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/post/create", h.isAuthenticated(h.createPost))
	mux.HandleFunc("/user/login", h.login)
	mux.HandleFunc("/user/signup", h.signup)
	mux.HandleFunc("/user/logout", h.isAuthenticated(h.logout))
	mux.HandleFunc("/user/posts", h.isAuthenticated(h.userPosts))

	return h.recoverPanic(h.logRequest(h.secureHeaders(mux)))
}
