# Keptn Debug-UI

The Debug-UI was built to debug keptn.

**Features:**
 - See all sequences of a project 
 - Get all the events of a specific sequence 
 - Get the sequences which are blocking a specific sequence 
 - Export a MongoDB collection

The Debug-UI can be exposed on port 9090 of the shipyard controller by setting 'DEBUGUI_ENABLED' in the deploy/service.yml to "true".
After enabling it, you can portforward port 9090 of the shipyard-controller to your local host: `kubectl -n keptn port-forward service/shipyard-controller 9090:9090`

All the endpoints are documented in a [Swagger-ui](http://localhost:9090/swagger-ui/)

![debugui](./debugui.png)

## View the events of specific sequence

By clicking on the View Events button you are able to see a list of all events relevant to that specific sequence.

![debugui](./viewevents.png)

## Get all blocking sequences 

If a sequence is currently waiting to be executed, the `getblockingsequences` button will show a list of all blocking sequences.

## Downloading a MongoDB collection

This will download the json data from the selected MongoDB collection. This can be then imported into MongoDB-compass for further debugging by opening an empty Database and then clicking on import.

![debugui](./dbdump.png)

# FAQ
## How to snakeshot in balanka

* Trap the ball under the man’s foot.
* Make sure that the point of contact between the foot and the ball is at, or very near, the top of the ball. This will prevent the ball from squeezing out prematurely.
* Once the ball is secure, place your wrist on the handle.
* Pull or push the handle with your wrist to move the ball sideways.
* Turn the handle with your wrist, making the man do a backward somersault before striking the ball.
