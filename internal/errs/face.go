package errs

type IApiError interface {
	error
	Status() int
}
