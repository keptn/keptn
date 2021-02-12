import {Trace} from './trace';
import {async} from "@angular/core/testing";

describe('Trace', () => {

  it('should create instances from json', async(() => {
    let rootTraces: Trace[] = [
      {
      "contenttype": "application/json",
      "data": {
        "deploymentStrategies": {},
        "eventContext": null,
        "helmChart": "H4sIFAAAAAAA/ykAK2FIUjBjSE02THk5NWIzVjBkUzVpWlM5Nk9WVjZNV2xqYW5keVRRbz1IZWxtAOxXS2/iSBDOGYn/UMplTjFtCElkaQ8RsDPsBrB4RFqtVqhjV6A3bXdvdxsNGvHfV35hYzyTkSYPReK7gLuqPle1q9qfPaqMbvXWVBlrSwN+9goghJDrbjf5JYRUf4ndvjqzO1edztV1u33dPSN2u9PungF5jWSqiLSh6oz88r2qxX0QUMnuUWkmQgc2dsNH7SkmTXJ9C1+QB+DFzQGPQoFGtWEeQtI0jZAG6GT/NzkHsWyLNN67qhN+Fun8byiPUL/WAfDc/JPry+r8d69O8/8mYAFdoQO+8J5QWUy0nlCaEL/SQHLUraQ9HGLZxLKbDYWSM4/2RBQaB+z3zv2EX0c6/wYDyalB3fJRcrENMHxBOfDM/MdjX5n/7mW3c5r/t8DFxQU0G2UVQKXUrY3dbDyx0Hegv2+IZiNAQ31qqNNsAJRe/82Gluglq9kRoR349g2s+/S9Uj43YLeL/bRR1OBqmwQBKME5C1cL6VOD2RpAQL8uQrqhjNMHjg6Q1GC2Eh2YlkMSSuToGaGy8IAab31HH5DrPSGVcp8yQN71eUBeHeTu/DC6Gg+wLzuGJ0JDWYiqiLg43KWcJjt0z0tblCzBbnde8XIjzl3Bmbd1YPg4FsZVqJOHkbtJER/RxXV+z7UxslgFkEoY4QnuwLznlg37vF2hjAM35IYUZgw3Ndz9+bK3mM0no6U7nRyQJUrCgfPkNbKUSvyLnvmtVGhisDID7HaQembS8tgz15yFp6GrOj+TbmDmVZxjx66F7XDD8+rcSX85vh0Njur6XYnAKa8CPDLk/hQfK8uZwaVm7cB53llWzF9zw/7AvZv8NRqM569037SR//5UVP7pn5o8epPx/HY4HkyXw9Ht5+M8nm3ZnOjPgTsfx73xx6A3/x5N+tE5psEPSWbzukxe5EloSb26x5HdeDC9H/ZqNyEZ55rAxfhucDv7koQOpsvF9K4uOh5Mp9WKQo5Ur63s98LHTYtKVqLlbIMhau0q8YAHZcUUn9FUSpVJja01Um7WFVPNaAOwkBlGeR853c7QE6GvHbg6cJGomPD3RvvAaFiAIjKFtVtYFVKffeDstYiUh/ogSc4CZvRRg3kyirkJCaqWAAOhtg60yeXNiJWtCv+LUH+PrPsDLpu0L2Oul3n/V/Vfdty+6Kfgc99/HbtT0X+X121y0n9vgWP9V0i/WdoLJd0H9cIvWU5VWY9H2qAaJqIglSaJ9ViV5EOdXR2rE0PVCs2hKNlrvEyjlQXZe2/lCSeccMKHwv8BAAD//7KJCEUAGgAA",
        "project": "sockshop",
        "service": "carts"
      },
      "id": "ade14d1a-338d-4a88-ad88-7b03ec9b5d8e",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T12:45:01.643Z",
      "type": "sh.keptn.event.service.create.started",
      "shkeptncontext": "56255df2-22a2-45ef-b2c4-cf4882096b3f",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "deploymentStrategies": {},
          "eventContext": null,
          "helmChart": "H4sIFAAAAAAA/ykAK2FIUjBjSE02THk5NWIzVjBkUzVpWlM5Nk9WVjZNV2xqYW5keVRRbz1IZWxtAOxXS2/iSBDOGYn/UMplTjFtCElkaQ8RsDPsBrB4RFqtVqhjV6A3bXdvdxsNGvHfV35hYzyTkSYPReK7gLuqPle1q9qfPaqMbvXWVBlrSwN+9goghJDrbjf5JYRUf4ndvjqzO1edztV1u33dPSN2u9PungF5jWSqiLSh6oz88r2qxX0QUMnuUWkmQgc2dsNH7SkmTXJ9C1+QB+DFzQGPQoFGtWEeQtI0jZAG6GT/NzkHsWyLNN67qhN+Fun8byiPUL/WAfDc/JPry+r8d69O8/8mYAFdoQO+8J5QWUy0nlCaEL/SQHLUraQ9HGLZxLKbDYWSM4/2RBQaB+z3zv2EX0c6/wYDyalB3fJRcrENMHxBOfDM/MdjX5n/7mW3c5r/t8DFxQU0G2UVQKXUrY3dbDyx0Hegv2+IZiNAQ31qqNNsAJRe/82Gluglq9kRoR349g2s+/S9Uj43YLeL/bRR1OBqmwQBKME5C1cL6VOD2RpAQL8uQrqhjNMHjg6Q1GC2Eh2YlkMSSuToGaGy8IAab31HH5DrPSGVcp8yQN71eUBeHeTu/DC6Gg+wLzuGJ0JDWYiqiLg43KWcJjt0z0tblCzBbnde8XIjzl3Bmbd1YPg4FsZVqJOHkbtJER/RxXV+z7UxslgFkEoY4QnuwLznlg37vF2hjAM35IYUZgw3Ndz9+bK3mM0no6U7nRyQJUrCgfPkNbKUSvyLnvmtVGhisDID7HaQembS8tgz15yFp6GrOj+TbmDmVZxjx66F7XDD8+rcSX85vh0Njur6XYnAKa8CPDLk/hQfK8uZwaVm7cB53llWzF9zw/7AvZv8NRqM569037SR//5UVP7pn5o8epPx/HY4HkyXw9Ht5+M8nm3ZnOjPgTsfx73xx6A3/x5N+tE5psEPSWbzukxe5EloSb26x5HdeDC9H/ZqNyEZ55rAxfhucDv7koQOpsvF9K4uOh5Mp9WKQo5Ur63s98LHTYtKVqLlbIMhau0q8YAHZcUUn9FUSpVJja01Um7WFVPNaAOwkBlGeR853c7QE6GvHbg6cJGomPD3RvvAaFiAIjKFtVtYFVKffeDstYiUh/ogSc4CZvRRg3kyirkJCaqWAAOhtg60yeXNiJWtCv+LUH+PrPsDLpu0L2Oul3n/V/Vfdty+6Kfgc99/HbtT0X+X121y0n9vgWP9V0i/WdoLJd0H9cIvWU5VWY9H2qAaJqIglSaJ9ViV5EOdXR2rE0PVCs2hKNlrvEyjlQXZe2/lCSeccMKHwv8BAAD//7KJCEUAGgAA",
          "project": "sockshop",
          "service": "carts"
        },
        "id": "ade14d1a-338d-4a88-ad88-7b03ec9b5d8e",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T12:45:01.643Z",
        "type": "sh.keptn.event.service.create.started",
        "shkeptncontext": "56255df2-22a2-45ef-b2c4-cf4882096b3f"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
        }
      },
      "id": "f95d2c20-2d89-4f52-8838-36743cf8835f",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T12:47:05.910Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
          }
        },
        "id": "f95d2c20-2d89-4f52-8838-36743cf8835f",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T12:47:05.910Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.1",
          "teststrategy": "functional"
        },
        "id": "511ebf66-a55c-4159-8866-073035bab8ff",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:48:21.398Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-03-16T12:48:24Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-03-16T12:48:21Z",
          "teststrategy": "functional"
        },
        "id": "fb78aa3c-e769-458a-96ca-d549368e9c2a",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:48:24.775Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T12:48:24Z",
            "timeStart": "2020-03-16T12:48:21Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "c5c64483-c53a-47d5-bce6-9deac74db415",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:48:24.794Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "1d95ba9c-924c-405a-b013-2e07ce0ccd81",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:48:24.799Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
          }
        },
        "id": "662c4b5b-9726-4932-9093-d8d0bbfae257",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:48:24.917Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.1",
          "teststrategy": "performance"
        },
        "id": "59320701-cfa8-442e-bd29-f3eda1255855",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:50:44.542Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T12:52:15Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T12:50:44Z",
          "teststrategy": "performance"
        },
        "id": "44c2456e-7a57-4499-83d8-4af64d01e9c5",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:52:15.772Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T12:52:15Z",
            "timeStart": "2020-03-16T12:50:44Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "teststrategy": "performance"
        },
        "id": "e8ace82f-3a6b-42a3-bc6b-4dcb93c741d7",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:52:15.818Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging"
        },
        "id": "d0905c49-a1dd-456c-aa04-efe3065709a2",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:52:15.829Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "production",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
          }
        },
        "id": "391dc227-99a0-4dc8-9176-282f618aa3bd",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:52:15.969Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-production",
          "deploymentURIPublic": "http://carts.sockshop-production.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "production",
          "tag": "0.10.1",
          "teststrategy": ""
        },
        "id": "49203fd8-456f-4fa9-a486-615a639ee334",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:54:41.891Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T12:54:41Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "production",
          "start": "2020-03-16T12:54:41Z",
          "teststrategy": ""
        },
        "id": "c5c6530d-7f2e-4cb4-bb49-2e82f4377507",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:54:41.978Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no test has been executed",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T12:54:41Z",
            "timeStart": "2020-03-16T12:54:41Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "production",
          "teststrategy": ""
        },
        "id": "90c644d1-76b8-433f-afe1-b32348693343",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:54:41.986Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "production"
        },
        "id": "e560fddd-793b-4328-b442-fb01d3d4cdf8",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T12:54:41.995Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
        }
      },
      "id": "d72b0a9d-24d3-48ec-832f-e36696c466d4",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T13:48:48.749Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "d72b0a9d-24d3-48ec-832f-e36696c466d4",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T13:48:48.749Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.2",
          "teststrategy": "functional"
        },
        "id": "d66841ef-cc33-4de4-ac77-1fe12a5e6af9",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:49:55.878Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-03-16T13:49:58Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-03-16T13:49:55Z",
          "teststrategy": "functional"
        },
        "id": "c7aeac48-322b-4fa3-af71-eb9ce9fb1111",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:49:58.766Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T13:49:58Z",
            "timeStart": "2020-03-16T13:49:55Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "078a59c1-cd7a-4df8-b993-7872191f2c1f",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:49:58.790Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "3e807c9b-3b1d-45ed-9b0e-655a8686842b",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:49:58.800Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "d9073b98-98f5-47de-966f-8d57902f2cf3",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:49:58.880Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.2",
          "teststrategy": "performance"
        },
        "id": "6a9f572e-4a62-484f-99a8-a578efd6cc15",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T13:51:12.181Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T14:08:06Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T13:51:12Z",
          "teststrategy": "performance"
        },
        "id": "b8520739-150b-4897-b912-6d4b9316b2da",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T14:08:06.568Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T14:08:06Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-03-16T13:51:12Z",
          "teststrategy": "performance"
        },
        "id": "6fc28d4e-77a5-496f-82da-a09b7f7709c5",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T14:08:06.683Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "6f8bf547-d588-4866-b5f5-1fe43a4b0c65"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
        }
      },
      "id": "21ccce4c-e77e-479c-91fa-7052551bdf48",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T15:09:02.740Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "21ccce4c-e77e-479c-91fa-7052551bdf48",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:02.740Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.2",
          "teststrategy": "functional"
        },
        "id": "9c85fbf9-0aaa-418b-b246-bf0e4683d169",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:05.979Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-03-16T15:09:08Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-03-16T15:09:05Z",
          "teststrategy": "functional"
        },
        "id": "36b028db-d51d-4828-b037-6c060773b5e2",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:08.808Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T15:09:08Z",
            "timeStart": "2020-03-16T15:09:05Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "72da36b1-c086-4362-988b-dfa2905c75bc",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:08.881Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "a15b3223-9d3b-4294-97b1-1cb68f72121b",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:08.891Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "9bcbce99-9e00-493f-8143-83f1a2dfaf82",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:08.982Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.2",
          "teststrategy": "performance"
        },
        "id": "005c2e55-5533-46f1-9ec7-d7b8d4274c92",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:09:14.177Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T15:26:09Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T15:09:14Z",
          "teststrategy": "performance"
        },
        "id": "4a05c4b5-6efa-4207-95e6-131cb93eb485",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:26:09.375Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T15:26:09Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-03-16T15:09:14Z",
          "teststrategy": "performance"
        },
        "id": "b59ce0aa-c5c4-40ab-9388-ccb70be22ba4",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:26:09.436Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "934b6d92-7605-4eaa-b471-d6022f3e6d72"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
        }
      },
      "id": "84235fd5-d009-4e42-96c2-4cfdcfb48b4d",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T15:54:52.249Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "84235fd5-d009-4e42-96c2-4cfdcfb48b4d",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:52.249Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.2",
          "teststrategy": "functional"
        },
        "id": "54a95360-32b5-4bb5-8fdd-ca3e91ee830d",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:55.321Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-03-16T15:54:58Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-03-16T15:54:55Z",
          "teststrategy": "functional"
        },
        "id": "0d1a50da-194d-409e-81be-21fb9dce1dac",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:58.074Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T15:54:58Z",
            "timeStart": "2020-03-16T15:54:55Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "9ed2ecaf-fc50-4fd7-83f9-7b3a3be09345",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:58.096Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "1cc077e9-2f6c-4329-ac0f-15f956f22f94",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:58.106Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.2"}
          }
        },
        "id": "ef7b141e-23a3-4542-86c1-ca3dcb8341cc",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:54:58.218Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.2",
          "teststrategy": "performance"
        },
        "id": "a110f4f6-aeff-4800-bf4d-6e72edd2e260",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T15:55:03.690Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:11:54Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T15:55:03Z",
          "teststrategy": "performance"
        },
        "id": "b7c07a50-ef29-4be8-9b87-93151b0ece40",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:11:54.302Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:11:54Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-03-16T15:55:03Z",
          "teststrategy": "performance"
        },
        "id": "9b99ca6f-cc66-44c1-8881-dac8cc8340db",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:11:54.377Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:11:54Z",
          "indicatorValues": [{
            "metric": "response_time_p95",
            "success": true,
            "value": 2013.11088577307
          }, {"metric": "response_time_p50", "success": true, "value": 2010.8712390903988}, {
            "metric": "error_rate",
            "success": true,
            "value": 139
          }, {"metric": "throughput", "success": true, "value": 3}, {
            "metric": "cpu_usage",
            "success": true,
            "value": 8.50832800771676
          }],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T15:55:03Z",
          "teststrategy": "performance"
        },
        "id": "41ee64d9-01f1-4fef-adcd-f22f67a24317",
        "source": "dynatrace-sli-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:13:56.182Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "evaluation": {
            "indicatorResults": [{
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "<=900", "targetValue": 900, "violated": true}, {
                "criteria": "<600",
                "targetValue": 600,
                "violated": true
              }],
              "value": {"metric": "response_time_p95", "success": true, "value": 2013.11088577307}
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "<=800", "targetValue": 800, "violated": true}, {
                "criteria": "<300",
                "targetValue": 300,
                "violated": true
              }],
              "value": {"metric": "response_time_p50", "success": true, "value": 2010.8712390903988}
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
              "value": {"metric": "error_rate", "success": true, "value": 139}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {"metric": "throughput", "success": true, "value": 3}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {"metric": "cpu_usage", "success": true, "value": 8.50832800771676}
            }],
            "result": "fail",
            "score": 0,
            "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
            "timeEnd": "2020-03-16T16:11:54Z",
            "timeStart": "2020-03-16T15:55:03Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "teststrategy": "performance"
        },
        "id": "cdcf2cf8-f2f9-4d7d-b62a-d27bba980389",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:13:56.288Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "discard"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging"
        },
        "id": "aa063da1-67e9-44e5-9275-08d87a7e0af1",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:13:56.300Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
        }
      },
      "id": "6eacae7f-7446-49a8-b351-0d95bb6be76b",
      "source": "https://github.com/keptn/keptn/cli#configuration-change",
      "specversion": "0.2",
      "time": "2020-05-28T07:49:29.742Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
          }
        },
        "id": "6eacae7f-7446-49a8-b351-0d95bb6be76b",
        "source": "https://github.com/keptn/keptn/cli#configuration-change",
        "specversion": "0.2",
        "time": "2020-05-28T07:49:29.742Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.1",
          "teststrategy": "functional"
        },
        "id": "07edecb8-3add-4c4c-859a-484d33554a94",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:50:45.299Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-05-28T07:50:51Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-05-28T07:50:45Z",
          "teststrategy": "functional"
        },
        "id": "1a092116-422f-4b14-b087-4be38377e9e9",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:50:51.975Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-05-28T07:50:51Z",
            "timeStart": "2020-05-28T07:50:45Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "e276d1dc-677c-4193-885c-e9fe5a95e0e5",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:50:52.018Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "791596f6-cf13-4b50-8baa-a2558972bb6c",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:50:52.030Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.1"}
          }
        },
        "id": "e36f86f7-b2d2-41cc-a9a1-10df93197623",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:50:52.181Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.1",
          "teststrategy": "performance"
        },
        "id": "acee690b-2616-4056-aeba-298355c38845",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:52:07.546Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-05-28T07:54:12Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-05-28T07:52:07Z",
          "teststrategy": "performance"
        },
        "id": "60b25fe0-110b-4d97-aa73-cfa3ea9bd734",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:54:12.734Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-05-28T07:54:12Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-05-28T07:52:07Z",
          "teststrategy": "performance"
        },
        "id": "6914860b-f244-45e5-9c4e-0286f52adb59",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:54:12.778Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-05-28T07:54:12Z",
          "indicatorValues": [{
            "metric": "response_time_p95",
            "success": true,
            "value": 345.8670349893803
          }, {"metric": "response_time_p50", "success": true, "value": 186.38390148907558}, {
            "metric": "error_rate",
            "success": true,
            "value": 5
          }, {
            "metric": "throughput",
            "success": true,
            "value": 4
          }, {
            "message": "Dynatrace API returned status code 403: The query involves 1794 metrics, but the limit for REST queries is 10. Consider splitting your query into multiple smaller queries.",
            "metric": "cpu_usage",
            "success": false,
            "value": 0
          }],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "start": "2020-05-28T07:52:07Z",
          "teststrategy": "performance"
        },
        "id": "360b819e-7050-4045-b665-121d26786c58",
        "source": "dynatrace-sli-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:56:15.777Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "evaluation": {
            "indicatorResults": [{
              "score": 1,
              "status": "pass",
              "targets": [{"criteria": "<600", "targetValue": 600, "violated": false}],
              "value": {"metric": "response_time_p95", "success": true, "value": 345.8670349893803}
            }, {
              "score": 1,
              "status": "pass",
              "targets": [{"criteria": "<300", "targetValue": 300, "violated": false}],
              "value": {"metric": "response_time_p50", "success": true, "value": 186.38390148907558}
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
              "value": {"metric": "error_rate", "success": true, "value": 5}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {"metric": "throughput", "success": true, "value": 4}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {
                "message": "Dynatrace API returned status code 403: The query involves 1794 metrics, but the limit for REST queries is 10. Consider splitting your query into multiple smaller queries.",
                "metric": "cpu_usage",
                "success": false,
                "value": 0
              }
            }],
            "result": "fail",
            "score": 66.66666666666666,
            "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
            "timeEnd": "2020-05-28T07:54:12Z",
            "timeStart": "2020-05-28T07:52:07Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "teststrategy": "performance"
        },
        "id": "214ca172-4080-4165-9a4d-39f399b17a45",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:56:15.840Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "discard"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging"
        },
        "id": "ab1ea10a-90e7-43c2-9824-222acf798b9e",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:56:15.852Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "",
        "end": "2020-05-28T08:59:00.000Z",
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "staging",
        "start": "2020-05-28T08:54:00.000Z",
        "teststrategy": "manual"
      },
      "id": "0de64842-b277-4629-a5b5-c87829ff1151",
      "source": "https://github.com/keptn/keptn/cli#configuration-change",
      "specversion": "0.2",
      "time": "2020-05-28T07:53:06.860Z",
      "type": "sh.keptn.event.start-evaluation",
      "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "",
          "end": "2020-05-28T08:59:00.000Z",
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "start": "2020-05-28T08:54:00.000Z",
          "teststrategy": "manual"
        },
        "id": "0de64842-b277-4629-a5b5-c87829ff1151",
        "source": "https://github.com/keptn/keptn/cli#configuration-change",
        "specversion": "0.2",
        "time": "2020-05-28T07:53:06.860Z",
        "type": "sh.keptn.event.start-evaluation",
        "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "",
          "deploymentstrategy": "",
          "end": "2020-05-28T08:59:00.000Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-05-28T08:54:00.000Z",
          "teststrategy": "manual"
        },
        "id": "128992d7-51ac-4f0d-bb66-e2282d3a33e3",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:53:06.919Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09"
      }, {
        "contenttype": "application/json",
        "data": {
          "deployment": "",
          "deploymentstrategy": "",
          "end": "2020-05-28T08:59:00.000Z",
          "indicatorValues": [{
            "message": "end time must not be in the future",
            "metric": "response_time_p95",
            "success": false,
            "value": 0
          }, {
            "message": "end time must not be in the future",
            "metric": "response_time_p50",
            "success": false,
            "value": 0
          }, {
            "message": "end time must not be in the future",
            "metric": "error_rate",
            "success": false,
            "value": 0
          }, {
            "message": "end time must not be in the future",
            "metric": "throughput",
            "success": false,
            "value": 0
          }, {"message": "end time must not be in the future", "metric": "cpu_usage", "success": false, "value": 0}],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "start": "2020-05-28T08:54:00.000Z",
          "teststrategy": "manual"
        },
        "id": "b2ac26ca-5a26-4a6b-8b6a-de98e8393540",
        "source": "dynatrace-sli-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:53:07.119Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "",
          "evaluation": {
            "indicatorResults": [{
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "<=900", "targetValue": 0, "violated": true}, {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }],
              "value": {
                "message": "end time must not be in the future",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              }
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "<=800", "targetValue": 0, "violated": true}, {
                "criteria": "<300",
                "targetValue": 0,
                "violated": true
              }],
              "value": {
                "message": "end time must not be in the future",
                "metric": "response_time_p50",
                "success": false,
                "value": 0
              }
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
              "value": {
                "message": "end time must not be in the future",
                "metric": "error_rate",
                "success": false,
                "value": 0
              }
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {
                "message": "end time must not be in the future",
                "metric": "throughput",
                "success": false,
                "value": 0
              }
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {
                "message": "end time must not be in the future",
                "metric": "cpu_usage",
                "success": false,
                "value": 0
              }
            }],
            "result": "fail",
            "score": 0,
            "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
            "timeEnd": "2020-05-28T08:59:00.000Z",
            "timeStart": "2020-05-28T08:54:00.000Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "teststrategy": "manual"
        },
        "id": "1f2c7acc-c09e-424c-b4a0-1f8748657ee8",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-05-28T07:53:07.193Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09"
      }]
    }, {
      "contenttype": "application/json",
      "data": {
        "canary": {"action": "set", "value": 100},
        "eventContext": null,
        "labels": null,
        "project": "sockshop",
        "service": "carts",
        "stage": "",
        "configurationChange": {
          "values": {"image": "docker.io/keptnexamples/carts:0.10.3"}
        }
      },
      "id": "d9949721-29de-47bf-9b04-b890eeb21e4d",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-03-16T16:14:29.785Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583",
      "rootTraces": [{
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "eventContext": null,
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.3"}
          }
        },
        "id": "d9949721-29de-47bf-9b04-b890eeb21e4d",
        "source": "https://github.com/keptn/keptn/api",
        "specversion": "0.2",
        "time": "2020-03-16T16:14:29.785Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts.sockshop-dev",
          "deploymentURIPublic": "http://carts.sockshop-dev.104.197.225.41.xip.io",
          "deploymentstrategy": "direct",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "tag": "0.10.3",
          "teststrategy": "functional"
        },
        "id": "aca28645-27db-4605-86ec-7a429b8d40c6",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:15:40.573Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "end": "2020-03-16T16:15:43Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "start": "2020-03-16T16:15:40Z",
          "teststrategy": "functional"
        },
        "id": "4f3cdd63-a80e-45a4-b518-4514aa17e4e1",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:15:43.542Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "direct",
          "evaluation": {
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO found for service carts",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2020-03-16T16:15:43Z",
            "timeStart": "2020-03-16T16:15:40Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "teststrategy": "functional"
        },
        "id": "fa2f322c-8ef6-4800-93da-e633408f8b0a",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:15:43.585Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "promote"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "dev"
        },
        "id": "44dbf3cf-d3e2-4115-b5f2-b12a774f2c8d",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:15:43.595Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "set", "value": 100},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "configurationChange": {
            "values": {"image": "docker.io/keptnexamples/carts:0.10.3"}
          }
        },
        "id": "0c2f02b6-4116-4bda-b774-4af2930f89d9",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:15:43.685Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentURILocal": "http://carts-canary.sockshop-staging",
          "deploymentURIPublic": "http://carts.sockshop-staging.104.197.225.41.xip.io",
          "deploymentstrategy": "blue_green_service",
          "image": "docker.io/keptnexamples/carts",
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "tag": "0.10.3",
          "teststrategy": "performance"
        },
        "id": "5c30bdda-6ca0-42d6-8db7-e81df0fd5ba1",
        "source": "helm-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:16:56.913Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:18:48Z",
          "labels": null,
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T16:16:56Z",
          "teststrategy": "performance"
        },
        "id": "f37ae270-4509-4c8f-a43c-48afbdd6f5d9",
        "source": "jmeter-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:18:48.795Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "customFilters": [],
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:18:48Z",
          "indicators": ["response_time_p95", "response_time_p50", "error_rate", "throughput", "cpu_usage"],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "sliProvider": "dynatrace",
          "stage": "staging",
          "start": "2020-03-16T16:16:56Z",
          "teststrategy": "performance"
        },
        "id": "ef822670-549b-43bf-8c17-8e951054f951",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:18:48.847Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deployment": "canary",
          "deploymentstrategy": "blue_green_service",
          "end": "2020-03-16T16:18:48Z",
          "indicatorValues": [{
            "metric": "response_time_p95",
            "success": true,
            "value": 339.15595320978224
          }, {"metric": "response_time_p50", "success": true, "value": 158.36151562776496}, {
            "metric": "error_rate",
            "success": true,
            "value": 3.3333333333333335
          }, {"metric": "throughput", "success": true, "value": 3}, {
            "metric": "cpu_usage",
            "success": true,
            "value": 15.433689541286892
          }],
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "start": "2020-03-16T16:16:56Z",
          "teststrategy": "performance"
        },
        "id": "6502134b-3a34-4807-a485-935ece25e3cf",
        "source": "dynatrace-sli-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:20:50.165Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "deploymentstrategy": "blue_green_service",
          "evaluation": {
            "indicatorResults": [{
              "score": 1,
              "status": "pass",
              "targets": [{"criteria": "<600", "targetValue": 600, "violated": false}],
              "value": {"metric": "response_time_p95", "success": true, "value": 339.15595320978224}
            }, {
              "score": 1,
              "status": "pass",
              "targets": [{"criteria": "<300", "targetValue": 300, "violated": false}],
              "value": {"metric": "response_time_p50", "success": true, "value": 158.36151562776496}
            }, {
              "score": 0,
              "status": "fail",
              "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
              "value": {"metric": "error_rate", "success": true, "value": 3.3333333333333335}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {"metric": "throughput", "success": true, "value": 3}
            }, {
              "score": 0,
              "status": "info",
              "targets": null,
              "value": {"metric": "cpu_usage", "success": true, "value": 15.433689541286892}
            }],
            "result": "fail",
            "score": 66.66666666666666,
            "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
            "timeEnd": "2020-03-16T16:18:48Z",
            "timeStart": "2020-03-16T16:16:56Z"
          },
          "labels": null,
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "teststrategy": "performance"
        },
        "id": "725206de-f7bf-4461-b758-070b77402dc3",
        "source": "lighthouse-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:20:50.233Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }, {
        "contenttype": "application/json",
        "data": {
          "canary": {"action": "discard"},
          "labels": null,
          "project": "sockshop",
          "service": "carts",
          "stage": "staging"
        },
        "id": "04ef469d-2ef0-403e-9a3c-81584fc23b30",
        "source": "gatekeeper-service",
        "specversion": "0.2",
        "time": "2020-03-16T16:20:50.280Z",
        "type": "sh.keptn.event.dev.artifact-delivery.triggered",
        "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
      }]
    }, {
      "data": {
        "configurationChange": {
          "values": {}
        },
        "deployment": {
          "deploymentstrategy": ""
        },
        "project": "keptn",
        "service": "control-plane",
        "stage": "dev"
      },
      "id": "3b209b06-597c-413e-9401-b80e4855a313",
      "source": "https://github.com/keptn/keptn/cli#configuration-change",
      "specversion": "1.0",
      "time": "2021-02-02T08:52:39.186Z",
      "type": "sh.keptn.event.dev.artifact-delivery.triggered",
      "shkeptncontext": "0ede19b7-dc65-4f04-9882-ddadf3703019"
    }].map(trace => Trace.fromJSON(trace));
    let evaluationTraces: Trace[] = [{
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "blue_green_service",
        "evaluation": {
          "indicatorResults": null,
          "result": "no evaluation performed by lighthouse because no SLO found for service carts",
          "score": 0,
          "sloFileContent": "",
          "timeEnd": "2020-03-16T12:52:15Z",
          "timeStart": "2020-03-16T12:50:44Z"
        },
        "labels": null,
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "teststrategy": "performance"
      },
      "id": "e8ace82f-3a6b-42a3-bc6b-4dcb93c741d7",
      "source": "lighthouse-service",
      "specversion": "0.2",
      "time": "2020-03-16T12:52:15.818Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "218ddbfa-ed09-4cf9-887a-167a334a76d0"
    }, {
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "blue_green_service",
        "evaluation": {
          "indicatorResults": [{
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "<=900", "targetValue": 900, "violated": true}, {
              "criteria": "<600",
              "targetValue": 600,
              "violated": true
            }],
            "value": {"metric": "response_time_p95", "success": true, "value": 2013.11088577307}
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "<=800", "targetValue": 800, "violated": true}, {
              "criteria": "<300",
              "targetValue": 300,
              "violated": true
            }],
            "value": {"metric": "response_time_p50", "success": true, "value": 2010.8712390903988}
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
            "value": {"metric": "error_rate", "success": true, "value": 139}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {"metric": "throughput", "success": true, "value": 3}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {"metric": "cpu_usage", "success": true, "value": 8.50832800771676}
          }],
          "result": "fail",
          "score": 0,
          "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
          "timeEnd": "2020-03-16T16:11:54Z",
          "timeStart": "2020-03-16T15:55:03Z"
        },
        "labels": null,
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "teststrategy": "performance"
      },
      "id": "cdcf2cf8-f2f9-4d7d-b62a-d27bba980389",
      "source": "lighthouse-service",
      "specversion": "0.2",
      "time": "2020-03-16T16:13:56.288Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "38bdaef6-60a2-48f3-b474-626a683d175c"
    }, {
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "blue_green_service",
        "evaluation": {
          "indicatorResults": [{
            "score": 1,
            "status": "pass",
            "targets": [{"criteria": "<600", "targetValue": 600, "violated": false}],
            "value": {"metric": "response_time_p95", "success": true, "value": 339.15595320978224}
          }, {
            "score": 1,
            "status": "pass",
            "targets": [{"criteria": "<300", "targetValue": 300, "violated": false}],
            "value": {"metric": "response_time_p50", "success": true, "value": 158.36151562776496}
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
            "value": {"metric": "error_rate", "success": true, "value": 3.3333333333333335}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {"metric": "throughput", "success": true, "value": 3}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {"metric": "cpu_usage", "success": true, "value": 15.433689541286892}
          }],
          "result": "fail",
          "score": 66.66666666666666,
          "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
          "timeEnd": "2020-03-16T16:18:48Z",
          "timeStart": "2020-03-16T16:16:56Z"
        },
        "labels": null,
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "teststrategy": "performance"
      },
      "id": "725206de-f7bf-4461-b758-070b77402dc3",
      "source": "lighthouse-service",
      "specversion": "0.2",
      "time": "2020-03-16T16:20:50.233Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "42e8e409-5afc-4ee5-abdb-f41926ab2583"
    }, {
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "",
        "evaluation": {
          "indicatorResults": [{
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "<=900", "targetValue": 0, "violated": true}, {
              "criteria": "<600",
              "targetValue": 0,
              "violated": true
            }],
            "value": {
              "message": "end time must not be in the future",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            }
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "<=800", "targetValue": 0, "violated": true}, {
              "criteria": "<300",
              "targetValue": 0,
              "violated": true
            }],
            "value": {
              "message": "end time must not be in the future",
              "metric": "response_time_p50",
              "success": false,
              "value": 0
            }
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
            "value": {
              "message": "end time must not be in the future",
              "metric": "error_rate",
              "success": false,
              "value": 0
            }
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "message": "end time must not be in the future",
              "metric": "throughput",
              "success": false,
              "value": 0
            }
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "message": "end time must not be in the future",
              "metric": "cpu_usage",
              "success": false,
              "value": 0
            }
          }],
          "result": "fail",
          "score": 0,
          "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
          "timeEnd": "2020-05-28T08:59:00.000Z",
          "timeStart": "2020-05-28T08:54:00.000Z"
        },
        "labels": null,
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "teststrategy": "manual"
      },
      "id": "1f2c7acc-c09e-424c-b4a0-1f8748657ee8",
      "source": "lighthouse-service",
      "specversion": "0.2",
      "time": "2020-05-28T07:53:07.193Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "3302455e-ffe6-4ae3-a514-31150557dd09"
    }, {
      "contenttype": "application/json",
      "data": {
        "deploymentstrategy": "blue_green_service",
        "evaluation": {
          "indicatorResults": [{
            "score": 1,
            "status": "pass",
            "targets": [{"criteria": "<600", "targetValue": 600, "violated": false}],
            "value": {"metric": "response_time_p95", "success": true, "value": 345.8670349893803}
          }, {
            "score": 1,
            "status": "pass",
            "targets": [{"criteria": "<300", "targetValue": 300, "violated": false}],
            "value": {"metric": "response_time_p50", "success": true, "value": 186.38390148907558}
          }, {
            "score": 0,
            "status": "fail",
            "targets": [{"criteria": "=0", "targetValue": 0, "violated": true}],
            "value": {"metric": "error_rate", "success": true, "value": 5}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {"metric": "throughput", "success": true, "value": 4}
          }, {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "message": "Dynatrace API returned status code 403: The query involves 1794 metrics, but the limit for REST queries is 10. Consider splitting your query into multiple smaller queries.",
              "metric": "cpu_usage",
              "success": false,
              "value": 0
            }
          }],
          "result": "fail",
          "score": 66.66666666666666,
          "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
          "timeEnd": "2020-05-28T07:54:12Z",
          "timeStart": "2020-05-28T07:52:07Z"
        },
        "labels": null,
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "teststrategy": "performance"
      },
      "id": "214ca172-4080-4165-9a4d-39f399b17a45",
      "source": "lighthouse-service",
      "specversion": "0.2",
      "time": "2020-05-28T07:56:15.840Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "fea3dc8c-5a85-435a-a86d-cee0b62f248e"
    }].map(trace => Trace.fromJSON(trace));

    expect(rootTraces[0] instanceof Trace).toBe(true, 'instance of Trace');

    expect(rootTraces[0].type).toBe('sh.keptn.event.service.create.started');
    expect(rootTraces[0].getLabel()).toBe('create', 'Label for trace "sh.keptn.event.service.create.started" should be "create"');
    expect(rootTraces[0].getShortImageName()).toBe(undefined);
    expect(rootTraces[0].getIcon()).toBe('information', 'Icon for trace "sh.keptn.event.service.create.started" should be "information"');
    expect(rootTraces[0].isFaulty()).toBe(null);
    expect(rootTraces[0].isWarning()).toBe(null);
    expect(rootTraces[0].isSuccessful()).toBe(null);
    expect(rootTraces[0].getProject()).toBe('sockshop');
    expect(rootTraces[0].getService()).toBe('carts');

    expect(rootTraces[1].type).toBe('sh.keptn.event.dev.artifact-delivery.triggered');
    expect(rootTraces[1].getLabel()).toBe('artifact-delivery', 'Label for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "artifact-delivery"');
    expect(rootTraces[1].getShortImageName()).toBe('carts:0.10.1', 'ShortImageName for first trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "carts:0.10.1"');
    expect(rootTraces[1].getIcon()).toBe('duplicate', 'Icon for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "duplicate"');
    expect(rootTraces[1].isFaulty()).toBe(null);
    expect(rootTraces[1].isWarning()).toBe(null);
    expect(rootTraces[1].isSuccessful()).toBe(null);
    expect(rootTraces[1].getProject()).toBe('sockshop');
    expect(rootTraces[1].getService()).toBe('carts');

    expect(rootTraces[2].type).toBe('sh.keptn.event.dev.artifact-delivery.triggered');
    expect(rootTraces[2].getLabel()).toBe('artifact-delivery', 'Label for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "artifact-delivery"');
    expect(rootTraces[2].getShortImageName()).toBe('carts:0.10.2', 'ShortImageName for second trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "carts:0.10.2"');
    expect(rootTraces[2].getIcon()).toBe('duplicate', 'Icon for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "duplicate"');
    expect(rootTraces[2].isFaulty()).toBe(null);
    expect(rootTraces[2].isWarning()).toBe(null);
    expect(rootTraces[2].isSuccessful()).toBe(null);
    expect(rootTraces[2].getProject()).toBe('sockshop');
    expect(rootTraces[2].getService()).toBe('carts');

    expect(rootTraces[8].type).toBe('sh.keptn.event.dev.artifact-delivery.triggered');
    expect(rootTraces[8].getLabel()).toBe('artifact-delivery', 'Label for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "artifact-delivery"');
    expect(rootTraces[8].getShortImageName()).toBe(undefined, 'ShortImageName for third trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "undefined"');
    expect(rootTraces[8].getIcon()).toBe('duplicate', 'Icon for trace "sh.keptn.event.dev.artifact-delivery.triggered" should be "duplicate"');
    expect(rootTraces[8].isFaulty()).toBe(null);
    expect(rootTraces[8].isWarning()).toBe(null);
    expect(rootTraces[8].isSuccessful()).toBe(null);
    expect(rootTraces[8].getProject()).toBe('keptn');
    expect(rootTraces[8].getService()).toBe('control-plane');

    expect(evaluationTraces[0].type).toBe('sh.keptn.event.evaluation.finished');
    expect(evaluationTraces[0].getLabel()).toBe('Evaluation finished', 'Label for trace "sh.keptn.event.evaluation.finished" should be "Evaluation finished"');
    expect(evaluationTraces[0].getIcon()).toBe('traffic-light', 'Icon for trace "sh.keptn.event.evaluation.finished" should be "traffic-light"');
    expect(evaluationTraces[0].isFaulty()).toBe(null);
    expect(evaluationTraces[0].isWarning()).toBe(null);
    expect(evaluationTraces[0].isSuccessful()).toBe(evaluationTraces[0].data.stage, 'Successful evaluation should return true');
    expect(evaluationTraces[0].getProject()).toBe('sockshop');
    expect(evaluationTraces[0].getService()).toBe('carts');
  }));
});
