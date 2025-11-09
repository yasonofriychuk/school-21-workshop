package current_weather

import (
	"encoding/json"
	"github.com/yasonofriychuk/school-21-workshop/pkg/logger"
	"net/http"
)

type Handler struct {
	log     logger.Log
	usecase usecase
}

func New(log logger.Log, usecase usecase) *Handler {
	return &Handler{
		log:     log,
		usecase: usecase,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		h.log.WithContext(ctx).WithError(err).Error("failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	temperature, err := h.usecase.GetCurrenTemperatureByIp(r.Context(), requestData.IP)
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Error("failed to get current weather temperature")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]float64{"temperature": temperature}); err != nil {
		h.log.WithContext(ctx).WithError(err).Error("failed to encode response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
