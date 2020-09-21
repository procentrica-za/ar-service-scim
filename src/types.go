package main

import "github.com/gorilla/mux"

type TokenResponse struct {
	Accesstoken  string `json:"access_token"`
	Refreshtoken string `json:"refresh_token"`
}

type UserResponse struct {
	Message      string `json:"message"`
	Accesstoken  string `json:"access_token"`
	Refreshtoken string `json:"refresh_token"`
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	KeySecret string `json:"keysecret"`
}

type Server struct {
	router *mux.Router
}

type Config struct {
	ISHost          string
	ISPort          string
	APIMHost        string
	APIMPort        string
	ListenServePort string
	ISUsername      string
	ISPassword      string
}
