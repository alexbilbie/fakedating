package payload

type ErrorResponse struct {
	Msg   string
	Error string `json:",omitempty"`
}
