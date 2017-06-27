package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/cloudfoundry-community/cfseeker/api"
	"github.com/cloudfoundry-community/cfseeker/commands"
)

// This function is meant to perform any request from the CLI to the API. This
// function is used for the command dispatcher, and therefore it has to return a
// function which is of the type commandFn. Different CLI commands can be
// implemented by making a function which takes in an input struct and returns
// the required HTTP method, URI to give to the HTTP request, and output struct
// to unmarshal the JSON response into (if the HTTP request is successful and
// returns a 2xx code). If the HTTP request fails to send, the error will be
// returned. If a non-2xx code is returned from the API, then an error
// containing the meta.error key in the JSON response will be returned.
func cliRequest(cmdInfo func(interface{}) (string, string, interface{})) commandFn {
	return func(input interface{}) (interface{}, error) {
		method, uri, output := cmdInfo(input)
		if output == nil {
			panic("cmdInfo gave back nil output interface")
		}
		req, err := http.NewRequest(method, uri, bytes.NewReader(nil))
		if err != nil {
			panic(fmt.Sprintf("Couldn't create HTTP request from cmdInfo function: %s", err))
		}
		if usernameFlag != nil && *usernameFlag != "" {
			if passwordFlag == nil || *passwordFlag == "" {
				password := promptForPassword()
				passwordFlag = &password
			}
			req.SetBasicAuth(*usernameFlag, *passwordFlag)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error sending HTTP request: %s", err)
		}

		if basicAuthRequested(resp) { //Do it again with some auth
			username, password := promptForBasicAuth()
			req.SetBasicAuth(username, password)
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("Error sending HTTP request")
			}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading response body: %s", err)
		}

		//Non-200 error codes don't make HTTP library errors. Check ourselves.
		var non2xxCode bool
		if resp.StatusCode/100 != 2 {
			non2xxCode = true
		}

		//Make the HTTP response into a struct we can use
		apiResponse := api.Response{Contents: output}
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal JSON response from server: %s", err)
		}

		if non2xxCode { //Tell the user why their request 400'd or 500'd or whatever
			return nil, fmt.Errorf("Error given from API Request: %s", apiResponse.Meta.Error)
		}
		//Otherwise, give them back what they asked for
		return apiResponse.Contents, nil
	}
}

func basicAuthRequested(r *http.Response) bool {
	return r.StatusCode == 401 && r.Header.Get("WWW-Authenticate") != ""
}

func promptForBasicAuth() (username, password string) {
	fmt.Fprintln(os.Stderr, "Please log in with Basic Auth")
	username = promptForUsername()
	password = promptForPassword()

	return
}

func promptForUsername() (username string) {
	fmt.Fprintf(os.Stderr, "Username: ")
	_, err := fmt.Fscanln(os.Stdin, &username)
	if err != nil {
		panic("Error while reading in username")
	}
	return
}

func promptForPassword() (password string) {
	fmt.Fprintf(os.Stderr, "Password: ")
	bPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic("Error while reading in password")
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "")
	password = string(bPassword)
	return
}

func getCLIFn(command string) (toRun commandFn, toInput interface{}) {
	switch command {
	case "find":
		toRun = cliRequest(findCLICommand)
		toInput = commands.FindInput{
			AppGUID:   *appGUIDFind,
			OrgName:   *orgFind,
			SpaceName: *spaceFind,
			AppName:   *appNameFind,
		}
	case "server":
		bailWith("Refusing to run server mode with --target (-t) flag set")
	}
	return
}

func findCLICommand(input interface{}) (method, uri string, output interface{}) {
	in := input.(commands.FindInput)

	//Form the request uri
	(*targetFlag).Path = api.FindEndpoint
	query := (*targetFlag).Query()
	query.Set(api.FindAppGUIDKey, in.AppGUID)
	query.Set(api.FindOrgNameKey, in.OrgName)
	query.Set(api.FindSpaceNameKey, in.SpaceName)
	query.Set(api.FindAppNameKey, in.AppName)
	(*targetFlag).RawQuery = query.Encode()

	return "GET", (*targetFlag).String(), commands.FindOutput{}
}
