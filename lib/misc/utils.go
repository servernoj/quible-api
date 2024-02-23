package misc

func Of[E any](e E) *E {
	return &e
}
