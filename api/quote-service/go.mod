module github.com/robertpbellamy/react-quote-calculator/api/quote-service

go 1.16

require (
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/robertpbellamy/react-quote-calculator/api/db v0.0.0-00010101000000-000000000000 // indirect
	go.mongodb.org/mongo-driver v1.5.2
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/robertpbellamy/react-quote-calculator/api/db => ../db
