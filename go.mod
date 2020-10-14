module emonP1

require (
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	scm.t-m-m.be/emonP1/P1 v0.0.0
)

replace scm.t-m-m.be/emonP1/P1 => ./P1

go 1.15
