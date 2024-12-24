ssh mimas@192.168.0.150 "mkdir -p /tmp/accountant"
scp bin/* .env mimas@192.168.0.150:/tmp/accountant
