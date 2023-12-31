package emailService

// This file defines association between PUG templates and their corresponding data inflators. The templates
// themselves are written to declare their inflator signature (read more about `filter` block in templates/README.md)
// so the purpose of this file is to "publish" those inflators in `Handlers` map accessible elsewhere.
//
// The `go:generate` lines below (add as many as the number of templates you want to handle) will be processed by
// `go generate` command and will result in creation of Go sources with defined handlers/inflators.

//go:generate jade -pkg=emailService -stdlib -stdbuf templates/demo.pug
//go:generate jade -pkg=emailService -stdlib -stdbuf templates/userActivation.pug

var Handlers = map[string]any{
	"Demo":           Demo,
	"UserActivation": UserActivation,
}

func ternary(condition bool, iftrue, iffalse any) any {
	if condition {
		return iftrue
	} else {
		return iffalse
	}
}
