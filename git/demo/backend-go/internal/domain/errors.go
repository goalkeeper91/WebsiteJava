package domain

import "errors"

// User Errors
var (
	ErrUserNotFound      = errors.New("user nicht gefunden")
	ErrUserAlreadyExists = errors.New("user existiert bereits")
)

// Auth Token Errors
var (
	ErrAuthTokenNotFound = errors.New("auth token nicht gefunden")
	ErrTokenExpired      = errors.New("token abgelaufen")
)

// Chat Command Errors
var (
	ErrCommandNotFound      = errors.New("command nicht gefunden")
	ErrCommandAlreadyExists = errors.New("command existiert bereits")
	ErrEmptyTrigger         = errors.New("trigger darf nicht leer sein")
	ErrTriggerTooLong       = errors.New("trigger zu lang (max 100 Zeichen)")
	ErrTriggerTaken         = errors.New("trigger bereits vergeben")
)

// Channel Errors
var (
	ErrChannelNotFound      = errors.New("channel nicht gefunden")
	ErrChannelNotRegistered = errors.New("channel nicht registriert")
	ErrChannelAlreadyExists = errors.New("channel existiert bereits")
)

// Subscription Errors
var (
	ErrSubscriptionNotFound  = errors.New("subscription nicht gefunden")
	ErrSubscriptionExpired   = errors.New("subscription abgelaufen")
	ErrSubscriptionCanceled  = errors.New("subscription wurde gekündigt")
	ErrFeatureNotAvailable   = errors.New("feature nicht in aktuellem plan verfügbar")
	ErrCommandLimitReached   = errors.New("command limit des plans erreicht")
	ErrWorkflowLimitReached  = errors.New("workflow limit des plans erreicht")
	ErrVoteLimitReached      = errors.New("vote limit des plans für diesen monat erreicht")
)

// n8n Integration Errors
var (
	ErrN8NIntegrationNotFound    = errors.New("n8n integration nicht gefunden")
	ErrN8NIntegrationDisabled    = errors.New("n8n integration ist deaktiviert")
	ErrN8NIntegrationNotReady    = errors.New("n8n integration nicht konfiguriert")
	ErrN8NWebhookFailed          = errors.New("n8n webhook aufruf fehlgeschlagen")
	ErrN8NWorkflowNotFound       = errors.New("n8n workflow nicht gefunden")
)

// Vote Errors
var (
	ErrVoteSessionNotFound   = errors.New("vote session nicht gefunden")
	ErrVoteSessionClosed     = errors.New("vote session ist bereits geschlossen")
	ErrVoteSessionNotActive  = errors.New("vote session ist nicht aktiv")
	ErrAlreadyVoted          = errors.New("du hast bereits in dieser session abgestimmt")
	ErrInvalidVoteOption     = errors.New("ungültige vote option")
	ErrNoVotesYet            = errors.New("noch keine stimmen abgegeben")
)

// Workflow Template Errors
var (
	ErrTemplateNotFound    = errors.New("workflow template nicht gefunden")
	ErrTemplateNotAllowed  = errors.New("workflow template erfordert höheren plan")
)

// General Errors
var (
	ErrUnauthorized   = errors.New("nicht autorisiert")
	ErrForbidden      = errors.New("zugriff verweigert")
	ErrInvalidInput   = errors.New("ungültige eingabe")
	ErrInternalServer = errors.New("interner serverfehler")
	ErrNotFound       = errors.New("nicht gefunden")
)