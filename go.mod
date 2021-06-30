module b00m.in/finisher

require (
	github.com/dhconnelly/rtreego v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/lib/pq v1.0.0
	github.com/paulmach/go.geojson v1.4.0
	github.com/plutov/paypal/v4 v4.0.0
	github.com/smira/go-point-clustering v1.0.1
	github.com/stripe/stripe-go/v72 v72.53.0
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	rsc.io/quote v1.5.2
)

replace b00m.in/finisher/data => ./data

replace b00m.in/finisher/comms => ./comms

replace b00m.in/finisher/subs => ./subs

//replace github.com/lib/pq => ../mygo/src/github.com/lib/pq

go 1.13
