package login

import (
	// loginInterface "cookiefieldcli/cmd/login/interface"
	"cookiefieldcli/cmd/login/loginResponse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJob struct {
	// iLogin   loginInterface.ILogin
	DeviceCodeData loginResponse.LResponse
	ClientID       string
	Domain         string
}

func NewLoginJob(clientid string, domain string) *LoginJob {
	return &LoginJob{
		ClientID: clientid,
		Domain:   domain,
	}
}

func (loginJob *LoginJob) SetDeviceCodeData(deviceCodeData loginResponse.LResponse) {
	loginJob.DeviceCodeData = deviceCodeData
}

// Gets a device code to be used later when getting a request token.
func (loginJob *LoginJob) GetDeviceCode() error {
	//Set up an http request to get a device code.
	url := loginJob.Domain + "/oauth/device/code"
	payload := strings.NewReader("client_id=" + loginJob.ClientID + "&scope=%7Bopenid profile$7D")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	// fmt.Println("GetDeviceCode request Header")
	// fmt.Println(req.Header)
	var result loginResponse.LResponse
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Failed to retreive device code.")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
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
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Panic("Can not unmarshal JSON")
	}
	// fmt.Println("Unmarshalled response body:")
	// fmt.Println(result)
	//Prompt the user to open a link for logging in.
	fmt.Println("1. On your computer or mobile device navigate to: ", result.VerificationURIComplete)
	fmt.Println("2. Enter the following code: ", result.UserCode)
	//Assign the Device Code Data to the loginJob so we don't have to pass it as a parameter multiple times in other functions.
	loginJob.SetDeviceCodeData(result)
	return nil
}

// Gets a request token.
func (loginJob *LoginJob) GetRequestToken() error {
	url := loginJob.Domain + "/oauth/token"
	method := "POST"
	authenticate := false
	var pollErr error
	//Keep polling Auth0 for a request token until status 200 or until invalid grant
	for !authenticate {
		authenticate, pollErr = loginJob.PollRequestTokenStatus(url, method)
		if pollErr != nil {
			log.Panic(pollErr)
		}
		//Sleep for the interval duration in the device code data. If we poll too fast, Auth0 will give 429 status.
		time.Sleep(time.Duration(loginJob.DeviceCodeData.Interval) * time.Second)
	}
	return nil
}

func (loginJob *LoginJob) PollRequestTokenStatus(url string, method string) (bool, error) {
	//Construct a new Reader. This must be done within this function, or else Go will read to the end of the REader and not
	//properly construct our http request.
	payload := strings.NewReader("grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code&device_code=" +
		loginJob.DeviceCodeData.DeviceCode + "&client_id=" + loginJob.ClientID)

	//Create a new http Client. This is done within this function because Go will not erase previous Client settings if we
	//pass a client instance as a parameter, which will mess up our post request.
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	//Construct a new http request. This must be done within this function because Go will nuke the request headers after every
	//call of &http.Client.Do().
	req, reqErr := http.NewRequest(method, url, payload)
	if reqErr != nil {
		fmt.Println(reqErr)
		return false, reqErr
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, resErr := client.Do(req)
	// var token_data map[string]interface{}
	if resErr != nil {
		fmt.Println(resErr)
		return false, resErr
	}
	defer res.Body.Close()
	// BodyDecodeErr := json.NewDecoder(res.Body).Decode(&token_data)
	// if BodyDecodeErr != nil {
	// 	fmt.Println("Could not convert res.Body to json: ")
	// 	log.Panic(BodyDecodeErr)
	// 	return false
	// }
	// resbody, ReadAllErr := io.ReadAll(res.Body)
	// if ReadAllErr != nil {
	// 	fmt.Println(ReadAllErr)
	// 	return false, ReadAllErr
	// }
	fmt.Println("res.StatusCode:")
	fmt.Println(res.StatusCode)
	// fmt.Println("res.Body: ")
	// fmt.Println(string(resbody))
	if res.StatusCode == 200 {
		fmt.Println("Authenticated!")
		fmt.Println("- Id Token: ")
		// fmt.Print(token_data["id_token"])
		return true, nil
	} else if res.StatusCode == 400 {
		// fmt.Println("res.StatusCode: ")
		// fmt.Print(res.StatusCode)
		return true, nil
	}
	// } else if token_data["error"] != "authorization_pending" && token_data["error"] != "slow_down" {
	// 	log.Panic(token_data["error"])
	// 	return true
	// }
	return false, nil
}
