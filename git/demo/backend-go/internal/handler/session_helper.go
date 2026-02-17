package handler

import (
	"fmt"
	"strconv"

	"github.com/gorilla/sessions"
)

func getUserIDFromSession(session *sessions.Session) (int64, error) {
	userIDRaw, exists := session.Values["user_id"]
	if !exists {
		return 0, fmt.Errorf("no user_id in session")
	}

	// Handle different types that might be stored
	switch v := userIDRaw.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse user_id string: %w", err)
		}
		return parsed, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unknown user_id type: %T", userIDRaw)
	}
}