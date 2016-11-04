package hipchat

const (
	COLOR_YELLOW = "yellow"
	COLOR_GREEN  = "green"
	COLOR_RED    = "red"
	COLOR_PURPLE = "purple"
	COLOR_GRAY   = "gray"
	COLOR_RANDOM = "random"

	TYPE_TEXT = "text"
	TYPE_HTML = "html"
)

type Response struct {
	Color         string `json:"color"`
	From          string `json:"from"`
	AttachTo      string `json:"attach_to"`
	Message       string `json:"message"`
	Notify        bool   `json:"notify"`
	MessageFormat string `json:"message_format"`
}

func NewResponse(color, message string) *Response {
	return &Response{
		Color:         color,
		Message:       message,
		Notify:        false,
		MessageFormat: TYPE_TEXT,
	}
}

func NewNotify(color, message string) *Response {
	return &Response{
		Color:         color,
		From:          "Snakey!",
		Message:       message,
		Notify:        true,
		MessageFormat: TYPE_TEXT,
	}
}
