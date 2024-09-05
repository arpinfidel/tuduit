package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims

	Email  string `json:"email"`
	UserID int    `json:"user_id"`
}

const (
	priv = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn\nSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i\ncqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC\nPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR\nap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA\nRdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3\nn6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy\nMaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9\nPOIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE\nKdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM\nIvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn\nFcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY\nmEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj\nFuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U\nI5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs\n2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn\n/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT\nOvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86\nEunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+\nhR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0\n4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb\nmDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry\neBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3\nCKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+\n9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq\n-----END RSA PRIVATE KEY-----\n"
	pub  = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41\nfGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7\nmCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp\nHssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2\nXrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b\nODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy\n7wIDAQAB\n-----END PUBLIC KEY-----\n"
)

func TestJWT_Sign(t *testing.T) {
	now := time.Date(2300, 2, 15, 0, 0, 0, 0, time.UTC)
	type fields struct {
		SigningMethod string
		privateKey    string
		publicKey     string
	}
	type args struct {
		claims jwt.Claims
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sign",
			fields: fields{
				SigningMethod: "RS256",
				privateKey:    priv,
				publicKey:     pub,
			},
			args: args{
				claims: &CustomClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject:   "123",
						Issuer:    "123",
						ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
					},
					Email:  "123",
					UserID: 123,
				},
			},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxMjMiLCJzdWIiOiIxMjMiLCJleHAiOjEwNDE3NjgzNjAwLCJlbWFpbCI6IjEyMyIsInVzZXJfaWQiOjEyM30.b-fO_QEwkm_orlbhbtcyyTcfo0M4s_WN2XAzuTb3tv89SPRQ6IxjIyhpnYKOKaIOMkgjlPE93zz9e63I5IVf7AEwC25eonFNGERER9nwRJQEJiy5g7rN9muHFAxFfHRIz_hjCiq0KuN-S0NBXFtwz3gt1neJs_ugwDGJSJ4Mz78P5N8AMST8a4rl8w6nEkpqERs1hxXayJA6h_vzPeIAZiElEPxfnq1ES3SrF3k4TWVddy1XjB_ZFpFAsjMTwbL4qPXwTPgMLPakV1k87McZqkLPEJbyks2m7MOvQuqTzV_BtJ5I-54tNYDnCV_xGWXePeHvw-4_sw13ktl-OiDFgA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := New(tt.fields.SigningMethod, []byte(tt.fields.privateKey), []byte(tt.fields.publicKey))
			if err != nil {
				panic(err)
			}
			got, err := j.Sign(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWT.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JWT.Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: deepequal is not working
// func TestJWT_Verify(t *testing.T) {
// 	now := time.Date(2300, 2, 15, 0, 0, 0, 0, time.UTC)
// 	time.Local = time.UTC

// 	type fields struct {
// 		SigningMethod string
// 		privateKey    string
// 		publicKey     string
// 	}
// 	type args struct {
// 		token        string
// 		targetClaims jwt.Claims
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *CustomClaims
// 		wantErr bool
// 	}{
// 		{
// 			name: "ok",
// 			fields: fields{
// 				SigningMethod: "RS256",
// 				privateKey:    priv,
// 				publicKey:     pub,
// 			},
// 			args: args{
// 				token:        "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxMjMiLCJzdWIiOiIxMjMiLCJleHAiOjEwNDE3NjgzNjAwLCJlbWFpbCI6IjEyMyIsInVzZXJfaWQiOjEyM30.b-fO_QEwkm_orlbhbtcyyTcfo0M4s_WN2XAzuTb3tv89SPRQ6IxjIyhpnYKOKaIOMkgjlPE93zz9e63I5IVf7AEwC25eonFNGERER9nwRJQEJiy5g7rN9muHFAxFfHRIz_hjCiq0KuN-S0NBXFtwz3gt1neJs_ugwDGJSJ4Mz78P5N8AMST8a4rl8w6nEkpqERs1hxXayJA6h_vzPeIAZiElEPxfnq1ES3SrF3k4TWVddy1XjB_ZFpFAsjMTwbL4qPXwTPgMLPakV1k87McZqkLPEJbyks2m7MOvQuqTzV_BtJ5I-54tNYDnCV_xGWXePeHvw-4_sw13ktl-OiDFgA",
// 				targetClaims: &CustomClaims{},
// 			},
// 			want: &CustomClaims{
// 				RegisteredClaims: jwt.RegisteredClaims{
// 					Subject:   "123",
// 					Issuer:    "123",
// 					ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
// 				},
// 				Email:  "123",
// 				UserID: 123,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			j, err := New(tt.fields.SigningMethod, []byte(tt.fields.privateKey), []byte(tt.fields.publicKey))
// 			if err != nil {
// 				panic(err)
// 			}
// 			_, err = j.Verify(tt.args.token, tt.args.targetClaims)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("JWT.Verify() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(tt.args.targetClaims, tt.want) {
// 				t.Errorf("JWT.Verify() = %v, want %v", tt.args.targetClaims, tt.want)
// 			}
// 		})
// 	}
// }
