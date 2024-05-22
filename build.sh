rm -rf build
mkdir build
cd src/extension
phpize
cd ../../build
go build -buildmode=c-archive -o libaikido_go.a ../src/lib/aikido_lib.go
CXX=g++ CXXFLAGS="-fPIC -I../include" LDFLAGS="-L./ -laikido_go" ../src/extension/configure
make
make install