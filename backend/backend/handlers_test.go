package backend

import (
	"bytes"
	"context"
	"dataconnector/command"
	"dataconnector/filter"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func mockEncryptor(plaintext []byte) (ciphertext []byte, err error) {
	ciphertext = []byte("encrypted commands")
	return
}
func mockDecryptor(ciphertext []byte) (plaintext []byte, err error) {
	plaintext = []byte(`"unencrypted commands"`)
	return
}

type mockBackender struct{}

func (m *mockBackender) upsertUser(ctx context.Context, email, googleKey string) (err error) {
	switch email {
	case "should_fail":
		return fmt.Errorf("random database error")
	}
	return
}
func (m *mockBackender) getCommands(ctx context.Context, googleKey string) (encryptedCommands []byte, err error) {
	encryptedCommands = []byte("encrypted commands")
	return
}
func (m *mockBackender) saveCommands(ctx context.Context, googleKey string, encryptedCommands []byte) (err error) {
	return
}
func (m *mockBackender) Setup() (err error) {
	return
}

type mockCommander struct{}

func (c *mockCommander) Valid() (err error) {
	return
}
func (c *mockCommander) DeParameterize(params []string) (err error) {
	return
}
func (c *mockCommander) Run() (bdy []byte, err error) {
	bdy = []byte(`"mock command was run"`)
	return
}

type mockFilterer struct{}

func (f *mockFilterer) StripUnsafe() (err error) {
	return
}
func (f *mockFilterer) Run(bdy []byte) (out interface{}, err error) {
	out = `"mock command was run"`
	return
}

func TestUpdateGoogleKeyHandler(t *testing.T) {
	tests := []struct {
		name             string
		invalidSignature bool
		claims           jwt.MapClaims
		key              string
		want             *Response
	}{
		{
			name: "valid token should update",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"email":      "a@b.com",
				"google_key": "a temporary google user key",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
			},
		},
		{
			name: "expired token should fail",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(-1 * time.Minute).Unix(),
				"email":      "a@b.com",
				"google_key": "a temporary google user key",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name:             "invalid signature should fail",
			invalidSignature: true,
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"email":      "a@b.com",
				"google_key": "a temporary google user key",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name: "missing email claim should fail",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"not_email":  "a@b.com",
				"google_key": "a temporary google user key",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name: "wrong datatype for email claim should fail",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"email":      1,
				"google_key": "1234567",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name: "missing google_key claim should fail",
			claims: jwt.MapClaims{
				"iat":            time.Now(),
				"exp":            time.Now().Add(1 * time.Minute).Unix(),
				"email":          "a@b.com",
				"not_google_key": "a temporary google user key",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name: "wrong datatype for google_key claim should fail",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"email":      "a@b.com",
				"google_key": 1,
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errInvalidJWT,
			},
		},
		{
			name: "failure to upsert should return an error",
			claims: jwt.MapClaims{
				"iat":        time.Now(),
				"exp":        time.Now().Add(1 * time.Minute).Unix(),
				"email":      "should_fail",
				"google_key": "1234567",
			},
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Error:    errUnableToUpdateGoogleKey,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Backender: &mockBackender{},
				JWTSecret: "myjwtsecret",
			}

			req, err := http.NewRequest(http.MethodPost, "/update_google_key", nil)
			if err != nil {
				t.Fatal(err)
			}

			var signature = cfg.JWTSecret
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			if tt.invalidSignature {
				signature = "thisisnotavalidsignature"
			}

			signed, err := token.SignedString([]byte(signature))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", signed))
			got := cfg.UpdateGoogleKeyHandler(req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRunHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		want        *Response
	}{
		{
			name:        "can run a user's command",
			requestBody: `{"google_key":"123","command_name":"first command","params":["1","2","3"]}`,
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Response: `"mock command was run"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Backender: &mockBackender{},
				Encrypt:   mockEncryptor,
				Decrypt:   mockDecryptor,
			}

			req, err := http.NewRequest(http.MethodPost, "/run", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			got := cfg.RunHandler(req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name      string
		googleKey string
		want      *Response
	}{
		{
			name:      "can get a user's commands",
			googleKey: "123",
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Response: commandFilterSlice{
					{
						Name: "first command",
						Command: &command.Command{
							Type:      "mock command",
							Commander: &mockCommander{},
						},
						Filter: &filter.Filter{
							Type:     "mock filter",
							Filterer: &mockFilterer{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Backender: &mockBackender{},
				Encrypt:   mockEncryptor,
				Decrypt:   mockDecryptor,
			}

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/get?google_key=%s", tt.googleKey), nil)
			if err != nil {
				t.Fatal(err)
			}

			got := cfg.GetHandler(req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		want        *Response
	}{
		{
			name:        "can save a command",
			requestBody: `"command1"`,
			want: &Response{
				status:   http.StatusOK,
				template: "json",
				Response: commandFilterSlice{
					{
						Name: "first command",
						Command: &command.Command{
							Type:      "mock command",
							Commander: &mockCommander{},
						},
						Filter: &filter.Filter{
							Type:     "mock filter",
							Filterer: &mockFilterer{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Backender: &mockBackender{},
				Encrypt:   mockEncryptor,
				Decrypt:   mockDecryptor,
			}

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(tt.requestBody)))
			if err != nil {
				t.Fatal(err)
			}

			got := cfg.SaveHandler(req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func (c *commandFilterSlice) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"unencrypted commands"`:
		*c = commandFilterSlice{
			{
				Name: "first command",
				Command: &command.Command{
					Type:      "mock command",
					Commander: &mockCommander{},
				},
				Filter: &filter.Filter{
					Type:     "mock filter",
					Filterer: &mockFilterer{},
				},
			},
		}

	}

	return nil
}

func (u *userCommands) UnmarshalJSON(b []byte) error {
	u.GoogleKey = "123"

	switch string(b) {
	case `"command1"`:
		u.Commands = []*commandFilter{
			{
				Name: string(b),
				Command: &command.Command{
					Type:      string(b),
					Commander: &mockCommander{},
				},
				Filter: &filter.Filter{
					Type:     string(b),
					Filterer: &mockFilterer{},
				},
			},
		}
	}

	return nil
}
