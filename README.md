[![Build Status](https://travis-ci.com/georlav/objectmap.svg?token=LUHt821atupKxCks2oys&branch=master)](https://travis-ci.com/georlav/objectmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/georlav/objectmap)](https://goreportcard.com/report/github.com/georlav/objectmap)
[![](https://img.shields.io/badge/unicorn-approved-ff69b4.svg)](https://www.youtube.com/watch?v=9auOCbH5Ns4)

# ObjectMap
A Simple command line tool that helps you check PHP and Java applications for insecure deserialization vulnerabilities.

Supported checks
 * PHP Object Injection
 * Java Deserialization
 
## Requirements

 * golang
 
## Basic usage examples

Load a request from a file.
```bash
objectmap -r request.file
```

Request data should be in valid format (HTTP/1.x wire representation)
```http request
POST /form HTTP/1.1
Host: 127.0.0.1:8056
Content-Length: 42
Content-Type: application/x-www-form-urlencoded
User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)
Cookie: PHPSESSID=298zf09hf012fh2; csrftoken=u32t4o3tb3gg43; _gat=1;

license=string&content=string&paramsXML=ss
```

Or you can initialize your target using command line params
```bash
objectmap -u 127.0.0.1:8056/form --body="license=string&content=string&paramsXML=ss" --method=post
```

Application analyzes target, calculates all the available insertion points and injects various payloads to detect
insecure deserialization vulnerabilities.

## Report example

```Results
INFO Calculating insertion points                 
INFO Found 10 insertion points                    
+--------------------+----------------------+------------+
|  INSERTION POINT   |    VULNERABILITY     |   STATUS   |
+--------------------+----------------------+------------+
| Param[paramsXML]   | PHP Object Injection | Clean      |
| Cookie[_gat]       | Java Deserialization | Clean      |
| Cookie[PHPSESSID]  | Java Deserialization | Vulnerable |
| Param[license]     | PHP Object Injection | Clean      |
| Cookie[PHPSESSID]  | PHP Object Injection | Clean      |
| Cookie[csrftoken]  | PHP Object Injection | Clean      |
| Param[license]     | Java Deserialization | Clean      |
| Cookie[csrftoken]  | Java Deserialization | Clean      |
| Param[content]     | PHP Object Injection | Clean      |
| Header[User-Agent] | PHP Object Injection | Clean      |
| Param[paramsXML]   | Java Deserialization | Clean      |
| Header[User-Agent] | Java Deserialization | Clean      |
| Cookie[_gat]       | PHP Object Injection | Clean      |
| Param[content]     | Java Deserialization | Clean      |
+--------------------+----------------------+------------+
|                         TOTAL REQUESTS    |     40     |
+--------------------+----------------------+------------+
```

## Available Options

```
--url value, -u value                    Target url
--url-scheme value, --us value           Set the URL scheme [http, https] (default: "http")
--method value, -m value                 Set the HTTP request method, supported methods are [GET POST PUT PATCH DELETE] (default: "GET")
--body value                             Set the request body
--request value, -r value                Load http request from a file
--request-concurrency value, --rc value  Set the number of concurrent requests (default: 1)
--request-retries value, --rr value      Set number of retries on request failure (default: 2)
--no-follow, --nf                        Do not follow http redirects (default: follows)
--timeout value, -t value                Set the max timeout limit in seconds for http requests (default: 10)
--user-agent value                       Set client user agent (default: "ObjectMap/1.0")
--random-agent                           Set client to use a random user agent
--banner, -b                             Retrieve server banner
--verbose value, -v value                Set the verbosity level [1-5] (default: 4)
--help, -h                               Show help
```

## Installing
```bash
go get github.com/georlav/objectmap/cmd/objectmap
```

## Compiling from sources
```
git clone git@github.com:georlav/objectmap.git
cd objectmap
make build
```

## Running the tests

```bash
cd $GOPATH/src/github.com/georlav/objectmap
make test
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. 

## Authors

* **georlav** - *Initial work*

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
