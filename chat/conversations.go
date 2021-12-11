package chat

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationsResource struct {
	db     *gorm.DB
	logger *log.Logger
}

func (rs conversationsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	/* r.Post("/", rs.Create) */

	r.Route("/{conversationID}", func(r chi.Router) {
		r.Use(rs.ConversationCtx)

		/*****
			r.Get("/", rs.Get)
			r.Patch("/", rs.Update)
			r.Delete("/", rs.Delete)

			r.Post("/join/{inviteLink}", rs.Join)
			r.Get("/members", rs.GetMembers)
			r.Get("/member/{memberID}", rs.GetMember)
			r.Patch("/member/{memberID}", rs.UpdateMember)

			r.Get("/messages", rs.GetMessages)
		*****/
	})

	return r
}

func (rs conversationsResource) ConversationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conversationID := chi.URLParam(r, "conversationID")
		conversationUUID, err := uuid.Parse(conversationID)
		if err != nil {
			rs.logger.Printf("error while converting conversation ID to UUID: %v", err)
			http.Error(w, "conversation ID must be a valid UUID", 400)
			return
		}

		conversation, err := getConversation(rs.db, conversationUUID)
		if err != nil {
			rs.logger.Printf("error while getting conversation: %v", err)

			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, http.StatusText(404), 404)
			} else {
				http.Error(w, "something went wrong", 500)
			}

			return
		}

		ctx := context.WithValue(r.Context(), conversationKey, conversation)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs conversationsResource) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := getAuthenticatedUserFromContext(ctx)

	if !ok {
		rs.logger.Printf("couldn't get user from context")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if user.Conversations == nil {
		user.Conversations = []*Conversation{}
	}

	marshaled, err := json.Marshal(user.Conversations)
	if err != nil {
		rs.logger.Printf("error while unmarshaling conversations: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(marshaled)
	return
}
