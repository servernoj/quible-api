package emailService

//go:generate jade -pkg=emailService -stdlib -stdbuf templates/activation.pug

var Handlers = map[string]any{
	"Activation": Activation,
}
