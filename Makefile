generate_certs:
	cfssl gencert -initca certs/ca-csr.json | cfssljson -bare certs/ca
	cfssl gencert -ca=certs/ca.pem -ca-key=certs/ca-key.pem -config=certs/ca-config.json -profile=mserv certs/server-csr.json | cfssljson -bare certs/server
	cfssl gencert -ca=certs/ca.pem -ca-key=certs/ca-key.pem -config=certs/ca-config.json -profile=mserv certs/client-csr.json | cfssljson -bare certs/client
	openssl req -x509 -nodes -newkey rsa:2048 -keyout certs/client-unkwn-key.pem -out certs/client-unkwn.pem -days 3650 -subj "/C=RU/ST=Saint-Petersburg/L=Saint-Petersburg/O=chapsuk/OU=mserv/CN=*"
