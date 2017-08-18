locap
=====

Locap stands for Local proxy. It is intended to help development of tools
that use an HTTP API but don't provide cross-origin header.

The actual project is very small an covers nearly all the basic HTTP API
queries but if you need more functionalities, feel free to create an issue.

## Installation

```
$> go get -u -v github.com/tuxlinuxien/locap
```

## Usage

        NAME:
           locap - locap

        USAGE:
           locap [options]

        AUTHOR:
           Yoann Cerda <tuxlinuxien@gmail.com>

        COMMANDS:
             help, h  Shows a list of commands or help for one command

        GLOBAL OPTIONS:
           --port value, -p value         (default: 1314)
           --destination value, -d value  
           --help, -h                     show help



## Use case example

Let's assume that we have a RESTful API on *api.domain.com* that doesn't allow
cross-origin for security reasons but you need to work with it on your
front-end project using AJAX.
Since out front-end project will probably use a local IP, then any HTTP
request sent to *api.domain.com* will be rejected.


**Locap** will be used as a proxy to forward your queries from your local domain
to *api.domain.com*.

By running this command:

```
$> locap --port 1314 --destination 'http://api.domain.com'
[...]
```

you will be able to request **http(s)://api.domain.com** from **http://localhost:1314**

If your production API provides a route that can list all the users via `GET /users`,
just change `http://api.domain.com/users` by `http://localhost:1314/users` then *locap*
will just forward the body and headers of your request to the destination and send you
back the result.
