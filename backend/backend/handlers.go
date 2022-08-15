package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/brentadamson/dataconnector/backend/crypto"
	"github.com/brentadamson/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	errUnauthorized             = fmt.Errorf("unauthorized")
	errInvalidGoogleKey         = fmt.Errorf("invalid Google key")
	errInvalidEmail             = fmt.Errorf("invalid email address")
	errInvalidJWT               = fmt.Errorf("invalid JSON Web Token")
	errUnableToRegisterUser     = fmt.Errorf("unable to register user")
	errUnableToUpdateGoogleKey  = fmt.Errorf("unable to update Google key")
	errInvalidCommand           = fmt.Errorf("invalid command")
	errDuplicateCommandName     = fmt.Errorf("duplicate command name")
	errUnableToRetrieveCommands = fmt.Errorf("unable to retrieve commands")
	errUnableToSaveCommands     = fmt.Errorf("unable to save commands")
)

func (cfg *Config) RegisterUserHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	user := struct {
		Email string `json:"email"`
		Key   string `json:"key"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Info.Println(err)
		rsp.Status = http.StatusInternalServerError
		return rsp
	}

	if user.Key != cfg.Key {
		rsp.Error = errUnauthorized
		return rsp
	}

	// not even sure if this is needed but I guess doesn't hurt
	user.Email, err = url.QueryUnescape(user.Email)
	if user.Email == "" || err != nil {
		rsp.Error = errInvalidEmail
		return rsp
	}

	if err := cfg.Backender.registerUser(context.Background(), user.Email); err != nil {
		rsp.Error = errUnableToRegisterUser
		return rsp
	}

	return
}

// UpdateGoogleKeyHandler updates a user's Google Sheets API Key
func (cfg *Config) UpdateGoogleKeyHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	// decode the JWT and make sure the signature is valid
	token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// a workaround for https://github.com/dgrijalva/jwt-go/issues/314
		mapClaims := token.Claims.(jwt.MapClaims)
		delete(mapClaims, "iat")
		return []byte(cfg.JWTSecret), nil
	})
	if token == nil || err != nil { // this will also validate the signature and timestamp
		rsp.Error = errInvalidJWT
		return rsp
	}

	claims := token.Claims.(jwt.MapClaims)
	var email string
	val, ok := claims["email"]
	if !ok {
		rsp.Error = errInvalidJWT
		return rsp
	}
	switch val.(type) {
	case string:
		email = val.(string)
	default:
		rsp.Error = errInvalidJWT
		return rsp
	}

	var googleKey string
	val, ok = claims["google_key"]
	if !ok {
		rsp.Error = errInvalidJWT
		return rsp
	}
	switch val.(type) {
	case string:
		googleKey, err = url.QueryUnescape(val.(string))
		if googleKey == "" || err != nil {
			rsp.Error = errInvalidGoogleKey
			return rsp
		}
	default:
		rsp.Error = errInvalidJWT
		return rsp
	}

	if err := cfg.Backender.upsertUser(context.Background(), email, googleKey); err != nil {
		rsp.Error = errUnableToUpdateGoogleKey
		return rsp
	}

	return
}

// GetHandler returns a user's commands
// curl 'http://127.0.0.1:8000/get?google_key=my_key'
func (cfg *Config) GetHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	googleKey := r.FormValue("google_key")
	if googleKey == "" {
		rsp.Error = errInvalidGoogleKey
		return rsp
	}

	rsp.Response, rsp.Error = cfg.getCommands(googleKey)
	return
}

func (cfg *Config) getCommands(googleKey string) (cmds commandFilterSlice, err error) {
	// get all of their commands
	encryptedCommands, err := cfg.Backender.getCommands(context.Background(), googleKey)
	if err != nil {
		err = errUnableToRetrieveCommands
		return
	}

	if len(encryptedCommands) > 0 {
		var decryptedCommands []byte
		decryptedCommands, err = cfg.Decrypt(encryptedCommands)
		if err != nil {
			return
		}

		if err = json.Unmarshal(decryptedCommands, &cmds); err != nil {
			err = errUnableToRetrieveCommands
			return
		}
	}

	return
}

// RunSheetsHandler runs a single command (for Google Sheets)
// curl -XPOST 'http://127.0.0.1:8000/sheets_run' -d '{"command_number":2, "params":["123"], "email":"me@example.com", "command":{"name":"my command","filter":{"type":"jmespath","filter":{"expression":""}},"command":{"type":"direct","command":{"method":"get","url":"https://example.com"}}},"key":"1234567"}'
func (cfg *Config) RunSheetsHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	uc := struct {
		Email         string        `json:"email"`
		CommandFilter commandFilter `json:"command"`
		CommandNumber int           `json:"command_number"`
		Params        []string      `json:"params"`
		Key           string        `json:"key"`
	}{}

	rsp.Error = json.NewDecoder(r.Body).Decode(&uc)
	if rsp.Error != nil {
		log.Info.Println(rsp.Error)
		rsp.Status = http.StatusInternalServerError
		return rsp
	}

	if uc.Key != cfg.Key {
		rsp.Error = errUnauthorized
		return rsp
	}

	// not even sure if this is needed but I guess doesn't hurt
	var err error
	uc.Email, err = url.QueryUnescape(uc.Email)
	if uc.Email == "" || err != nil {
		rsp.Error = errInvalidEmail
		return rsp
	}

	// Optional check user status
	if cfg.UserFn != nil {
		if rsp.Error = cfg.UserFn(uc.Email, uc.CommandNumber); rsp.Error != nil {
			return rsp
		}
	}

	if rsp.Error = uc.CommandFilter.Command.Valid(); rsp.Error != nil {
		return rsp
	}

	if rsp.Error = uc.CommandFilter.Command.DeParameterize(uc.Params); rsp.Error != nil {
		return rsp
	}

	var bdy []byte
	bdy, rsp.Error = uc.CommandFilter.Command.Run()
	if rsp.Error != nil {
		rsp.Error = fmt.Errorf("%s:%q", rsp.Error, string(bdy))
		return rsp
	}

	// filter the data
	if rsp.Error = uc.CommandFilter.Filter.StripUnsafe(); rsp.Error != nil {
		return rsp
	}

	var out interface{}
	out, rsp.Error = uc.CommandFilter.Filter.Run(bdy)
	if rsp.Error != nil {
		rsp.Error = fmt.Errorf("%s:%q", rsp.Error, string(bdy))
		return rsp
	}

	rsp.Response = out
	return rsp
}

// RunHandler runs a single, named command
// curl -XPOST 'http://127.0.0.1:8000/run' -d '{"google_key":"my_key", "command_name":"curl_command", "params":["1"]}'
func (cfg *Config) RunHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	var bdy []byte
	bdy, rsp.Error = ioutil.ReadAll(r.Body)
	if rsp.Error != nil {
		rsp.Status = http.StatusInternalServerError
		return rsp
	}
	defer r.Body.Close()

	var user *userCommand
	if rsp.Error = json.Unmarshal(bdy, &user); rsp.Error != nil {
		return rsp
	}

	var err error
	user.GoogleKey, err = url.QueryUnescape(user.GoogleKey)
	if user.GoogleKey == "" || err != nil {
		rsp.Error = errInvalidGoogleKey
		return rsp
	}

	// get their selected command
	commands, err := cfg.getCommands(user.GoogleKey)
	if errors.Is(err, crypto.ErrMalformedCiphertext) {
		// in this case, decrypted commands is likely "" (e.g. the user has no saved commands)...any other scenarios?
		rsp.Error = nil
		rsp.Response = "No saved commands ;(. Please create a command: Add-ons -> Data Connector -> Manage connections"
		return rsp
	}

	for _, cmd := range commands {
		if cmd.Name != strings.ToLower(user.CommandName) {
			continue
		}

		if rsp.Error = cmd.Command.Valid(); rsp.Error != nil {
			return rsp
		}

		if rsp.Error = cmd.Command.DeParameterize(user.Params); rsp.Error != nil {
			return rsp
		}

		cmd.Command.AddCredentials(user.Credentials)

		bdy, rsp.Error = cmd.Command.Run()
		if rsp.Error != nil {
			rsp.Error = fmt.Errorf("%s:%q", rsp.Error, string(bdy))
			return rsp
		}

		// filter the data
		if rsp.Error = cmd.Filter.StripUnsafe(); rsp.Error != nil {
			return rsp
		}

		var out interface{}
		out, rsp.Error = cmd.Filter.Run(bdy)
		if rsp.Error != nil {
			rsp.Error = fmt.Errorf("%s:%q", rsp.Error, string(bdy))
			return rsp
		}

		rsp.Response = out
		return
	}
	return rsp
}

// SaveHandler saves a user's commands
// curl -XPOST 'http://127.0.0.1:8000/save' -d '{"google_key":"my_key","commands":[{"name":"api_command","command":{"type":"direct","command":{"method":"get","url":"https://api.chucknorris.io/jokes/random", "headers":[{"key":"User-Agent","value":"Data Connector"}]}},"filter":{"type":"jmespath","filter":{"expression":"[[value]]"}}}]}'
func (cfg *Config) SaveHandler(r *http.Request) (rsp *Response) {
	rsp = &Response{
		Status:   http.StatusOK,
		Template: "json",
	}

	var bdy []byte
	bdy, rsp.Error = ioutil.ReadAll(r.Body)
	if rsp.Error != nil {
		rsp.Status = http.StatusInternalServerError
		return rsp
	}
	defer r.Body.Close()

	// ensure the JSON is correct and the commands are valid
	var user *userCommands
	if rsp.Error = json.Unmarshal(bdy, &user); rsp.Error != nil {
		return rsp
	}

	var err error
	user.GoogleKey, err = url.QueryUnescape(user.GoogleKey)
	if user.GoogleKey == "" || err != nil {
		rsp.Error = errInvalidGoogleKey
		return rsp
	}

	uniqueNames := make(map[string]struct{})
	for _, cmd := range user.Commands {
		if cmd.Name == "" {
			rsp.Error = errInvalidCommand
			return rsp
		}

		cmd.Name = strings.ToLower(cmd.Name)
		if _, ok := uniqueNames[cmd.Name]; ok {
			rsp.Error = errDuplicateCommandName
			return rsp
		}

		uniqueNames[cmd.Name] = struct{}{}
		if rsp.Error = cmd.Command.Valid(); rsp.Error != nil {
			return rsp
		}
	}

	j, err := json.Marshal(user.Commands)
	if err != nil {
		return rsp
	}

	var encrypted []byte
	encrypted, rsp.Error = cfg.Encrypt(j)
	if rsp.Error != nil {
		return rsp
	}

	if err := cfg.Backender.saveCommands(context.Background(), user.GoogleKey, encrypted); err != nil {
		rsp.Error = errUnableToSaveCommands
		return rsp
	}

	// Maybe we don't have to return their commands? Probably not a big deal.
	// We could probably just return the commands passed in but good to double-ensure everything works, I guess.
	rsp.Response, rsp.Error = cfg.getCommands(user.GoogleKey)
	return
}
