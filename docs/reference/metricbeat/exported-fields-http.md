---
mapped_pages:
  - https://www.elastic.co/guide/en/beats/metricbeat/current/exported-fields-http.html
---

% This file is generated! See scripts/generate_fields_docs.py

# HTTP fields [exported-fields-http]

HTTP module

## http [_http]



## request [_request]

HTTP request information

**`http.request.headers`**
:   The HTTP headers sent

type: object


## response [_response]

HTTP response information

**`http.response.headers`**
:   The HTTP headers received

type: object


**`http.response.code`**
:   The HTTP status code

type: keyword

example: 404


**`http.response.phrase`**
:   The HTTP status phrase

type: keyword

example: Not found


## json [_json]

json metricset

## server [_server]

server

