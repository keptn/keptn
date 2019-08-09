# Shipyard Service

The shipyard-service is a keptn core component. It is responsible for creating a project and processing a shipyard file that defines the stages each deployment has to go through until it is released into production. The definition of a shipyard file is provided [here](https://github.com/keptn/keptn/blob/develop/specification/shipyard.md).

The shipyard-service listens to keptn events of type:
- [`sh.keptn.internal.events.project.create`](https://github.com/keptn/keptn/blob/develop/specification/cloudevents.md#create-project)

When receiving such an event, the shipyard-service processes the payload in the data block of the event. Thereby, it uses the API of the configuration-service to create the specified entities and to finally store the payload as shipyad.yaml.