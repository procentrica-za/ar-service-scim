package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (s *Server) verifycredentials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handle VerifyCredentials in IS with SCIM Has Been Called!")
		//get JSON payload
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)

		//handle for bad JSON provided
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println(err.Error())
			return
		}

		if user.KeySecret == "" {
			invalidUser := UserResponse{}
			invalidUser.Message = "Bad Json provided! No application authorization token provided..."
			js, jserr := json.Marshal(invalidUser)
			if jserr != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, jserr.Error())
				fmt.Println("Error occured when trying to marshal the response to validate user credentials...")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(js)
			return
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		// TODO: Set InsecureSkipVerify as config in environment.env
		client := &http.Client{}
		data := url.Values{}
		data.Set("grant_type", "password")
		data.Add("username", user.Username)
		data.Add("password", user.Password)
		req, err := http.NewRequest("POST", "https://"+config.APIMHost+":"+config.APIMPort+"/token", bytes.NewBufferString(data.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Basic "+"VkVWUVgzWnFqOGFCTmxxVTA1aUZZdjhaWm4wYTpIYmZRamlBa3E3UGRkbk5zZ3JreFVOb1ZwVnNh")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode == 400 {
			invalidUser := UserResponse{}
			invalidUser.Message = "Invalid user credentials for application / Invalid application authorization token."

			js, jserr := json.Marshal(invalidUser)
			if jserr != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, jserr.Error())
				fmt.Println("Error occured when trying to marshal the response to verify user credentials when incorrect credential details for the application were recieved...")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(js)
			return
		}

		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		identityServerResponse := TokenResponse{}
		validUser := UserResponse{}
		err = json.Unmarshal(bodyText, &identityServerResponse)

		if identityServerResponse.Accesstoken == "" {
			validUser.Message = "User found but invalid application authorization token provided..."
		} else {
			validUser.Message = "User credentials successfully validated!"
		}

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding validate credentials response...")
			return
		}

		validUser.Accesstoken = identityServerResponse.Accesstoken
		validUser.Refreshtoken = identityServerResponse.Accesstoken

		js, jserr := json.Marshal(validUser)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the response to validate user credentials...")
			return
		}

		//return back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)

	}
}
