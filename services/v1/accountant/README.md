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


## Workflow
### Register an owner
```sh
hurl.exe  --variable email=wer@asd.com --variable password=1234 owner_register.hurl
# 200 OK 
```

### Login
```sh
hurl.exe  --variable email=wer@asd.com --variable password=1234 owner_login.hurl
# 200 OK 
# {"Token":"233ad02b-0f26-4f56-908c-c6fdf3c46ca8"}
```

### Add carousel
```sh
hurl.exe --variable token=9382be15-0870-4a87-b0ee-d9fb11ee74fd --variable carid=550e8400-e29b-41d4-a716-446655440000 .\carousel_add.hurl
# 200 OK 
```

### Assing security key associated with the owner
```sh
hurl.exe --variable token=9382be15-0870-4a87-b0ee-d9fb11ee74fd --variable skey=sk_test_51PajXoRubpSlGSkxRr6WpEzbhLnZH7fV8ly3yhPNWKsHG7ArdsKQAjVXj6iftvOIiBs5Prp5732t4YbBTQ54v9zI00tAccea11 .\owner_skey.hurl
# 200 OK 
```

### Assing product id to the carousel
```sh
hurl.exe --variable token=9382be15-0870-4a87-b0ee-d9fb11ee74fd --variable carid=550e8400-e29b-41d4-a716-446655440000 --variable prodid=prod_RFg9LDezJP3edQ .\carousel_prodid.hurl
# 200 OK 
```


### Refresh webhook key associated with the owner
```sh
hurl.exe --variable token=9382be15-0870-4a87-b0ee-d9fb11ee74fd .\owner_whook_refresh.hurl
# 200 OK 
```