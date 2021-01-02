module m0v.in/finisher

require (
	github.com/dhconnelly/rtreego v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/lib/pq v1.0.0
	github.com/paulmach/go.geojson v1.4.0
	github.com/smira/go-point-clustering v1.0.1
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	rsc.io/quote v1.5.2
)

replace m0v.in/finisher/data => ./data

replace m0v.in/finisher/comms => ./comms

//replace github.com/lib/pq => ../mygo/src/github.com/lib/pq

go 1.13
