To create typescript for proto files do as follows (https://github.com/thesayyn/protoc-gen-ts)
- npm install -g protoc-gen-ts
- protoc  --ts_out=./ --ts_opt=no_grpc --ts_opt=no_namespace .\cluster\proto\cluster-rpc.proto  