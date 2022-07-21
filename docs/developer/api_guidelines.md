# Introduction
This document serves as a guideline for adding endpoints to the Keptn API,
and should help to provide a consistent user experience across all of our APIs.

```
Notation Conventions and Compliance
The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", 
"MAY", and "OPTIONAL" in this document are to be interpreted as described in
[BCP 14](https://tools.ietf.org/html/bcp14) [[RFC2119](https://tools.ietf.org/html/rfc2119)] [[RFC8174](https://tools.ietf.org/html/rfc8174)]
when, and only when, they appear in all capitals, as shown here.
``` 

# API Grouping
We define three different *types* of APIs:

* `public` - customer facing public API exposed by the API gateway
* `internal` - internally APIs used by internal services never exposed by the API gateway
* `operations` - internal API providing non-functional features like statistics, debug and testing never exposed by the API gateway

## Public APIs
Public APIs MUST BE provided by the API gateway and mapped to the K8s Service as follows:

API Gateway: `http://cluster.keptn.sh/api/<component>/<version>/<api resources>`

K8S Service: `<component>/public/<version>/<api resources>`

## Internal APIs
Internal APIs MUST NOT BE exposed by the API gateway and SHOULD be accessed by other services in Keptn as follows:

K8S Service:  `<component>/internal/<version>/<api resources>`

## Operations APIs
Operations APIs MUST NOT be exposed by the API gateway and SHOULD be cover non-functional features
like statistics, debug or testing APIs etc...

K8S Service: `<component>/operations/<version>/<api resources>` 

---

# Authentication
Each request MUST have a *Authorization* header containing a base64 encoded token string:

## Example

```
GET /controlplane/v1/sequences
Authorization: token dGhlc2VjcmV0dG9rZW5mb3JrZXB0bmFwaQ==
```

Internal communication MAY NOT require a token.

---

# Custom Methods

Sometimes, functionality cannot be expressed using the default request methods on classical REST resources.
For this use cases *custom methods* SHOULD be used:

```
https://domain/api/<component>/<version>/<api resources>:<custom-method>
```

- Custom methods SHOULD use the `HTTP POST` method.
- The URL MUST end with a verb separated by a colon.

## Example
```
POST cluster.keptn.sh/api/controlplane/v1/sequences/94dbe764-e186-48bf-b90d-a53b21fca847:pause
POST cluster.keptn.sh/api/tokens/v1/token:rotate
```

 ---
# Filtering

API endpoints supporting result filtering MUST accept a filter
as a query parameter named *filter* containing a *filter expression*:

In general each field SHOULD be filter-able. Exceptions MUST be documented.
*Filter expressions* follow the following form: `<fieldname> operator <value>`

## Filters

### Filter based on numbers

**Operators:**

```
= != < <= > >=
```

**Format**:

```
int - decimal and hexadecimal (leading 0x)`
float - scientific notation (optional exponent e or E) 
```

**Example:**

```
failedAttempt >= 3
average >= 1.0E2
responseTimeSec > 0.23
```

### Filter based on strings

**Operators:**

```
= != contains, starts-with, ends-with
```

**Format**:

```
strings wrapped in single quotes
special characters escaped with \
```

**Example:**

```
projectName = 'myproject'
```

### Filter based on booleans

**Operators:**

```
= !=
```

**Format**:

```
true false
```

**Example:**

```
`active != true`
```

### Filter based on Date / Time

**Operators:**

```
= != < <= > >=
```

**Format**:

```
ISO 8601 string
```

**Example:**

```
lastModified >= '2020-01-01T01:00:00Z
```

---
# JSON Format

Payload sent to and returned by API endpoints SHOULD be represented as JSON objects and MUST have the
`Accept` and `Content-Type` headers set to `application/json`.
The structure of these objects SHOULD be documented in the swagger API docs.
When accepting JSON payload from a client, the API endpoint SHOULD validate if the payload adheres to the documented structure, and shall return an error
if it does not.
However, if a client sends an object with an unknown property, the API endpoint SHOULD just ignore it, and proceed with the execution of the request.

Groups of properties that are mutually exclusive from each other SHOULD be stored in nested objects, instead of having all properties in a flat JSON object.

*Example:*

*NOT OK*:
```json
{
  "name": "my-credential-config",
  "remoteURL": "my-url:8080", 
  "httpsToken": "some-value",
  "httpsUser": "my-user",
  "sshKey": "another-value",
  "sshKeyPassword": "something"
}
```

*OK:*
```json
{
  "name": "my-credential-config",
  "remoteURL": "my-url:8080", 
  "https": {
    "token": "some-value",
    "user": "my-user"
  },
  "ssh": {
    "key": "another-value",
    "keyPassword": "something"
  }
}
```
---
# Asynchronous Operations

For requests that trigger a long running/asynchronous operation, the API endpoint SHOULD respond with a `202 - Accepted` HTTP response code.
The returned payload SHOULD include an identifier for the job triggered by the request. This identifier should enable clients to track the current state of the
job, or to perform other operations on it (e.g. cancelling, pausing). It SHOULD be suffixed with `ID`

*Example:*

```
POST /projects/19340-29304:migrate

Status: 202 - Accepted

Body:

{
  "jobId": "123-456-789"
}
```

After accepting and starting/queueing a long running operation, the server MAY provide the current state of the operation via a `GET` endpoint. Example:

```
GET /projects/123-456-789:migrate?job-id=123-456-789

200 OK
{
  "status": {...}
}
```
---
# Naming Conventions

## Component Names

Service names SHOULD be concise and as simple as possible. They SHOULD be written in *kebap-case*. Also, since the name of a service will also be used
for the internal domain name in K8s, they SHOULD adhere to the DNS label standard defined in [RFC 1035](https://datatracker.ietf.org/doc/html/rfc3339), which means that service names SHOULD:

- Contain at most 63 characters
- Contain only lowercase alphanumeric characters or '-'
- Start with an alphabetic character
- End with an alphanumeric character

## Path Names

Paths to an API for a specific type of resource MUST contain the resource type (in singular). E.g.:

`v1/project`, and NOT `v1/projects`

## Query parameters

Query parameters SHOULD be in *kebab-case* and SHOULD not contain any upper case letters. For time based properties, please also refer to the [time formats](#time-formats) section.

## Field names

Field names SHOULD be in camel case and start with a lower case character. Further, field names SHOULD only contain alphanumeric characters.
Property names of JSON payloads sent to and from the API SHOULD always be well defined, and arbitrary property names (e.g. keys in a map) are not allowed.
If a map should be represented, the corresponding property SHOULD rather be an array containing objects with properties for the key and value, respectively
Examples:

*NOT OK:*

```json
{
  "receivedEvents": {
    "my-arbitrary-event-type": 4,
    "my-other-event-type": 2
  }
}
```

*OK:*

```json
{
  "receivedEvents": [
    {
      "key": "my-arbitrary-event-type",
      "value": 4
    },
    {
      "key": "my-other-event-type",
      "value": 2
    }
  ]
}
```
---
# Pagination

As soon as a *listable collection* is provided via the API, it MUST
support pagination.

The endpoint
* MUST provide a query parameter called `page-key` that points to the page to fetch.
  If the parameter is missing, the first page MUST be returned
* MUST provide a query parameter called `page-size` that sets the
  requested number of elements for the page.
* MUST provide a documentation mentioning the max page size that is supported
  If missing, there MUST be a reasonable default size.
* MUST contain a field `nextPageKey` in the response body that points to the next page to fetch.
* MAY contain a field `totalCount` that represents the total number of elements in the collection.

## Example
```
GET " cluster.keptn.sh/api/controlplane/v1/log?&page-size=100"
HTTP/1.1 200 OK
Content-Type: application/json
{
  "documents" : [
     …
  ],
  "nextPageKey" : 101,
  "totalCount" : 500
}
```
---
# Standard Methods

In most of the cases it is **recommended** to prefer standard methods that map naturally to HTTP methods.
However, when there is no such natural mapping available, or the API starts to look peculiar when trying to enforce it,
custom methods MAY be used.

We distinguish between the following standard methods:

| Standard Method | HTTP Method                |
|-----------------|----------------------------|
| list            | GET <resource collection>  |
| get             | GET <resource>             |
| create          | POST <resource collection> |
| full update     | PUT <resource>             |
| partial update  | PATCH <resource>           |
| delete          | DELETE <resource>          |

## list
Used to search for resources and get back a collection
* MUST map to `HTTP GET` method
* MUST NOT use a request body
* MUST return a response body that contains a list of resources that MAY be empty
* MAY take additional HTTP query parameters


## get
Used for accessing a single resource
* MUST map to `HTTP GET` method
* MUST NOT use a request body
* MUST return a response body that contains the resource if successful
* MAY take additional HTTP query parameters

## create
Used to create a new resource
* MUST map to `HTTP POST` method
* MAY accept a resource id to choose for the caller. If the resource id already exists, the method MUST fail
* MAY contain the fields necessary to construct the resource as request body
* MUST return `HTTP 201 - Created` on success
* SHOULD return the resource id in the *Location* header as a relative URL
* SHOULD return a response body that contains the created resource

## full update
Used to do an update on a resource
* MUST map to `HTTP PUT` method
* MAY be used to create a new resource
* MUST contain all fields of the resource
* MUST return `HTTP 201 - Created` if it was used for creating a new resource
* SHOULD return the resource id in the *Location* header as a relative URL if it was used for creating of a new resource
* SHOULD return a response body that contains the resource if it was used for creating of a new resource
* MUST return an empty response body and `HTTP 204 - No content` if it was used for updating an existing resource

## partial update
Used do a partial update on a resource
* MUST map to `HTTP PATCH` method
* MUST contain updated fields in JSON Merge Patch format (https://datatracker.ietf.org/doc/html/rfc7396) in the request body
* MUST not be used to create a new resource
* MUST return `HTTP 204` -No Content on success
* MUST return an empty response body

## delete
Used to remove a resource
* MUST map to `HTTP DELETE` method
* MAY be applicable to a single resource or a collection of resources. A filter criteria MUST be provided as query parameters to delete a collection of resources.
* MUST not use a request body
* MUST not use a response body
* MUST return `HTTP 204 - No Content` on success
* MUST return `HTTP 202 - Accepted` on async. deletion request

---
# Time Values Conventions

## Time Formats

Time values SHOULD be represented as UTC times and MUST be encoded according to the *ISO 8601* standard.
This standard accepts multiple variations of time formats, however to avoid confusion all time values in the APIs MUST be stored and represented using the following format:

- `2022-06-29T08:56:38.547Z`

When ingesting time values, i.e. by accepting them as request parameters, APIs SHOULD be able to handle multiple variations of the RFC3339 standard. Examples:

- `1985-04-12T23:20:50.52Z`
- `1996-12-19T16:39:57-08:00`
- `1990-12-31T23:59:60Z`
- `1990-12-31T15:59:60-08:00`
- `1937-01-01T12:00:27.87+00:20`
- `2022-06-28T13:43:14.71657055Z`

## Time Value Property Naming

All properties referring to a time value MUST end with the suffix `Time`, e.g. `modifyTime`, `triggerTime` etc. This suffix shall be used consistently across all APIs, meaning that
we SHOULD not have properties like `creationTimestamp`, `modifiedAt` etc.
For time frames, the property names `startTime` and `endTime` MUST be consistently used - other variations like `fromTime` or `toTime` therefore MUST NOT be used. 

---
# Versioning

Each service has its own individual API version, i.e., there is no shared version for all services within Keptn. The versioning scheme MUST follow
the semantic versioning conventions by maintaining a `major`, `minor` and `patch` version (`<major>.<minor>.<patch>`).

These versions MUST be incremented according to the following rules:

## Major Version Updates
The major version is increased whenever there is a breaking change, such as a modified structure of a request model, or changed HTTP status codes.

## Minor Version Updates
The minor version is increased when a new, but backwards compatible feature is introduced.

## Patch Version Updates
The patch version is increased for changes that do not change the functionality of the API.

## Version Info in API URLs
The exact API version, consisting of the major, minor and patch version shall be documented in the APIs documentation.
The URLs pointing to released API endpoints shall only include the major version, such as:

```
/<service-name>/v<major version>/<api-resource>
```

For unreleased APIs, the version SHOULD be provided in the following format:

```
/<service-name>/v0.<minor version>/<api-resource>
```