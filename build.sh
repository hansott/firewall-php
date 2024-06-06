rm -rf build
mkdir build
cd src/agent
go build -o ../../build/aikido
cd ../extension
phpize
cd ../lib
go build -buildmode=c-archive -o ../../build/libaikido_go.a
cd ../../build
CXX=g++ CXXFLAGS="-fPIC -std=c++20 -g -O0 -I../include" LDFLAGS="-L./ -laikido_go" ../src/extension/configure
make
make install
cd ..
