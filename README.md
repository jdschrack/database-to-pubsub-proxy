# Database Proxy POC
This is a proof of concept for a database proxy that intercepts SQL queries and publishes them to Google PUB/SUB. 

The purpose of this POC is to allow for the usa of sql queries to trigger events outside of the normal application flow. This could allow for things like notifications, triggers to other applications to refresh data, etc...

### Create a Self Signed Certificate
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
```
You will need to generate your own SSL certificates and to enable SSL support.

This is a work in progress and should not be used for production use.

