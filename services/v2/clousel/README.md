
```ps1
go mod init closuel
go mod tidy
go work init .
go work use .
```


## CREATE TABLES 
### SQLITE
```ps1
sqlite3.exe .\clousel.db ".read .\scripts\sql\create_clousel_tables.sql"
```


```
POST http://localhost:4321/path/to
[FormParams]
Param1: {{p1}}
Param2: {{p2}}
Param3: {{p3}}
HTTP 200
--------------------

hurl.exe --variable p1=<specify> --variable p2=<specify> --variable p3=<specify> scripts\hurl\path-to-req.hurl
```

## HTTP REUESTS


### REGISTER A COMPANY
```ps1
hurl.exe --variable cname=default --variable email=default@mail.org --variable password=1234 scripts\hurl\company\business_register.hurl
```
### LOGIN AS A COMPANY
```ps1
hurl.exe --variable cname=default --variable password=1234 scripts\hurl\company\business_login.hurl

{"Tocken":"g0vbB9dTTjCAwJxE3OrV5w=="}
$env:BTOCKEN="L3dPLURQRT2cdbkRo/BvJg=="

```


### ADD SECRET KEYS AS A COMPANY
```ps1
hurl.exe --variable tocken=$(echo $env:BTOCKEN)  --variable skey=rk_test_51PajXoRubpSlGSkxjkTLxXvVYhBpBDGEYjfwAwLDEY71FwqbrjX1Umnx6q0kcimmiu9q41J5O2J64Bhg0OyW6neC00uCUisHI7 --variable prodid=prod_RFg9LDezJP3edQ scripts\hurl\company\business_keys.hurl
```

### ADD MACHINE
```ps1
hurl.exe --variable tocken=$(echo $env:BTOCKEN)  --variable mid=550e8400-e29b-41d4-a716-446655440000  --variable cost=1  scripts\hurl\company\business_machine_add.hurl
```

### GET MACHINE(S)

```ps1
# get all
hurl.exe --variable tocken=$(echo $env:BTOCKEN) scripts\hurl\company\business_machine_get_all.hurl

# get by id
hurl.exe --variable tocken=$(echo $env:BTOCKEN)  --variable mid=550e8400-e29b-41d4-a716-446655440000   scripts\hurl\company\business_machine_get_mid.hurl

# get by status
hurl.exe --variable tocken=$(echo $env:BTOCKEN)  --variable status=new   scripts\hurl\company\business_machine_get_status.hurl
```

### UPDATE MACHINE
```ps1
hurl.exe --variable tocken=$(echo $env:BTOCKEN)  --variable mid=550e8400-e29b-41d4-a716-446655440000  --variable cost=N  scripts\hurl\company\business_machine_update.hurl
```



### REGISTER AN USER
```ps1
hurl.exe --variable uname=redrabbit --variable company=default --variable email=redrabbit@mail.com --variable password=1qasw2  .\scripts\hurl\client\client_register.hurl
```

### LOGIN AS AN USER
```ps1
hurl.exe --variable uname=redrabbit --variable password=1qasw2  .\scripts\hurl\client\client_login.hurl

$env:UTOCKEN="jdYiXSo4RqS2TfypgvljMA=="
```

### READ USER BALANCE
```ps1
hurl.exe --variable tocken=$(echo $env:UTOCKEN)  scripts\hurl\client\client_balance.hurl
```

### READ PRICES
```ps1
hurl.exe --variable tocken=$(echo $env:UTOCKEN)  scripts\hurl\client\client_price.hurl
```
### BUY TICKETS
```ps1
hurl.exe --variable tocken=$(echo $env:UTOCKEN)  --variable home="https://www.i.ua" --variable prid=price_1QNAnQRubpSlGSkxgO3f7pSg .\scripts\hurl\client\client_buy.hur
```


### USER & MACHINE
```ps1
# Play
hurl.exe --variable tocken=$(echo $env:UTOCKEN) --variable mid=550e8400-e29b-41d4-a716-446655440000  .\scripts\hurl\client\client_machine_play.hurl
{"EventId":"956fc277-2fce-40da-b4b2-c59b7bb67a53"}

# Poll
hurl.exe --variable tocken=$(echo $env:UTOCKEN) --variable eid=956fc277-2fce-40da-b4b2-c59b7bb67a53  .\scripts\hurl\client\client_machine_poll.hurl

```


## STRIPE WEB HOOK
```ps1
stripe.exe listen --forward-to localhost:4321/webhook/dev
```



# NGINX CONFIGURATION
## HTTPS SUPPORT
Certbot `https://certbot.eff.org/instructions?ws=nginx&os=snap`
Easy to run with nginx
`sudo snap install --classic certbot`
configure NGINX firectly
`sudo certbot --nginx`
or just get certs
`sudo certbot certonly --nginx`

## REDIRECT HTTPS->HTTP
### Nginx reverse proxy redirection configuration:

```conf
  location / {
    proxy_pass http://localhost:4321;
    proxy_set_header Host $host;
  }
```

### CERTBOT
Certbot `https://certbot.eff.org/instructions?ws=nginx&os=snap`  
Easy to run with nginx  `sudo snap install --classic certbot`  
configure NGINX firectly `sudo certbot --nginx` or just get certs `sudo certbot certonly --nginx`

