package emailService

// The `go:generate` lines below (add as many as the number of templates you want to handle) will be processed by
// `go generate` command and will result in creation of Go sources with defined inflators.

//go:generate jade -pkg=emailService -stdlib -stdbuf templates/inviteToPrivateChatGroup.pug
