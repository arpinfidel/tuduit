package http

import (
	"net/http"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/errs"
)

func (s *HTTPServer) AuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			NewJSONResponseErr("", errs.ErrUnauthorized.WithTrace().WithUserMessagef("Authorization header is required")).Response(w, r)
			return
		}

		bearer := "Bearer "
		n := len(bearer)
		if len(token) < n || token[:n] != bearer {
			NewJSONResponseErr("", errs.ErrUnauthorized.WithTrace().WithUserMessagef("Authorization header not in correct format")).Response(w, r)
			return
		}
		token = token[n:]

		claims, err := s.d.App.VerifyToken(token)
		if err != nil {
			NewJSONResponseErr("", errs.ErrUnauthorized.WithTrace().WithUserMessagef("Invalid token")).Response(w, r)
			return
		}

		users, _, err := s.d.App.GetUserByIDs(r.Context(), []int64{claims.UserID}, entity.Pagination{Page: 1, PageSize: 1})
		if err != nil {
			NewJSONResponseErr("", errs.ErrInternalServerError.WithTrace()).Response(w, r)
			return
		}
		if len(users) == 0 {
			NewJSONResponseErr("", errs.ErrInternalServerError.WithTrace()).Response(w, r)
			return
		}
		user := users[0]

		ctx := ctxx.New(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
