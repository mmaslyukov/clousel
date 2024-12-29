bun run build
zip -r versel.zip .\dist\
scp .\versel.zip mimas@192.168.0.150:/tmp