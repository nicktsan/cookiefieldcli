package login

import (
	"cookiefieldcli/cmd/login/loginResponse"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func RequestDeviceCode() (string, error) {
	url := "https://{yourDomain}/oauth/device/code"

	payload := strings.NewReader("client_id={yourClientId}&scope=%7Bopenid profile$7D")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Error generating the device code")
		os.Exit(1)
	}
	fmt.Println("Device code successful")

	//Convert the response body into a []byte so we can prepare it for conversion to a custom Response struct.
	body, _ := io.ReadAll(res.Body)
	// fmt.Println("res:")
	// fmt.Println(res)
	// fmt.Println("res.Body:")
	// fmt.Println(string(body))

	//Converts the body from []byte to a custom Response struct.
	var result loginResponse.Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println("Unmarshalled response body:")
	fmt.Println(result)

	return string(body), nil
}
