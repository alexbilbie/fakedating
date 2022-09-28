package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fakedating/pkg/middleware"
	"fakedating/pkg/model"
	"fakedating/pkg/payload"
	"fakedating/pkg/util"
)

func (h Handler) Swipe(w http.ResponseWriter, r *http.Request) {
	// Read payload
	body, readErr := io.ReadAll(io.LimitReader(r.Body, 500))
	defer r.Body.Close()
	if readErr != nil {
		util.WriteErrorResponse("Failed to read the request body", readErr, http.StatusInternalServerError, w)
		return
	}

	var swipePayload payload.SwipeRequest
	if unmarshallErr := json.Unmarshal(body, &swipePayload); unmarshallErr != nil {
		util.WriteErrorResponse("Failed to decode the request body", unmarshallErr, http.StatusBadRequest, w)
		return
	}

	// Validate recipient ID exists
	_, getRecipientErr := h.userRepository.GetByID(swipePayload.Recipient)
	if getRecipientErr != nil {
		log.Printf("Failed to lookup recipient user: %v", getRecipientErr)
		util.WriteErrorResponse("Failed to lookup recipient user", getRecipientErr, http.StatusInternalServerError, w)
		return
	}

	// Save swipe
	mutualMatch, saveErr := h.userRepository.SaveSwipe(
		middleware.GetUserIDFromContext(r.Context()),
		swipePayload.Recipient,
		swipePayload.Matched,
	)
	if saveErr != nil {
		log.Printf("Failed to save swipe: %v", saveErr)
		util.WriteErrorResponse("Failed to save swipe", saveErr, http.StatusInternalServerError, w)
		return
	}

	util.WriteJSONResponse(
		payload.SwipeResponse{MutualMatch: mutualMatch == model.ProfileMatchMutual},
		http.StatusOK,
		w,
	)
}
