package v1

const (
	GROUP_PREFIX string = "chat:"
)

var (
	AccessReadOnly  = []string{"subscribe", "history"}
	AccessReadWrite = []string{"subscribe", "publish", "history"}
)
