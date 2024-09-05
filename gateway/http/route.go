package http

import "github.com/go-chi/chi/v5"

func (s *HTTPServer) route(r chi.Router) {
	api := r.Route("/api", func(r chi.Router) {})
	api.Route("/v1", s.v1)
}

func (s *HTTPServer) v1(r chi.Router) {
	usr := r.Route("/user", func(r chi.Router) {})
	usr.Route("/registration", func(r chi.Router) {
		r.Post("/otp/send", wrapHandler(s, s.d.App.OTPSend))
		r.Post("/otp/verify", wrapHandler(s, s.d.App.OTPVerify))
		r.Post("/register", wrapHandler(s, s.d.App.UserRegister))
	})
	usr.Route("/auth", func(r chi.Router) {
		r.Post("/login", wrapHandler(s, s.d.App.UserLogin))
	})

	task := r.Route("/task", func(r chi.Router) {})
	task.Route("/", func(r chi.Router) {
		r.With(s.AuthJWT).Get("/", wrapHandler(s, s.d.App.GetTaskList))
	})
}
