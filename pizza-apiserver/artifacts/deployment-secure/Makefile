SHELL = bash

OUTPUT := serving-cert-secret.yaml v1alpha1-apiservice.yaml v1beta1-apiservice.yaml

all: tls.key tls.crt $(OUTPUT)

tls.key tls.crt:
		# https://stackoverflow.com/questions/10175812/how-to-generate-a-self-signed-ssl-certificate-using-openssl/41366949#41366949
		/opt/homebrew/opt/openssl@3.1/bin/openssl req -new -x509 \
			-nodes \
			-newkey rsa:4096 \
			-days 365 \
			-subj "/CN=api.pizza-apiserver.svc" \
			-addext 'subjectAltName=DNS:api.pizza-apiserver.svc' \
			-keyout tls.key \
			-out tls.crt

$(OUTPUT): tls.key tls.crt
		sed 's,CERT,$(shell base64 tls.crt),;s,KEY,$(shell base64 tls.key),' $@.template > $@

clean:
		rm -f tls.key tls.crt $(OUTPUT)
