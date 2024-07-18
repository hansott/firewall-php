export PATH="$PATH:$HOME/go/bin:$HOME/.local/bin"
rm -rf build
mkdir build
cd lib
cd php-extension
phpize
cd ..
protoc --go_out=agent --go-grpc_out=agent ipc.proto
protoc --go_out=request-processor --go-grpc_out=request-processor ipc.proto
cd agent
go get google.golang.org/grpc
go build -gcflags "all=-N -l" -buildmode=c-shared  -o ../../build/aikido_agent.so
cd ../request-processor
go get google.golang.org/grpc
go build -gcflags "all=-N -l" -buildmode=c-shared  -o ../../build/aikido_request_processor.so
cd ../../build
CXX=g++ CXXFLAGS="-fPIC -g -O0 -I../include" LDFLAGS="-lstdc++" ../lib/php-extension/configure
make
cd ..

#sudo cp -f ./build/aikido_agent.so /opt/aikido/aikido_agent.so
#sudo cp -f ./build/aikido_request_processor.so /opt/aikido/aikido_request_processor.so