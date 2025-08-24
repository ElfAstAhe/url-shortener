package test

type HTTPTestCase struct {
	Name               string
	Method             string
	Path               string
	ExpectedStatusCode int
	ExpectedBody       string
	ExpectedStrValue   string
	ExpectedBytesValue []byte
}
