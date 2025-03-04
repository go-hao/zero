package xjwt

type Algorithm string

const (
	HS256 Algorithm = "HS256"
	HS384 Algorithm = "HS384"
	HS512 Algorithm = "HS512"
	RS256 Algorithm = "RS256"
	RS384 Algorithm = "RS384"
	RS512 Algorithm = "RS512"
	ES256 Algorithm = "ES256"
	ES384 Algorithm = "ES384"
	ES512 Algorithm = "ES512"
)

const (
	ClaimsKeyAudience  string = "aud"
	ClaimsKeyExpiresAt string = "exp"
	ClaimsKeyId        string = "jti"
	ClaimsKeyIssuedAt  string = "iat"
	ClaimsKeyIssuer    string = "iss"
	ClaimsKeyNotBefore string = "nbf"
	ClaimsKeySubject   string = "sub"
)
