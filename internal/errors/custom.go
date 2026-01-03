package errors

import "net/http"

var (
	ErrRoundNotFound     = New("round não encontrado", http.StatusNotFound)
	ErrServerNotFound    = New("servidor não encontrado", http.StatusNotFound)
	ErrServerExists      = New("servidor já existente", http.StatusConflict)
	ErrJsonInvalidFormat = New("formato do json inválido", http.StatusBadRequest)
	ErrPlayerNotFound    = New("player não encontrado", http.StatusNotFound)
	ErrPlayerExists      = New("player já existente", http.StatusConflict)
	ErrUUIDError         = New("o uuid está em um padrão inválido", http.StatusInternalServerError)

	ErrNoKillsReceived  = New("nenhum registro de kill foi recebido", http.StatusBadRequest)
	ErrInvalidKillerID  = New("o id do killer está inválido ou ausente", http.StatusBadRequest)
	ErrInvalidVictimID  = New("o id da vítima está inválido ou ausente", http.StatusBadRequest)
	ErrInvalidServerID  = New("o id do servidor está inválido ou ausente", http.StatusBadRequest)
	ErrInvalidRoundID   = New("o id do round está inválido ou ausente", http.StatusBadRequest)
	ErrServerNotFoundDB = New("o servidor informado não foi encontrado", http.StatusNotFound)
	ErrRoundNotFoundDB  = New("o round informado não foi encontrado", http.StatusNotFound)
	ErrPlayerLookupFail = New("falha ao buscar dados do player", http.StatusInternalServerError)
	ErrBatchSaveFailed  = New("falha ao salvar batch de kills", http.StatusInternalServerError)

	ErrRoundsNotFound = New("nenhuma rodada encontrada", http.StatusNotFound)
	ErrStatsNotFound  = New("estatísticas não encontradas", http.StatusNotFound)
	ErrConvertParam   = New("erro ao converter parâmetro", http.StatusBadRequest)

	ErrHeaderMissing = New("cabeçalho obrigatório ausente", http.StatusBadRequest)
	ErrUnauthorized  = New("não autorizado", http.StatusUnauthorized)

	ErrJSONMarshalFail = New("falha ao serializar json", http.StatusInternalServerError)

	ErrRedisOperationFail = New("falha ao executar operação no redis", http.StatusInternalServerError)
)
