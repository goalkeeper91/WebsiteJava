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

// Automation Settings Errors
var (
	ErrInvalidClipDurationRange = errors.New("min_clip_duration darf nicht größer als max_clip_duration sein")
	ErrClipDurationTooShort     = errors.New("min_clip_duration muss mindestens 5 sekunden betragen")
	ErrClipDurationTooLong      = errors.New("max_clip_duration darf 300 sekunden (5 minuten) nicht überschreiten")
	ErrInvalidHashtagCount      = errors.New("max_hashtags muss zwischen 0 und 30 liegen")
)

// Scheduled Message Errors
var (
	ErrScheduledMessageNotFound = errors.New("automatisierte nachricht nicht gefunden")
	ErrEmptyMessage             = errors.New("nachricht darf nicht leer sein")
	ErrIntervalTooShort         = errors.New("intervall muss mindestens 60 sekunden betragen")
)

// Giveaway Errors
var (
	ErrGiveawayAlreadyOpen    = errors.New("es läuft bereits ein giveaway für diesen kanal")
	ErrNoOpenGiveaway         = errors.New("kein offenes giveaway für diesen kanal")
	ErrNoGiveawayEntries      = errors.New("das giveaway hat noch keine teilnehmer")
	ErrGiveawayKeywordInvalid = errors.New("codewort muss zwischen 2 und 50 zeichen lang sein")
)

// General Errors
var (
	ErrUnauthorized   = errors.New("nicht autorisiert")
	ErrForbidden      = errors.New("zugriff verweigert")
	ErrInvalidInput   = errors.New("ungültige eingabe")
	ErrInternalServer = errors.New("interner serverfehler")
	ErrNotFound       = errors.New("nicht gefunden")
)