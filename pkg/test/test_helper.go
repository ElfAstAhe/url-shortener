package test

type HttpTestCase struct {
	Name               string
	Method             string
	Path               string
	ExpectedStatusCode int
	ExpectedBody       string
	ExpectedStrValue   string
	ExpectedBytesValue []byte
}
