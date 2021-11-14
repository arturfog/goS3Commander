module gos3commander

go 1.17

require (
	github.com/arturfog/ui v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.41.18
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

require (
	github.com/arturfog/colors v0.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
)

replace github.com/arturfog/colors => ./colors

replace github.com/arturfog/ui => ./ui
