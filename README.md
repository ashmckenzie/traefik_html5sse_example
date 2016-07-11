Traefik with an HTML5 SSE golang app
====================================

Usage
-----

Clone the repo:

```
git clone https://github.com/ashmckenzie/traefik_html5sse_example.git
```

Run!

```
cd traefik_html5sse_example

# customise TRAEFIK_HTTP_PORT or TRAEFIK_DASHBOARD_HTTP_PORT env variables

./run.sh
```

Once `main.go` is built, Traefik / forego is downloaded and the Docker container built, the following URL's will be available:

* Traefik Dashboard - [http://localhost:8000/dashboard/](http://localhost:8000/dashboard/)
* Traefik HTTP - [http://localhost:8080](http://localhost:8080)
* HTML5 SSE application - [http://localhost:8080/html5sse](http://localhost:8080/html5sse)

The problem
-----------

The HTML5 SSE application is a stripped down version of https://github.com/kljensen/golang-html5-sse-example to contain the bare minimum to illustrate the problem.  The app listens on port `3000` by default, and the main (and only) URL to access is `/` (only accessible within the Docker container).

The HTML5 SSE application blocks requests until there is something to send, which it will do every 5 secs, e.g.:

```
data: Message: 107 - the time is 2016-07-11 07:29:23.205371368 +0000 UTC

data: Message: 108 - the time is 2016-07-11 07:29:28.205811451 +0000 UTC
```

However, when accessing the HTML5 SSE application through Traefik ([http://localhost:8080/html5sse](http://localhost:8080/html5sse)), data is not received every 5 secs, instead the connection blocks and every 3 mins or so the backlog of events is received in one big payload.
