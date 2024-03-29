package tokenPayload

type TPayload struct {
	GrantType  string `json:"grant_type"`  // 'grant_type': 'urn:ietf:params:oauth:grant-type:device_code',
	DeviceCode string `json:"device_code"` //     'device_code': device_code_data['device_code'],
	ClientId   string `json:"client_id"`   //     'client_id': AUTH0_CLIENT_ID
}
