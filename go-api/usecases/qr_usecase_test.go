package usecases

import (
	"errors"
	"math"
	"testing"
)

// multiply calcula el producto de dos matrices para verificar A ≈ Q*R.
func multiply(a, b [][]float64) [][]float64 {
	m := len(a)
	n := len(b[0])
	inner := len(b)
	out := make([][]float64, m)
	for i := 0; i < m; i++ {
		out[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			var sum float64
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			out[i][j] = sum
		}
	}
	return out
}

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestQR_ReconstructsMatrix(t *testing.T) {
	uc := NewQRUseCase()
	a := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}

	res, err := uc.Execute(a)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// Q*R debe reconstruir A dentro de una tolerancia razonable.
	reconstructed := multiply(res.Q, res.R)
	for i := range a {
		for j := range a[i] {
			if !almostEqual(reconstructed[i][j], a[i][j], 1e-3) {
				t.Errorf("Q*R[%d][%d] = %v, se esperaba %v", i, j, reconstructed[i][j], a[i][j])
			}
		}
	}
}

func TestQR_OrthonormalColumns(t *testing.T) {
	uc := NewQRUseCase()
	a := [][]float64{{1, 2}, {3, 4}, {5, 6}}

	res, err := uc.Execute(a)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// Las columnas de Q deben ser ortonormales: q_i · q_j = 1 si i==j, 0 si i!=j.
	n := len(res.Q[0])
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var dot float64
			for r := 0; r < len(res.Q); r++ {
				dot += res.Q[r][i] * res.Q[r][j]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if !almostEqual(dot, expected, 1e-3) {
				t.Errorf("q%d·q%d = %v, se esperaba %v", i, j, dot, expected)
			}
		}
	}
}

func TestQR_UpperTriangularR(t *testing.T) {
	uc := NewQRUseCase()
	a := [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}}

	res, err := uc.Execute(a)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// R debe ser triangular superior: por debajo de la diagonal todo cero.
	for i := 0; i < len(res.R); i++ {
		for j := 0; j < i; j++ {
			if res.R[i][j] != 0 {
				t.Errorf("R[%d][%d] = %v, se esperaba 0 (triangular superior)", i, j, res.R[i][j])
			}
		}
	}
}

func TestQR_RejectsNonRectangular(t *testing.T) {
	uc := NewQRUseCase()
	a := [][]float64{{1, 2, 3}, {4, 5}}
	if _, err := uc.Execute(a); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("se esperaba ErrNotRectangular, se obtuvo %v", err)
	}
}

func TestQR_RejectsEmpty(t *testing.T) {
	uc := NewQRUseCase()
	if _, err := uc.Execute([][]float64{}); !errors.Is(err, ErrEmptyMatrix) {
		t.Errorf("se esperaba ErrEmptyMatrix, se obtuvo %v", err)
	}
}

func TestQR_RejectsRankDeficient(t *testing.T) {
	uc := NewQRUseCase()
	// Dos columnas idénticas → linealmente dependientes.
	a := [][]float64{{1, 1}, {2, 2}, {3, 3}}
	if _, err := uc.Execute(a); !errors.Is(err, ErrRankDeficient) {
		t.Errorf("se esperaba ErrRankDeficient, se obtuvo %v", err)
	}
}

// TestQR_ErrorsAreClientErrors verifica que los errores de entrada se reporten como
// errores del cliente (para que el controlador responda 422).
func TestQR_ErrorsAreClientErrors(t *testing.T) {
	if !ErrEmptyMatrix.IsClientError() || !ErrNotRectangular.IsClientError() {
		t.Error("los errores de entrada QR deben ser errores de cliente")
	}
}
