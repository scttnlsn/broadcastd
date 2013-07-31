# broadcastd

Simple HTTP-based pubsub server

## Getting Started

**Install:**

Ensure Go is installed and then run:

    $ go get
    $ make
    $ sudo make install

**Run:**

    $ broadcastd

## API

**Subscribe:**

    $ curl http://localhost:5454

Messages will be newline separated.

**Publish:**

    $ curl -X POST http://localhost:5454 -d "Hello World"

Publish message to connected clients.

## CLI Options

* **-auth=""** - HTTP basic auth password required for all requests
* **-port=5454** - port on which to listen