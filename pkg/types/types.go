package types

type SecretData struct {
	Id string
	Secret string
}

type SecretResponse struct {
	Data string `json:"data"`
}
