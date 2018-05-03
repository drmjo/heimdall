### Heimdall

Heimdall is a reverse proxy that will take an Open API 3.0.1 definition and spin up
a reverse proxy to a backend service using https://golang.org/pkg/net/http/httputil/#ReverseProxy

This POC aims to prove the concept of abstracted authentication and authorization
layer from the actual service provider, it acts as a Man in the Middle sniffing
the request and the response to control access to a service and also make sure
the service provider is compliant with the predefined json interface. (similar to AWS API Gateway)

It sets up a router from the json and proxy passes the requests to a backend
server if will also validate the response back from the server against
the API schema.

Also has support for `$refs` in the json schema it will combine all refs to a
single large minified json.

It does the sniffing part right now and logs the request/response which is the
main point, there is not a lot of validation going on, but since we have access
to the request we can do anything we please when a request comes it and when a
response is being sent out.

### docker-compose

compose up if you just wanna see it in action

```bash
docker-compose -p h up --build
```
then go to `http://localhost:8844/status` in your browser

you should see the logs from `heimdall_1` and `dummy_service_1` in the terminal where you ran the docker-compose command

### Run for Dev

the below command will grant you a shell with root access
inside of a golang 1.8 container with the volumes mounted

```bash
# get a shell
bash cli.sh

# once inside to download the deps
go get -v github.com/drmjo/heimdall

# install the bin
go install -v github.com/drmjo/heimdall

# run
heimdall --help
```
