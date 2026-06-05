package errors

import "net/http"

var (
	ErrRoundNotFound     = New("round not found", http.StatusNotFound)
	ErrServerNotFound    = New("server not found", http.StatusNotFound)
	ErrServerExists      = New("server already exists", http.StatusConflict)
	ErrJsonInvalidFormat = New("invalid json format", http.StatusBadRequest)
	ErrPlayerNotFound    = New("player not found", http.StatusNotFound)
	ErrPlayerExists      = New("player already exists", http.StatusConflict)
	ErrUUIDError         = New("uuid is in an invalid format", http.StatusBadRequest)

	ErrNoKillsReceived  = New("no kill records were received", http.StatusBadRequest)
	ErrInvalidKillerID  = New("killer id is invalid or missing", http.StatusBadRequest)
	ErrInvalidVictimID  = New("victim id is invalid or missing", http.StatusBadRequest)
	ErrInvalidServerID  = New("server id is invalid or missing", http.StatusBadRequest)
	ErrInvalidRoundID   = New("round id is invalid or missing", http.StatusBadRequest)
	ErrServerNotFoundDB = New("the specified server was not found", http.StatusNotFound)
	ErrRoundNotFoundDB  = New("the specified round was not found", http.StatusNotFound)
	ErrPlayerLookupFail = New("failed to retrieve player data", http.StatusInternalServerError)
	ErrBatchSaveFailed  = New("failed to save kills batch", http.StatusInternalServerError)

	ErrRoundsNotFound = New("no rounds found", http.StatusNotFound)
	ErrStatsNotFound  = New("stats not found", http.StatusNotFound)
	ErrConvertParam   = New("failed to convert parameter", http.StatusBadRequest)

	ErrHeaderMissing = New("required header is missing", http.StatusBadRequest)
	ErrUnauthorized  = New("unauthorized", http.StatusUnauthorized)

	ErrJSONMarshalFail = New("failed to serialize json", http.StatusInternalServerError)

	ErrRedisOperationFail = New("failed to execute redis operation", http.StatusInternalServerError)
)
