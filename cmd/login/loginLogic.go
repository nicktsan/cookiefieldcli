package login

import (
	"cookiefieldcli/cmd/login/loginResponse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetDeviceCode() (loginResponse.LResponse, error) {
	// url := "https://{yourDomain}/oauth/device/code"
	url := "https://dev-nucixn2420u6r4t4.us.auth0.com/oauth/device/code"

	// payload := strings.NewReader("client_id={yourClientId}&scope=%7Bopenid profile$7D")
	payload := strings.NewReader("client_id=I4d0XcAXPue9sFGTQMPEboZEyYsTZwBG&scope=%7Bopenid profile$7D")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	// fmt.Println("GetDeviceCode request Header")
	// fmt.Println(req.Header)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Failed to retreive device code.")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// fmt.Println("Error generating the device code")
		log.Panic("Error generating the device code.")
	}
	fmt.Println("Device code successful")

	//Convert the response body into a []byte so we can prepare it for conversion to a custom Response struct.
	body, _ := io.ReadAll(res.Body)
	// fmt.Println("res:")
	// fmt.Println(res)
	// fmt.Println("res.Body:")
	// fmt.Println(string(body))

	//Converts the body from []byte to a custom Response struct.
	var result loginResponse.LResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	// fmt.Println("Unmarshalled response body:")
	// fmt.Println(result)
	fmt.Println("1. On your computer or mobile device navigate to: ", result.VerificationURIComplete)
	fmt.Println("2. Enter the following code: ", result.UserCode)

	return result, nil
}

func PostRequestToken(deviceCodeData loginResponse.LResponse) {
	url := "https://dev-nucixn2420u6r4t4.us.auth0.com/oauth/token"
	method := "POST"
	payload := strings.NewReader("grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code&device_code=" + deviceCodeData.DeviceCode + "&client_id=I4d0XcAXPue9sFGTQMPEboZEyYsTZwBG")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
