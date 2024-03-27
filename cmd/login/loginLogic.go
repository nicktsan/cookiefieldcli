package login

import (
	"cookiefieldcli/cmd/login/loginResponse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetDeviceCode() (loginResponse.LResponse, error) {
	// url := "https://{yourDomain}/oauth/device/code"
	url := "https://dev-nucixn2420u6r4t4.us.auth0.com/oauth/device/code"

	// payload := strings.NewReader("client_id={yourClientId}&scope=%7Bopenid profile$7D")
	payload := strings.NewReader("client_id=I4d0XcAXPue9sFGTQMPEboZEyYsTZwBG&scope=%7Bopenid profile$7D")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

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
	fmt.Println("DeviceCodeData:")
	fmt.Println(deviceCodeData)
	payload := strings.NewReader("grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code&device_code=%7B" + deviceCodeData.DeviceCode + "%7D&client_id=I4d0XcAXPue9sFGTQMPEboZEyYsTZwBG")
	fmt.Println("Payload: ")
	fmt.Println(payload)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	var data map[string]interface{}
	authenticated := false
	for !authenticated {
		fmt.Println("Checking if the user completed the flow...")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Panic("Failed to post request token: ", err)
		}
		defer res.Body.Close()
		jsonErr := json.NewDecoder(res.Body).Decode(&data)
		if jsonErr != nil {
			log.Panic("Failed to decode res.Body to json: ", jsonErr)
		}
		if res.StatusCode == 200 {
			fmt.Println("Authenticated!")
			fmt.Println("- Id Token: { ", data["id_token"])
			authenticated = true
		} else if data["error"] != "authorization_pending" && data["error"] != "slow_down" {
			fmt.Println(data["error"])
			// fmt.Println(data["error_description"])
			log.Panic(data["error_description"])
		} else {
			time.Sleep(time.Duration(deviceCodeData.Interval) * time.Second)
		}

		// jsonErr := json.NewDecoder(res.Body).Decode(&data)
		// body, _ := io.ReadAll(res.Body)
		// fmt.Println("token data res: ")
		// fmt.Println(res)
		// fmt.Println("token data body: ")
		// fmt.Println(string(body))
		// fmt.Println("data interface: ", data)
	}
}
