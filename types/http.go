package types

type HttpRequestType struct {
	ContentType      string
	Method           string
	FullURL          string
	HeaderAdditional []HttpHeaderAdditionalType
	Data             []byte
}

type HttpHeaderAdditionalType struct {
	Name  string
	Value string
}
