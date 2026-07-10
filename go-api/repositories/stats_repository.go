package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-api/models"
)

// StatsRepositoryImpl es el cliente HTTP hacia la API de Node que calcula las estadísticas.
type StatsRepositoryImpl struct {
	baseURL string
	client  *http.Client
}

// NewStatsRepository construye el repositorio con la URL base de la API en Node.js.
func NewStatsRepository(baseURL string) *StatsRepositoryImpl {
	return &StatsRepositoryImpl{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// CalcularEstadisticas envía las matrices a la API de Node, reenviando el JWT, y devuelve las estadísticas.
func (r *StatsRepositoryImpl) CalcularEstadisticas(token string, matrices [][][]float64) (models.Statistics, error) {
	payload := models.StatsRequest{Matrices: matrices}
	body, err := json.Marshal(payload)
	if err != nil {
		return models.Statistics{}, err
	}

	url := r.baseURL + "/api/v1/statistics"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return models.Statistics{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		return models.Statistics{}, fmt.Errorf("no se pudo contactar a la API de estadísticas: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return models.Statistics{}, fmt.Errorf("la API de estadísticas respondió %d: %s", resp.StatusCode, string(respBody))
	}

	// La API en Node.js responde con el envoltorio { data: { ...estadísticas } }.
	var wrapper struct {
		Data models.Statistics `json:"data"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return models.Statistics{}, fmt.Errorf("respuesta inválida de la API de estadísticas: %w", err)
	}

	return wrapper.Data, nil
}
