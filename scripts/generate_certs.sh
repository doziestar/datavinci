#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Directory to store the certificates
CERT_DIR="./certs"

# Create the certificate directory if it doesn't exist
mkdir -p $CERT_DIR

# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout $CERT_DIR/ca-key.pem -out $CERT_DIR/ca-cert.pem -subj "/C=US/ST=State/L=City/O=Organization/OU=DataVinci/CN=DataVinciCA"

echo "CA's self-signed certificate generated"

# Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout $CERT_DIR/server-key.pem -out $CERT_DIR/server-req.pem -subj "/C=US/ST=State/L=City/O=Organization/OU=DataVinci/CN=localhost"

# Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in $CERT_DIR/server-req.pem -days 360 -CA $CERT_DIR/ca-cert.pem -CAkey $CERT_DIR/ca-key.pem -CAcreateserial -out $CERT_DIR/server-cert.pem -extfile <(echo -e "subjectAltName=DNS:localhost,IP:127.0.0.1")

echo "Server's signed certificate generated"

# Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout $CERT_DIR/client-key.pem -out $CERT_DIR/client-req.pem -subj "/C=US/ST=State/L=City/O=Organization/OU=DataVinci/CN=DataVinciClient"

# Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in $CERT_DIR/client-req.pem -days 360 -CA $CERT_DIR/ca-cert.pem -CAkey $CERT_DIR/ca-key.pem -CAcreateserial -out $CERT_DIR/client-cert.pem

echo "Client's signed certificate generated"

echo "All certificates have been generated in $CERT_DIR"