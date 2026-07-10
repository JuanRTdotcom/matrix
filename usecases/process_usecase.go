package usecases

import "go-api/models"

// Factorizer abstrae la factorización QR. La satisface QRUseCase.
type Factorizer interface {
	Execute(matrix [][]float64) (models.QRResult, error)
}

// StatsRepository es la dependencia externa (API de Node) que calcula estadísticas.
type StatsRepository interface {
	CalcularEstadisticas(token string, matrices [][][]float64) (models.Statistics, error)
}

// ProcessUseCase factoriza la matriz y pide las estadísticas de Q y R a la API de Node.
type ProcessUseCase struct {
	qr    Factorizer
	stats StatsRepository
}

// NewProcessUseCase construye el caso de uso de procesamiento.
func NewProcessUseCase(qr Factorizer, stats StatsRepository) *ProcessUseCase {
	return &ProcessUseCase{qr: qr, stats: stats}
}

// Execute factoriza la matriz y solicita las estadísticas de Q y R a la API de Node.
func (uc *ProcessUseCase) Execute(token string, matrix [][]float64) (models.ProcessResponse, error) {
	qr, err := uc.qr.Execute(matrix)
	if err != nil {
		return models.ProcessResponse{}, err
	}

	stats, err := uc.stats.CalcularEstadisticas(token, [][][]float64{qr.Q, qr.R})
	if err != nil {
		return models.ProcessResponse{}, err
	}

	return models.ProcessResponse{
		Original:   matrix,
		QR:         qr,
		Statistics: stats,
	}, nil
}
