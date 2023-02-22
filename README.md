## SBZU

SBZU is a high performance stock recorder which record previous price, open price, the highest price, the lowest price, volume, value, and average price of a stock within a day

It's written in go, using redis as database and kafka as message broker.

There is layer for logic under usecase folder, and for accessing database under repo folder.

Redis and kafka is wrapped in its own interface, so we can mock them for unit test.

Below is the general diagram for this app

![alt text](sbzu.jpg)


## Unit Test

Unit tests are included, with help of https://github.com/vektra/mockery, so please install this to generate mock if you want to add or change the mock

for generate mock, you can write go generate ./...


## Communication protocol

Communication protocols used is grpc, with proto file under deliver/grpc/stock.proto

Default server run port is 50051, which is hardcoded in deliver/grpc/server.go

Please follow https://grpc.io/docs/languages/go/quickstart/ to generate protobuff from proto file

command to generate, go under delivery folder and execute:
`protoc --go_out=. --go-grpc_out=. stock.proto`
