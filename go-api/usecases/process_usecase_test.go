package usecases

import (
	"errors"
	"testing"

	"go-api/models"
)

// --- Dobles de prueba (mocks) ---

// stubFactorizer devuelve una QR fija o un error, sin cálculo real.
type stubFactorizer struct {
	result models.QRResult
	err    error
}

func (s *stubFactorizer) Execute(_ [][]float64) (models.QRResult, error) {
	return s.result, s.err
}

// spyStatsRepo registra lo que recibe y devuelve una respuesta fija.
type spyStatsRepo struct {
	gotToken    string
	gotMatrices [][][]float64
	result      models.Statistics
	err         error
}

func (s *spyStatsRepo) CalcularEstadisticas(token string, matrices [][][]float64) (models.Statistics, error) {
	s.gotToken = token
	s.gotMatrices = matrices
	return s.result, s.err
}

// --- Tests ---

func TestProcess_ComposesQRAndStats(t *testing.T) {
	qr := models.QRResult{Q: [][]float64{{1, 0}, {0, 1}}, R: [][]float64{{2, 0}, {0, 3}}}
	factorizer := &stubFactorizer{result: qr}
	stats := &spyStatsRepo{result: models.Statistics{Max: 3, Min: 0, Sum: 6, Average: 1.5, IsDiagonal: true}}

	uc := NewProcessUseCase(factorizer, stats)
	res, err := uc.Execute("tok-123", [][]float64{{2, 0}, {0, 3}})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// El caso de uso debe reenviar el token y ambas matrices (Q y R) a Node.
	if stats.gotToken != "tok-123" {
		t.Errorf("token reenviado = %q, se esperaba tok-123", stats.gotToken)
	}
	if len(stats.gotMatrices) != 2 {
		t.Fatalf("se esperaban 2 matrices (Q y R), se enviaron %d", len(stats.gotMatrices))
	}

	// La respuesta debe combinar QR + estadísticas.
	if res.Statistics.Max != 3 || !res.Statistics.IsDiagonal {
		t.Errorf("estadísticas no propagadas correctamente: %+v", res.Statistics)
	}
	if len(res.QR.Q) != 2 {
		t.Errorf("la QR no se incluyó en la respuesta")
	}
}

func TestProcess_PropagatesQRError(t *testing.T) {
	factorizer := &stubFactorizer{err: ErrNotRectangular}
	stats := &spyStatsRepo{}

	uc := NewProcessUseCase(factorizer, stats)
	_, err := uc.Execute("tok", [][]float64{{1, 2}, {3}})

	if !errors.Is(err, ErrNotRectangular) {
		t.Errorf("se esperaba ErrNotRectangular, se obtuvo %v", err)
	}
	// Si la QR falla, NO debe llamarse a la API de estadísticas.
	if stats.gotToken != "" {
		t.Error("no se debía llamar a la API de estadísticas cuando la QR falla")
	}
}

func TestProcess_PropagatesStatsError(t *testing.T) {
	qr := models.QRResult{Q: [][]float64{{1}}, R: [][]float64{{1}}}
	factorizer := &stubFactorizer{result: qr}
	statsErr := errors.New("node caído")
	stats := &spyStatsRepo{err: statsErr}

	uc := NewProcessUseCase(factorizer, stats)
	_, err := uc.Execute("tok", [][]float64{{1}})

	if !errors.Is(err, statsErr) {
		t.Errorf("se esperaba el error de stats, se obtuvo %v", err)
	}
}
