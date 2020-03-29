module shitty.moe/satelit-project/kumos

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	google.golang.org/grpc v1.28.0
	shitty.moe/satelit-project/satelit-scraper v0.0.0
)

replace shitty.moe/satelit-project/satelit-scraper v0.0.0 => ../..
