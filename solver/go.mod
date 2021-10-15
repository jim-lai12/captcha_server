module solver

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/onsi/gomega v1.16.0 // indirect
	request v0.0.0
)

replace request => ../request

replace cookiejar => ../request/cookiejar
