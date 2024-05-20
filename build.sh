rm -rf build
mkdir build
cd build
go build -buildmode=c-archive -o libaikido_go.a ../src/lib/aikido_lib.go
CXX=g++ CXXFLAGS="-fPIC" LIBS="libaikido_go.a" ../src/extension/configure
make
make install