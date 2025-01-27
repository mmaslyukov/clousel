## Cross compile to arm64
```cmd
set GOARCH=arm64
set GOOS=linux
go build -o binary_name
```


## HTTPS setting up

Generate private key (.key)
```
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

- or - 

# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
# https://safecurves.cr.yp.to/
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
```
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```