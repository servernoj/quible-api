package emailService

// The `go:generate` lines below (add as many as the number of templates you want to handle) will be processed by
// `go generate ./...` command and will result in creation of Go sources with defined inflators.

//go:generate jade -pkg=emailService -stdlib -stdbuf templates/userActivation.pug
//go:generate jade -pkg=emailService -stdlib -stdbuf templates/passwordReset.pug
//go:generate jade -pkg=emailService -stdlib -stdbuf templates/userInvitation.pug

func ternary(condition bool, iftrue, iffalse any) any {
	if condition {
		return iftrue
	} else {
		return iffalse
	}
}
