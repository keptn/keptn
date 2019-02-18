# CloudEvent Specification

CloudEvents is a vendor-neutral specification for defining the format of event data.

## Table of Contents
- [Specifications](#specifications)

## Specifications <a id="specifications"></a>

The following example shows a CloudEvent serialized as JSON:

``` JSON
{
    "specversion" : "0.2",
    "type" : "com.github.pull.create",
    "source" : "https://github.com/cloudevents/spec/pull/123",
    "id" : "A234-1234-1234",
    "time" : "2018-04-05T17:31:00Z",
    "comexampleextension1" : "value",
    "comexampleextension2" : {
        "othervalue": 5
    },
    "contenttype" : "text/xml",
    "data" : "<much wow=\"xml\"/>"
}
```

#### Pull Request Event

#### Push Event

#### Pipeline Event
