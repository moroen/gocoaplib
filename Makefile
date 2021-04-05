static_lib = lib/static/libgocoap.a
shared_lib = lib/shared/libgocoap.so

src = main.go

libs: $(static_lib) $(shared_lib)

$(static_lib): $(src)
	go mod tidy
	go build -buildmode=c-archive -o $(static_lib) 

$(shared_lib): $(src)
	go mod tidy
	go build -buildmode=c-shared -o $(shared_lib)

clean:
	rm -rf lib/*