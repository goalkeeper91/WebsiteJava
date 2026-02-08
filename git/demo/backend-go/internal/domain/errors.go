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
	ErrChannelNotFound        = errors.New("channel nicht gefunden")
	ErrChannelNotRegistered   = errors.New("channel nicht registriert")
	ErrChannelAlreadyExists   = errors.New("channel existiert bereits")
)

// General Errors
var (
	ErrUnauthorized       = errors.New("nicht autorisiert")
	ErrForbidden          = errors.New("zugriff verweigert")
	ErrInvalidInput       = errors.New("ungültige eingabe")
	ErrInternalServer     = errors.New("interner serverfehler")
)