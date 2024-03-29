package backend

import (
	"context"
	"dataconnector/command"
	"dataconnector/crypto"
	"dataconnector/filter"
)

// Config holds configuration settings for the api
type Config struct {
	Backender
	Encrypt   crypto.Encryptor
	Decrypt   crypto.Decryptor
	JWTSecret string
}

// Backender outlines methods to store and retrieve saved commands
type Backender interface {
	upsertUser(ctx context.Context, email, googleKey string) (err error)
	getCommands(ctx context.Context, googleKey string) (encryptedCommands []byte, err error)
	saveCommands(ctx context.Context, googleKey string, encryptedCommands []byte) (err error)
	Setup() (err error)
}

// Response is the http response to a request
type Response struct {
	status      int
	template    string
	Response    interface{} `json:"response,omitempty"`
	Error       error       `json:"-"`
	ErrorString string      `json:"error,omitempty"`
}

// User is a registered user
type User struct {
	Email string `json:"email"`
}

type userCommand struct {
	GoogleKey   string   `json:"google_key"`
	CommandName string   `json:"command_name"`
	Params      []string `json:"params"`
}

type userCommands struct {
	GoogleKey string           `json:"google_key"`
	Commands  []*commandFilter `json:"commands"`
}

type commandFilterSlice []*commandFilter
type commandFilter struct {
	Name    string           `json:"name"`
	Command *command.Command `json:"command"`
	Filter  *filter.Filter   `json:"filter"`
}
