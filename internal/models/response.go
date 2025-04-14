package models

type APIResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommandResponse struct {
	Executed      bool   `json:"executed"`
	ExitCode      int    `json:"exit_code,omitempty"`
	Stdout        string `json:"stdout,omitempty"`
	Stderr        string `json:"stderr,omitempty"`
	ExecutionTime string `json:"execution_time,omitempty"`
}

type Response struct {
	API     APIResponse     `json:"api"`
	Command CommandResponse `json:"command"`
}

type CommandsResponse struct {
	API      APIResponse `json:"api"`
	Commands []Command   `json:"commands"`
}

// Creates a response for successful execution
func NewSuccessResponse(executed bool, exitCode int, stdout, stderr, execTime string) *Response {
	return &Response{
		API: APIResponse{
			Status:  "success",
			Code:    "SUCCESS",
			Message: "Request processed successfully",
		},
		Command: CommandResponse{
			Executed:      executed,
			ExitCode:      exitCode,
			Stdout:        stdout,
			Stderr:        stderr,
			ExecutionTime: execTime,
		},
	}
}

// Creates a response for API processing errors
func NewAPIErrorResponse(code, message string) *Response {
	return &Response{
		API: APIResponse{
			Status:  "error",
			Code:    code,
			Message: message,
		},
		Command: CommandResponse{
			Executed: false,
		},
	}
}

// Creates a response for successful command listing
func NewCommandsSuccessResponse(commands []Command) *CommandsResponse {
	return &CommandsResponse{
		API: APIResponse{
			Status:  "success",
			Code:    "SUCCESS",
			Message: "Request processed successfully",
		},
		Commands: commands,
	}
}

// Error code constants
const (
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeCommandNotAllowed  = "COMMAND_NOT_ALLOWED"
	ErrCodeArgumentNotAllowed = "ARGUMENT_NOT_ALLOWED"
	ErrCodeCommandTimeout     = "COMMAND_TIMEOUT"
	ErrCodeServerError        = "SERVER_ERROR"
)
