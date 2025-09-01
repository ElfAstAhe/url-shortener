package test

type HTTPTestCase struct {
	Name               string
	Method             string
	Path               string
	Body               []byte
	ExpectedStatusCode int
	ExpectedBody       string
	ExpectedStrValue   string
	ExpectedBytesValue []byte
}
