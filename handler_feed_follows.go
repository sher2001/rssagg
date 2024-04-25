package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sher2001/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("unable to create the feed_follow: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedFollowToCustomFeedFollow(feed_follow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("unable to fetch feed_followss: %v", err))
	}
	respondWithJSON(w, 200, databaseFeedFolowsToCustomFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follow_id_str := chi.URLParam(r, "feedFollowId")
	feed_follow_id, err := uuid.Parse(feed_follow_id_str)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("unable to parse id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feed_follow_id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("unable to delete feed_follow: %v", err))
	}

	respondWithJSON(w, 200, struct{}{})
}
