module github.com/DimaKoz/practicummetrics/cmd/staticlint

go 1.19

require github.com/DimaKoz/practicummetrics/pkg/mainexit v0.0.0
replace github.com/DimaKoz/practicummetrics/pkg/mainexit => ../../pkg/mainexit

require (
	github.com/fatih/errwrap v1.5.0
	github.com/sashamelentyev/usestdlibvars v1.24.0
	golang.org/x/tools v0.14.0
	honnef.co/go/tools v0.4.6
	github.com/DimaKoz/practicummetrics/pkg/mainexit iter19
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	golang.org/x/exp/typeparams v0.0.0-20221208152030-732eee02a75a // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
