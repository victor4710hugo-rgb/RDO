module github.com/victor4710hugo-rgb/RDO/backend

go 1.27

// Dependências diretas (importadas pelo código do projeto).
//
// Dependências indiretas (exigidas pelas diretas acima).
require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // geração e validação de tokens JWT
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.59.0 // hash de senhas (bcrypt); direta, por isso não consta em "indirect"
	modernc.org/sqlite v1.63.0 // driver SQLite em Go puro (sem CGO)
)

// Dependências indiretas (exigidas pelas diretas acima).
require (
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect

	golang.org/x/crypto v0.59.0 // indirect


	golang.org/x/sys v0.50.0 // indirect
	modernc.org/libc v1.80.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)
