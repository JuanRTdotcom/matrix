package usecases

import (
	"math"

	"go-api/models"
)

// QRInputError representa una matriz de entrada inválida (error del cliente).
type QRInputError struct {
	Message string
}

func (e *QRInputError) Error() string       { return e.Message }
func (e *QRInputError) IsClientError() bool { return true }

// Errores de entrada de la factorización QR.
var (
	ErrEmptyMatrix    = &QRInputError{Message: "la matriz no puede estar vacía"}
	ErrNotRectangular = &QRInputError{Message: "la matriz debe ser rectangular (todas las filas del mismo tamaño)"}
	ErrRankDeficient  = &QRInputError{Message: "la matriz tiene columnas linealmente dependientes; no admite factorización QR estable"}
)

// QRUseCase encapsula la lógica de negocio para la factorización QR de una matriz.
type QRUseCase struct{}

// NewQRUseCase construye el caso de uso de factorización QR.
func NewQRUseCase() *QRUseCase {
	return &QRUseCase{}
}

// Validate verifica que la matriz sea no vacía y rectangular.
func (uc *QRUseCase) Validate(a [][]float64) error {
	if len(a) == 0 || len(a[0]) == 0 {
		return ErrEmptyMatrix
	}
	cols := len(a[0])
	for _, row := range a {
		if len(row) != cols {
			return ErrNotRectangular
		}
	}
	return nil
}

// Execute calcula la factorización QR (Q y R) con el algoritmo de Gram-Schmidt.
func (uc *QRUseCase) Execute(a [][]float64) (models.QRResult, error) {
	if err := uc.Validate(a); err != nil {
		return models.QRResult{}, err
	}

	m := len(a)    // filas
	n := len(a[0]) // columnas

	// v[j] es la columna j de A; se va ortogonalizando.
	v := make([][]float64, n)
	for j := 0; j < n; j++ {
		v[j] = make([]float64, m)
		for i := 0; i < m; i++ {
			v[j][i] = a[i][j]
		}
	}

	// Q se guarda por columnas para facilitar los productos punto.
	qCols := make([][]float64, n)
	r := makeMatrix(n, n)

	const eps = 1e-12
	for j := 0; j < n; j++ {
		norm := vectorNorm(v[j])
		if norm < eps {
			return models.QRResult{}, ErrRankDeficient
		}
		r[j][j] = norm

		// Normaliza la columna j.
		qCols[j] = make([]float64, m)
		for i := 0; i < m; i++ {
			qCols[j][i] = v[j][i] / norm
		}

		// Ortogonaliza las columnas restantes contra q_j.
		for k := j + 1; k < n; k++ {
			dot := dotProduct(qCols[j], v[k])
			r[j][k] = dot
			for i := 0; i < m; i++ {
				v[k][i] -= dot * qCols[j][i]
			}
		}
	}

	// Reconstruye Q en formato fila (m x n) para la respuesta.
	q := makeMatrix(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			q[i][j] = round(qCols[j][i])
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			r[i][j] = round(r[i][j])
		}
	}

	return models.QRResult{Q: q, R: r}, nil
}

// makeMatrix crea una matriz de ceros de dimensiones rows x cols.
func makeMatrix(rows, cols int) [][]float64 {
	mtx := make([][]float64, rows)
	for i := range mtx {
		mtx[i] = make([]float64, cols)
	}
	return mtx
}

// vectorNorm calcula la norma euclidiana de un vector.
func vectorNorm(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

// dotProduct calcula el producto punto de dos vectores de igual longitud.
func dotProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// round redondea a 6 decimales para evitar ruido de punto flotante en la respuesta.
func round(x float64) float64 {
	const factor = 1e6
	r := math.Round(x*factor) / factor
	if r == 0 {
		return 0 // evita el -0
	}
	return r
}
