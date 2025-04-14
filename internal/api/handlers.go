package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zinrai/sevalet/internal/config"
	"github.com/zinrai/sevalet/internal/executor"
	"github.com/zinrai/sevalet/internal/models"
)

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		config: config,
	}
}

// Processes command execution requests
func (h *Handler) ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	request, err := models.NewExecuteRequestFromJSON(r.Body)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		respondWithError(w, models.ErrCodeInvalidRequest, err.Error())
		return
	}

	// Log the command request
	log.Printf("Command request: %s %v (timeout: %ds)", request.Command, request.Args, request.Timeout)

	// Validate request
	if err := request.Validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		respondWithError(w, models.ErrCodeInvalidRequest, err.Error())
		return
	}

	// Validate if command is allowed
	if err := executor.ValidateCommand(request.Command, request.Args, h.config.GetCommandList()); err != nil {
		code := models.ErrCodeCommandNotAllowed
		if len(err.Error()) >= 5 && err.Error()[0:5] == "argum" {
			code = models.ErrCodeArgumentNotAllowed
		}
		log.Printf("Command validation error: %v", err)
		respondWithError(w, code, err.Error())
		return
	}

	// Execute command
	result := executor.ExecuteCommand(r.Context(), request.Command, request.Args, request.Timeout)

	// Error handling
	if result.Error != nil {
		if result.ExitCode == -1 { // Timeout
			log.Printf("Command timed out: %s %v", request.Command, request.Args)
			respondWithError(w, models.ErrCodeCommandTimeout, result.Error.Error())
		} else {
			log.Printf("Command execution error: %v", result.Error)
			respondWithError(w, models.ErrCodeServerError, result.Error.Error())
		}
		return
	}

	// Log command execution result
	log.Printf("Command executed: %s %v (exit: %d, time: %s)",
		request.Command, request.Args, result.ExitCode, result.ExecutionTime)

	// Return success response
	response := models.NewSuccessResponse(true, result.ExitCode, result.Stdout, result.Stderr, result.ExecutionTime)
	respondWithJSON(w, http.StatusOK, response)
}

// Returns the list of allowed commands
func (h *Handler) CommandsHandler(w http.ResponseWriter, r *http.Request) {
	// Get command list
	commandList := h.config.GetCommandList()
	if commandList == nil {
		log.Printf("Command list not loaded")
		respondWithError(w, models.ErrCodeServerError, "Command list not loaded")
		return
	}

	log.Printf("Commands list requested, returning %d commands", len(commandList.Commands))

	// Return success response
	response := models.NewCommandsSuccessResponse(commandList.Commands)
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code string, message string) {
	response := models.NewAPIErrorResponse(code, message)
	respondWithJSON(w, getHTTPStatusForErrorCode(code), response)
}

// Returns the HTTP status code corresponding to the error code
func getHTTPStatusForErrorCode(code string) int {
	switch code {
	case models.ErrCodeInvalidRequest:
		return http.StatusBadRequest
	case models.ErrCodeCommandNotAllowed, models.ErrCodeArgumentNotAllowed:
		return http.StatusForbidden
	case models.ErrCodeCommandTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
