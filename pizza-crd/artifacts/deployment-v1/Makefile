SHELL = bash

OUTPUT := pizza-crd.yaml serving-cert-secret.yaml mutatingadmissionregistration.yaml validatingadmissionregistration.yaml

all: tls.key tls.crt $(OUTPUT)

tls.key tls.crt:
		# https://stackoverflow.com/questions/10175812/how-to-generate-a-self-signed-ssl-certificate-using-openssl/41366949#41366949
		/opt/homebrew/opt/openssl@3.1/bin/openssl req -new -x509 \
			-subj "/CN=webhook.pizza-crd.svc" \
			-nodes \
			-newkey rsa:4096 \
			-days 365 \
			-addext 'subjectAltName=DNS:webhook.pizza-crd.svc' \
			-keyout tls.key \
			-out tls.crt

%.yaml: tls.key tls.crt %.yaml.template
		sed 's,CERT,$(shell base64 tls.crt),;s,KEY,$(shell base64 tls.key),' $@.template > $@

clean:
		rm -f tls.key tls.crt $(OUTPUT)
