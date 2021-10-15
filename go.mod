module captcha_server

go 1.17

require (
	db v0.0.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.3
	gopkg.in/yaml.v2 v2.4.0
	solver v0.0.0

)

require (
	cookiejar v0.0.0 // indirect
	request v0.0.0 // indirect
)

replace db => ./db

replace solver => ./solver

replace request => ./request

replace cookiejar => ./request/cookiejar
