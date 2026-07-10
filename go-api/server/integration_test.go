package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api/config"
)

// newTestServer levanta un stub de la API en Node.js y construye el servidor apuntando a él.
// Devuelve el *Server y una función de limpieza.
func newTestServer(t *testing.T) (*Server, func()) {
	t.Helper()

	// Stub de Node: responde estadísticas fijas al endpoint /api/v1/statistics.
	nodeStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statistics" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"status":"success","data":{"max":175,"min":-70,"average":8.97,"sum":161.44,"isDiagonal":false}}`))
	}))

	cfg := &config.Config{
		ServerPort:         "0",
		JWTSecret:          "test-secret",
		AuthUser:           "admin",
		AuthPass:           "admin123",
		NodeAPIURL:         nodeStub.URL,
		TokenExpiryMinutes: 60,
		CORSOrigins:        "*",
	}

	return New(cfg), nodeStub.Close
}

// doJSON ejecuta una petición contra la app Fiber y devuelve status + cuerpo decodificado.
func doJSON(t *testing.T, s *Server, method, path, token string, body any) (int, map[string]any) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.app.Test(req, -1)
	if err != nil {
		t.Fatalf("error ejecutando la petición: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed map[string]any
	_ = json.Unmarshal(raw, &parsed)
	return resp.StatusCode, parsed
}

// login autentica y devuelve el token JWT.
func login(t *testing.T, s *Server) string {
	t.Helper()
	status, body := doJSON(t, s, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"username": "admin", "password": "admin123",
	})
	if status != http.StatusOK {
		t.Fatalf("login status = %d, se esperaba 200", status)
	}
	data := body["data"].(map[string]any)
	return data["token"].(string)
}

func TestIntegration_LoginAndProcess(t *testing.T) {
	s, cleanup := newTestServer(t)
	defer cleanup()

	token := login(t, s)

	// Flujo feliz: matriz válida -> QR + estadísticas.
	status, body := doJSON(t, s, http.MethodPost, "/api/v1/matrix/process", token, map[string]any{
		"matrix": [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
	})
	if status != http.StatusOK {
		t.Fatalf("process status = %d, se esperaba 200 (body: %v)", status, body)
	}

	data := body["data"].(map[string]any)
	if _, ok := data["qr"]; !ok {
		t.Error("la respuesta no incluye la factorización QR")
	}
	stats := data["statistics"].(map[string]any)
	if stats["max"].(float64) != 175 {
		t.Errorf("max = %v, se esperaba 175 (del stub de Node)", stats["max"])
	}
}

func TestIntegration_ProcessRequiresJWT(t *testing.T) {
	s, cleanup := newTestServer(t)
	defer cleanup()

	status, _ := doJSON(t, s, http.MethodPost, "/api/v1/matrix/process", "", map[string]any{
		"matrix": [][]float64{{1}},
	})
	if status != http.StatusUnauthorized {
		t.Errorf("sin token: status = %d, se esperaba 401", status)
	}
}

func TestIntegration_ProcessRejectsInvalidMatrix(t *testing.T) {
	s, cleanup := newTestServer(t)
	defer cleanup()

	token := login(t, s)

	// Matriz no rectangular -> 422 (error del cliente propagado desde la QR).
	status, _ := doJSON(t, s, http.MethodPost, "/api/v1/matrix/process", token, map[string]any{
		"matrix": [][]float64{{1, 2, 3}, {4, 5}},
	})
	if status != http.StatusUnprocessableEntity {
		t.Errorf("matriz inválida: status = %d, se esperaba 422", status)
	}
}

func TestIntegration_RejectsBadCredentials(t *testing.T) {
	s, cleanup := newTestServer(t)
	defer cleanup()

	status, _ := doJSON(t, s, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"username": "admin", "password": "wrong",
	})
	if status != http.StatusUnauthorized {
		t.Errorf("credenciales malas: status = %d, se esperaba 401", status)
	}
}
