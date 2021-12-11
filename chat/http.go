package chat

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

type key int

const (
	authenticatedUserKey key = 0
	conversationKey      key = 1
)

type HTTPListener struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewHTTPListener(db *gorm.DB, logger *log.Logger) *HTTPListener {
	return &HTTPListener{
		db:     db,
		logger: logger,
	}
}

func getAuthenticatedUserFromContext(ctx context.Context) (User, bool) {
	authenticatedUser, ok := ctx.Value(authenticatedUserKey).(User)
	return authenticatedUser, ok
}

func (l HTTPListener) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := User{}

		// TODO: add real authentication
		result := l.db.First(&user)
		if result.Error != nil {
			l.logger.Printf("error while getting first user: %v", result.Error)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		ctx := context.WithValue(r.Context(), authenticatedUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func (l HTTPListener) AutoMigrate() {
	l.db.AutoMigrate(&User{}, &Conversation{}, &ConversationMembership{}, &Message{})
}

func (l HTTPListener) Listen() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(l.authenticationMiddleware)

	convRs := conversationsResource{db: l.db, logger: log.Default()}
	r.Mount("/conversations", convRs.Routes())

	usersRs := usersResource{db: l.db, logger: log.Default()}
	r.Mount("/users", usersRs.Routes())

	l.logger.Print("Listening on Port 3000\n")
	http.ListenAndServe(":3000", r)
}
