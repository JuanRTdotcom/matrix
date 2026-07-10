package models

// MatrixRequest es la matriz de entrada del endpoint de procesamiento.
// @Description Matriz rectangular de entrada.
type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix" validate:"required,min=1"`
}

// QRResult contiene la factorización QR (matrices Q y R) de la matriz de entrada.
// @Description Resultado de la factorización QR.
type QRResult struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// Statistics son las estadísticas de Q y R calculadas por la API de Node.
// @Description Estadísticas de las matrices resultantes.
type Statistics struct {
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	Average    float64 `json:"average"`
	Sum        float64 `json:"sum"`
	IsDiagonal bool    `json:"isDiagonal"`
}

// ProcessResponse es la respuesta: matriz original, su factorización QR y las estadísticas.
// @Description Respuesta de procesamiento.
type ProcessResponse struct {
	Original   [][]float64 `json:"original"`
	QR         QRResult    `json:"qr"`
	Statistics Statistics  `json:"statistics"`
}

// StatsRequest es el payload que la API en Go envía a la API en Node.js.
// @Description Payload de matrices para calcular estadísticas.
type StatsRequest struct {
	Matrices [][][]float64 `json:"matrices"`
}
