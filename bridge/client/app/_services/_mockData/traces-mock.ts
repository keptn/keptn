import {Trace} from "../../_models/trace";

export const traceData = [
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "message": "Finished release",
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "afeaf817-6600-4925-a89b-1086b9baf18b",
        "source": "shipyard-controller",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:33.553Z",
        "type": "sh.keptn.event.dev.delivery.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "shkeptnspecversion": "0.2.0",
        "triggeredid": "cd9015ff-fad0-43d3-abc0-d50455aa8a82",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "message": "Finished release",
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "afeaf817-6600-4925-a89b-1086b9baf18b",
          "source": "shipyard-controller",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:33.553Z",
          "type": "sh.keptn.event.dev.delivery.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "shkeptnspecversion": "0.2.0",
          "triggeredid": "cd9015ff-fad0-43d3-abc0-d50455aa8a82"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentURIsLocal": null,
        "deploymentstrategy": ""
      },
      "project": "sockshop",
      "service": "carts",
      "stage": "dev"
    },
    "id": "cd9015ff-fad0-43d3-abc0-d50455aa8a82",
    "source": "https://github.com/keptn/keptn/cli#configuration-change",
    "specversion": "1.0",
    "time": "2021-03-16T12:39:53.602Z",
    "type": "sh.keptn.event.dev.delivery.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentURIsLocal": null,
          "deploymentstrategy": ""
        },
        "project": "sockshop",
        "service": "carts",
        "stage": "dev"
      },
      "id": "cd9015ff-fad0-43d3-abc0-d50455aa8a82",
      "source": "https://github.com/keptn/keptn/cli#configuration-change",
      "specversion": "1.0",
      "time": "2021-03-16T12:39:53.602Z",
      "type": "sh.keptn.event.dev.delivery.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "fdc4c95c-163c-4c3b-94a2-6e57e314fca5",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:39:54.219Z",
        "type": "sh.keptn.event.deployment.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "a030c58a-9446-45f6-a6cd-ba00e9fc3729",
        "plainEvent": {
          "data": {
            "project": "sockshop",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "fdc4c95c-163c-4c3b-94a2-6e57e314fca5",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:39:54.219Z",
          "type": "sh.keptn.event.deployment.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "a030c58a-9446-45f6-a6cd-ba00e9fc3729"
        }
      },
      {
        "traces": [],
        "data": {
          "deployment": {
            "deploymentNames": [
              "direct"
            ],
            "deploymentURIsLocal": [
              "http://carts.sockshop-dev:80"
            ],
            "deploymentURIsPublic": [
              "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
            ],
            "deploymentstrategy": "direct",
            "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
          },
          "message": "Successfully deployed",
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "fde6a8c8-a2f0-4bbb-9abf-1d36208d6ed9",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:41:12.039Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "a030c58a-9446-45f6-a6cd-ba00e9fc3729",
        "plainEvent": {
          "data": {
            "deployment": {
              "deploymentNames": [
                "direct"
              ],
              "deploymentURIsLocal": [
                "http://carts.sockshop-dev:80"
              ],
              "deploymentURIsPublic": [
                "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
              ],
              "deploymentstrategy": "direct",
              "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
            },
            "message": "Successfully deployed",
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "fde6a8c8-a2f0-4bbb-9abf-1d36208d6ed9",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:41:12.039Z",
          "type": "sh.keptn.event.deployment.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "a030c58a-9446-45f6-a6cd-ba00e9fc3729"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentURIsLocal": null,
        "deploymentstrategy": "direct"
      },
      "message": "",
      "project": "sockshop",
      "result": "",
      "service": "carts",
      "stage": "dev",
      "status": ""
    },
    "id": "a030c58a-9446-45f6-a6cd-ba00e9fc3729",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:39:54.198Z",
    "type": "sh.keptn.event.deployment.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentURIsLocal": null,
          "deploymentstrategy": "direct"
        },
        "message": "",
        "project": "sockshop",
        "result": "",
        "service": "carts",
        "stage": "dev",
        "status": ""
      },
      "id": "a030c58a-9446-45f6-a6cd-ba00e9fc3729",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:39:54.198Z",
      "type": "sh.keptn.event.deployment.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true,
    "image": "carts:0.12.3"
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "19dad0eb-166e-4bc9-b972-9488c7df6e9a",
        "source": "jmeter-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:41:12.842Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "111654f4-ee87-4650-99d2-dbf1b612209b",
        "plainEvent": {
          "data": {
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "19dad0eb-166e-4bc9-b972-9488c7df6e9a",
          "source": "jmeter-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:41:12.842Z",
          "type": "sh.keptn.event.test.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "111654f4-ee87-4650-99d2-dbf1b612209b"
        }
      },
      {
        "traces": [],
        "data": {
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded",
          "test": {
            "end": "2021-03-16T12:41:24Z",
            "gitCommit": "",
            "start": "2021-03-16T12:41:12Z"
          }
        },
        "id": "40ada5cd-a5ef-491b-8a6a-21d3aff6d730",
        "source": "jmeter-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:41:24.389Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "111654f4-ee87-4650-99d2-dbf1b612209b",
        "plainEvent": {
          "data": {
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded",
            "test": {
              "end": "2021-03-16T12:41:24Z",
              "gitCommit": "",
              "start": "2021-03-16T12:41:12Z"
            }
          },
          "id": "40ada5cd-a5ef-491b-8a6a-21d3aff6d730",
          "source": "jmeter-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:41:24.389Z",
          "type": "sh.keptn.event.test.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "111654f4-ee87-4650-99d2-dbf1b612209b"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "direct"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-dev:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "direct",
        "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
      },
      "message": "",
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "dev",
      "status": "succeeded",
      "test": {
        "teststrategy": "functional"
      }
    },
    "id": "111654f4-ee87-4650-99d2-dbf1b612209b",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:41:12.837Z",
    "type": "sh.keptn.event.test.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "direct"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-dev:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "direct",
          "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
        },
        "message": "",
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "dev",
        "status": "succeeded",
        "test": {
          "teststrategy": "functional"
        }
      },
      "id": "111654f4-ee87-4650-99d2-dbf1b612209b",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:41:12.837Z",
      "type": "sh.keptn.event.test.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "37f838ec-ebe8-440e-aeb9-f8153a6cb462",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:41:25.256Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a",
        "plainEvent": {
          "data": {
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "37f838ec-ebe8-440e-aeb9-f8153a6cb462",
          "source": "lighthouse-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:41:25.256Z",
          "type": "sh.keptn.event.evaluation.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a"
        }
      },
      {
        "traces": [],
        "data": {
          "evaluation": {
            "gitCommit": "",
            "indicatorResults": null,
            "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
            "score": 0,
            "sloFileContent": "",
            "timeEnd": "2021-03-16T12:41:24Z",
            "timeStart": "2021-03-16T12:41:12Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev"
        },
        "id": "ccee2d2d-61ae-40ab-ac60-325d55f38e69",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:30.350Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a",
        "plainEvent": {
          "data": {
            "evaluation": {
              "gitCommit": "",
              "indicatorResults": null,
              "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
              "score": 0,
              "sloFileContent": "",
              "timeEnd": "2021-03-16T12:41:24Z",
              "timeStart": "2021-03-16T12:41:12Z"
            },
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev"
          },
          "id": "ccee2d2d-61ae-40ab-ac60-325d55f38e69",
          "source": "lighthouse-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:30.350Z",
          "type": "sh.keptn.event.evaluation.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "direct"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-dev:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "direct",
        "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
      },
      "evaluation": null,
      "message": "",
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "dev",
      "status": "succeeded",
      "test": {
        "end": "2021-03-16T12:41:24Z",
        "gitCommit": "",
        "start": "2021-03-16T12:41:12Z"
      }
    },
    "id": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:41:25.249Z",
    "type": "sh.keptn.event.evaluation.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "direct"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-dev:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "direct",
          "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
        },
        "evaluation": null,
        "message": "",
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "dev",
        "status": "succeeded",
        "test": {
          "end": "2021-03-16T12:41:24Z",
          "gitCommit": "",
          "start": "2021-03-16T12:41:12Z"
        }
      },
      "id": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:41:25.249Z",
      "type": "sh.keptn.event.evaluation.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "ed677716-7d7a-44c9-ba5c-7baff85fd98e",
        "source": "dynatrace-sli-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:41:26.147Z",
        "type": "sh.keptn.event.get-sli.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "640157d9-7c33-4429-94b5-1532b6d06c33",
        "plainEvent": {
          "data": {
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "ed677716-7d7a-44c9-ba5c-7baff85fd98e",
          "source": "dynatrace-sli-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:41:26.147Z",
          "type": "sh.keptn.event.get-sli.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "640157d9-7c33-4429-94b5-1532b6d06c33"
        }
      },
      {
        "traces": [],
        "data": {
          "get-sli": {
            "end": "2021-03-16T12:41:24Z",
            "indicatorValues": [
              {
                "message": "Couldnt retrieve any SLI Results",
                "metric": "no metric",
                "success": false,
                "value": 0
              }
            ],
            "start": "2021-03-16T12:41:12Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "0480badf-c3ab-4356-a132-b6f75b757eb3",
        "source": "dynatrace-sli-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:29.931Z",
        "type": "sh.keptn.event.get-sli.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "640157d9-7c33-4429-94b5-1532b6d06c33",
        "plainEvent": {
          "data": {
            "get-sli": {
              "end": "2021-03-16T12:41:24Z",
              "indicatorValues": [
                {
                  "message": "Couldnt retrieve any SLI Results",
                  "metric": "no metric",
                  "success": false,
                  "value": 0
                }
              ],
              "start": "2021-03-16T12:41:12Z"
            },
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "0480badf-c3ab-4356-a132-b6f75b757eb3",
          "source": "dynatrace-sli-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:29.931Z",
          "type": "sh.keptn.event.get-sli.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "640157d9-7c33-4429-94b5-1532b6d06c33"
        }
      }
    ],
    "data": {
      "deployment": "direct",
      "get-sli": {
        "end": "2021-03-16T12:41:24Z",
        "sliProvider": "dynatrace",
        "start": "2021-03-16T12:41:12Z"
      },
      "project": "sockshop",
      "service": "carts",
      "stage": "dev"
    },
    "id": "640157d9-7c33-4429-94b5-1532b6d06c33",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-16T12:41:26.123Z",
    "type": "sh.keptn.event.get-sli.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "plainEvent": {
      "data": {
        "deployment": "direct",
        "get-sli": {
          "end": "2021-03-16T12:41:24Z",
          "sliProvider": "dynatrace",
          "start": "2021-03-16T12:41:12Z"
        },
        "project": "sockshop",
        "service": "carts",
        "stage": "dev"
      },
      "id": "640157d9-7c33-4429-94b5-1532b6d06c33",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-16T12:41:26.123Z",
      "type": "sh.keptn.event.get-sli.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "19ea144d-3c08-429f-8338-48388c99cefb",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:31.111Z",
        "type": "sh.keptn.event.release.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "52bab387-1738-40a5-83d5-aebbe630bb4b",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "19ea144d-3c08-429f-8338-48388c99cefb",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:31.111Z",
          "type": "sh.keptn.event.release.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "52bab387-1738-40a5-83d5-aebbe630bb4b"
        }
      },
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "message": "Finished release",
          "project": "sockshop",
          "release": {
            "gitCommit": ""
          },
          "result": "pass",
          "service": "carts",
          "stage": "dev",
          "status": "succeeded"
        },
        "id": "78694385-3ab3-4c91-94fc-9bdf2c87f519",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:31.120Z",
        "type": "sh.keptn.event.release.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "52bab387-1738-40a5-83d5-aebbe630bb4b",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "message": "Finished release",
            "project": "sockshop",
            "release": {
              "gitCommit": ""
            },
            "result": "pass",
            "service": "carts",
            "stage": "dev",
            "status": "succeeded"
          },
          "id": "78694385-3ab3-4c91-94fc-9bdf2c87f519",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:31.120Z",
          "type": "sh.keptn.event.release.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "52bab387-1738-40a5-83d5-aebbe630bb4b"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "direct"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-dev:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "direct",
        "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
      },
      "evaluation": {
        "gitCommit": "",
        "indicatorResults": null,
        "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
        "score": 0,
        "sloFileContent": "",
        "timeEnd": "2021-03-16T12:41:24Z",
        "timeStart": "2021-03-16T12:41:12Z"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "message": "",
      "project": "sockshop",
      "release": null,
      "result": "pass",
      "service": "carts",
      "stage": "dev",
      "status": "",
      "test": {
        "end": "2021-03-16T12:41:24Z",
        "gitCommit": "",
        "start": "2021-03-16T12:41:12Z"
      }
    },
    "id": "52bab387-1738-40a5-83d5-aebbe630bb4b",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:43:31.097Z",
    "type": "sh.keptn.event.release.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "direct"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-dev:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "direct",
          "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
        },
        "evaluation": {
          "gitCommit": "",
          "indicatorResults": null,
          "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
          "score": 0,
          "sloFileContent": "",
          "timeEnd": "2021-03-16T12:41:24Z",
          "timeStart": "2021-03-16T12:41:12Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "message": "",
        "project": "sockshop",
        "release": null,
        "result": "pass",
        "service": "carts",
        "stage": "dev",
        "status": "",
        "test": {
          "end": "2021-03-16T12:41:24Z",
          "gitCommit": "",
          "start": "2021-03-16T12:41:12Z"
        }
      },
      "id": "52bab387-1738-40a5-83d5-aebbe630bb4b",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:43:31.097Z",
      "type": "sh.keptn.event.release.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true,
    "started": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "2b63e4c4-9278-4231-98bc-eb9f20e5ebbd",
        "source": "shipyard-controller",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:06.468Z",
        "type": "sh.keptn.event.staging.delivery.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "shkeptnspecversion": "0.2.0",
        "triggeredid": "3c55bace-67c6-40e3-83db-cf5255481450",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "fail",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "2b63e4c4-9278-4231-98bc-eb9f20e5ebbd",
          "source": "shipyard-controller",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:06.468Z",
          "type": "sh.keptn.event.staging.delivery.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "shkeptnspecversion": "0.2.0",
          "triggeredid": "3c55bace-67c6-40e3-83db-cf5255481450"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentURIsLocal": null,
        "deploymentstrategy": ""
      },
      "project": "sockshop",
      "service": "carts",
      "stage": "staging"
    },
    "id": "3c55bace-67c6-40e3-83db-cf5255481450",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:43:33.554Z",
    "type": "sh.keptn.event.staging.delivery.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentURIsLocal": null,
          "deploymentstrategy": ""
        },
        "project": "sockshop",
        "service": "carts",
        "stage": "staging"
      },
      "id": "3c55bace-67c6-40e3-83db-cf5255481450",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:43:33.554Z",
      "type": "sh.keptn.event.staging.delivery.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "93e14eb4-7854-4fe7-8610-4a093451aa94",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:43:33.568Z",
        "type": "sh.keptn.event.deployment.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "4d8fd810-1417-496e-bc15-0e90992cadd4",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "93e14eb4-7854-4fe7-8610-4a093451aa94",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:43:33.568Z",
          "type": "sh.keptn.event.deployment.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "4d8fd810-1417-496e-bc15-0e90992cadd4"
        }
      },
      {
        "traces": [],
        "data": {
          "deployment": {
            "deploymentNames": [
              "canary"
            ],
            "deploymentURIsLocal": [
              "http://carts.sockshop-staging:80"
            ],
            "deploymentURIsPublic": [
              "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
            ],
            "deploymentstrategy": "duplicate",
            "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "message": "Successfully deployed",
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "6cd8e6bd-30c9-4523-8876-c5e2f0ce7b85",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:44:53.983Z",
        "type": "sh.keptn.event.deployment.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "4d8fd810-1417-496e-bc15-0e90992cadd4",
        "plainEvent": {
          "data": {
            "deployment": {
              "deploymentNames": [
                "canary"
              ],
              "deploymentURIsLocal": [
                "http://carts.sockshop-staging:80"
              ],
              "deploymentURIsPublic": [
                "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
              ],
              "deploymentstrategy": "duplicate",
              "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
            },
            "labels": {
              "DtCreds": "dynatrace"
            },
            "message": "Successfully deployed",
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "6cd8e6bd-30c9-4523-8876-c5e2f0ce7b85",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:44:53.983Z",
          "type": "sh.keptn.event.deployment.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "4d8fd810-1417-496e-bc15-0e90992cadd4"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "direct"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-dev:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "blue_green_service",
        "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
      },
      "evaluation": {
        "gitCommit": "",
        "indicatorResults": null,
        "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
        "score": 0,
        "sloFileContent": "",
        "timeEnd": "2021-03-16T12:41:24Z",
        "timeStart": "2021-03-16T12:41:12Z"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "message": "",
      "project": "sockshop",
      "release": {
        "gitCommit": ""
      },
      "result": "",
      "service": "carts",
      "stage": "staging",
      "status": "",
      "test": {
        "end": "2021-03-16T12:41:24Z",
        "gitCommit": "",
        "start": "2021-03-16T12:41:12Z"
      }
    },
    "id": "4d8fd810-1417-496e-bc15-0e90992cadd4",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:43:33.562Z",
    "type": "sh.keptn.event.deployment.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "direct"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-dev:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-dev.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "blue_green_service",
          "gitCommit": "42a5bdbc0cc2b0c4d258c2a74c0afbcfaaab7620"
        },
        "evaluation": {
          "gitCommit": "",
          "indicatorResults": null,
          "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
          "score": 0,
          "sloFileContent": "",
          "timeEnd": "2021-03-16T12:41:24Z",
          "timeStart": "2021-03-16T12:41:12Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "message": "",
        "project": "sockshop",
        "release": {
          "gitCommit": ""
        },
        "result": "",
        "service": "carts",
        "stage": "staging",
        "status": "",
        "test": {
          "end": "2021-03-16T12:41:24Z",
          "gitCommit": "",
          "start": "2021-03-16T12:41:12Z"
        }
      },
      "id": "4d8fd810-1417-496e-bc15-0e90992cadd4",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:43:33.562Z",
      "type": "sh.keptn.event.deployment.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "652e89da-d54e-4f2f-95fd-774446d07be3",
        "source": "jmeter-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:44:54.801Z",
        "type": "sh.keptn.event.test.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "92418029-b9fd-4dde-80be-ddcc07d03bf6",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "652e89da-d54e-4f2f-95fd-774446d07be3",
          "source": "jmeter-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:44:54.801Z",
          "type": "sh.keptn.event.test.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "92418029-b9fd-4dde-80be-ddcc07d03bf6"
        }
      },
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded",
          "test": {
            "end": "2021-03-16T12:46:58Z",
            "gitCommit": "",
            "start": "2021-03-16T12:44:54Z"
          }
        },
        "id": "c6d4ad34-fb77-4d73-b695-44e17d4e58ae",
        "source": "jmeter-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:46:58.078Z",
        "type": "sh.keptn.event.test.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "92418029-b9fd-4dde-80be-ddcc07d03bf6",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded",
            "test": {
              "end": "2021-03-16T12:46:58Z",
              "gitCommit": "",
              "start": "2021-03-16T12:44:54Z"
            }
          },
          "id": "c6d4ad34-fb77-4d73-b695-44e17d4e58ae",
          "source": "jmeter-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:46:58.078Z",
          "type": "sh.keptn.event.test.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "92418029-b9fd-4dde-80be-ddcc07d03bf6"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "canary"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-staging:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "duplicate",
        "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "message": "",
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "staging",
      "status": "succeeded",
      "test": {
        "teststrategy": "performance"
      }
    },
    "id": "92418029-b9fd-4dde-80be-ddcc07d03bf6",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:44:54.796Z",
    "type": "sh.keptn.event.test.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "canary"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-staging:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "duplicate",
          "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "message": "",
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded",
        "test": {
          "teststrategy": "performance"
        }
      },
      "id": "92418029-b9fd-4dde-80be-ddcc07d03bf6",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:44:54.796Z",
      "type": "sh.keptn.event.test.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "63134434-6f6e-4743-abb7-e2b03177da6f",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:46:58.940Z",
        "type": "sh.keptn.event.evaluation.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "63134434-6f6e-4743-abb7-e2b03177da6f",
          "source": "lighthouse-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:46:58.940Z",
          "type": "sh.keptn.event.evaluation.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89"
        }
      },
      {
        "traces": [],
        "data": {
          "evaluation": {
            "comparedEvents": [
              "8ad74029-5ceb-4680-b697-d1701078faff"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "This is my very very very very very very very very very very very very long displayName",
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 11.663352564204615,
                    "violated": true
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 0.5,
                "status": "warning",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 331.4225933914111
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "fail",
            "score": 50,
            "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-16T12:46:58Z",
            "timeStart": "2021-03-16T12:44:54Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "843a4a6b-4dec-4328-8579-6af3d1a6c7cb",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:05.732Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
        "plainEvent": {
          "data": {
            "evaluation": {
              "comparedEvents": [
                "8ad74029-5ceb-4680-b697-d1701078faff"
              ],
              "gitCommit": "",
              "indicatorResults": [
                {
                  "displayName": "This is my very very very very very very very very very very very very long displayName",
                  "keySli": false,
                  "passTargets": [
                    {
                      "criteria": "<=+10%",
                      "targetValue": 11.663352564204615,
                      "violated": true
                    },
                    {
                      "criteria": "<600",
                      "targetValue": 600,
                      "violated": false
                    }
                  ],
                  "score": 0.5,
                  "status": "warning",
                  "value": {
                    "metric": "response_time_p95",
                    "success": true,
                    "value": 331.4225933914111
                  },
                  "warningTargets": [
                    {
                      "criteria": "<=800",
                      "targetValue": 800,
                      "violated": false
                    }
                  ]
                }
              ],
              "result": "fail",
              "score": 50,
              "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
              "timeEnd": "2021-03-16T12:46:58Z",
              "timeStart": "2021-03-16T12:44:54Z"
            },
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "fail",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "843a4a6b-4dec-4328-8579-6af3d1a6c7cb",
          "source": "lighthouse-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:05.732Z",
          "type": "sh.keptn.event.evaluation.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "canary"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-staging:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "duplicate",
        "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
      },
      "evaluation": null,
      "labels": {
        "DtCreds": "dynatrace"
      },
      "message": "",
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "staging",
      "status": "succeeded",
      "test": {
        "end": "2021-03-16T12:46:58Z",
        "gitCommit": "",
        "start": "2021-03-16T12:44:54Z"
      }
    },
    "id": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:46:58.923Z",
    "type": "sh.keptn.event.evaluation.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "canary"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-staging:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "duplicate",
          "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
        },
        "evaluation": null,
        "labels": {
          "DtCreds": "dynatrace"
        },
        "message": "",
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded",
        "test": {
          "end": "2021-03-16T12:46:58Z",
          "gitCommit": "",
          "start": "2021-03-16T12:44:54Z"
        }
      },
      "id": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:46:58.923Z",
      "type": "sh.keptn.event.evaluation.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true,
    "icon": "traffic-light",
    "label": "evaluation"
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "06b0b784-70bf-494d-a4be-8339cac1d273",
        "source": "dynatrace-sli-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:46:59.663Z",
        "type": "sh.keptn.event.get-sli.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "f4a6c143-d433-4451-b2b6-b67686775e11",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "06b0b784-70bf-494d-a4be-8339cac1d273",
          "source": "dynatrace-sli-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:46:59.663Z",
          "type": "sh.keptn.event.get-sli.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "f4a6c143-d433-4451-b2b6-b67686775e11"
        }
      },
      {
        "traces": [],
        "data": {
          "get-sli": {
            "end": "2021-03-16T12:46:58Z",
            "indicatorValues": [
              {
                "metric": "response_time_p95",
                "success": true,
                "value": 331.4225933914111
              }
            ],
            "start": "2021-03-16T12:44:54Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "11618479-0c21-447d-aef2-03fef0a5ee04",
        "source": "dynatrace-sli-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:04.871Z",
        "type": "sh.keptn.event.get-sli.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "f4a6c143-d433-4451-b2b6-b67686775e11",
        "plainEvent": {
          "data": {
            "get-sli": {
              "end": "2021-03-16T12:46:58Z",
              "indicatorValues": [
                {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 331.4225933914111
                }
              ],
              "start": "2021-03-16T12:44:54Z"
            },
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "result": "pass",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "11618479-0c21-447d-aef2-03fef0a5ee04",
          "source": "dynatrace-sli-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:04.871Z",
          "type": "sh.keptn.event.get-sli.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "f4a6c143-d433-4451-b2b6-b67686775e11"
        }
      }
    ],
    "data": {
      "deployment": "canary",
      "get-sli": {
        "end": "2021-03-16T12:46:58Z",
        "indicators": [
          "response_time_p95"
        ],
        "sliProvider": "dynatrace",
        "start": "2021-03-16T12:44:54Z"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "project": "sockshop",
      "service": "carts",
      "stage": "staging"
    },
    "id": "f4a6c143-d433-4451-b2b6-b67686775e11",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-16T12:46:59.660Z",
    "type": "sh.keptn.event.get-sli.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "plainEvent": {
      "data": {
        "deployment": "canary",
        "get-sli": {
          "end": "2021-03-16T12:46:58Z",
          "indicators": [
            "response_time_p95"
          ],
          "sliProvider": "dynatrace",
          "start": "2021-03-16T12:44:54Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "service": "carts",
        "stage": "staging"
      },
      "id": "f4a6c143-d433-4451-b2b6-b67686775e11",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-16T12:46:59.660Z",
      "type": "sh.keptn.event.get-sli.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "message": "Finished rollback",
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "9589b4d7-fa63-4dfe-a15e-a9528ed6e76c",
        "source": "shipyard-controller",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:16.571Z",
        "type": "sh.keptn.event.staging.rollback.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "shkeptnspecversion": "0.2.0",
        "triggeredid": "247b2d4b-a178-4dc6-aaef-e60b2ea2eda7",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "message": "Finished rollback",
            "project": "sockshop",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "9589b4d7-fa63-4dfe-a15e-a9528ed6e76c",
          "source": "shipyard-controller",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:16.571Z",
          "type": "sh.keptn.event.staging.rollback.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "shkeptnspecversion": "0.2.0",
          "triggeredid": "247b2d4b-a178-4dc6-aaef-e60b2ea2eda7"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentURIsLocal": null,
        "deploymentstrategy": ""
      },
      "project": "sockshop",
      "service": "carts",
      "stage": "staging"
    },
    "id": "247b2d4b-a178-4dc6-aaef-e60b2ea2eda7",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:48:06.469Z",
    "type": "sh.keptn.event.staging.rollback.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentURIsLocal": null,
          "deploymentstrategy": ""
        },
        "project": "sockshop",
        "service": "carts",
        "stage": "staging"
      },
      "id": "247b2d4b-a178-4dc6-aaef-e60b2ea2eda7",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:48:06.469Z",
      "type": "sh.keptn.event.staging.rollback.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true
  },
  {
    "traces": [
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "c0fc1417-667b-4ca2-aaca-4ef9ec126e00",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:06.487Z",
        "type": "sh.keptn.event.rollback.started",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "c111f095-74dd-4a33-bfcb-b3d6440cd239",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "project": "sockshop",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "c0fc1417-667b-4ca2-aaca-4ef9ec126e00",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:06.487Z",
          "type": "sh.keptn.event.rollback.started",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "c111f095-74dd-4a33-bfcb-b3d6440cd239"
        }
      },
      {
        "traces": [],
        "data": {
          "labels": {
            "DtCreds": "dynatrace"
          },
          "message": "Finished rollback",
          "project": "sockshop",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "f24967fb-786d-48dd-bb2f-956923e89788",
        "source": "helm-service",
        "specversion": "1.0",
        "time": "2021-03-16T12:48:16.124Z",
        "type": "sh.keptn.event.rollback.finished",
        "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
        "triggeredid": "c111f095-74dd-4a33-bfcb-b3d6440cd239",
        "plainEvent": {
          "data": {
            "labels": {
              "DtCreds": "dynatrace"
            },
            "message": "Finished rollback",
            "project": "sockshop",
            "service": "carts",
            "stage": "staging",
            "status": "succeeded"
          },
          "id": "f24967fb-786d-48dd-bb2f-956923e89788",
          "source": "helm-service",
          "specversion": "1.0",
          "time": "2021-03-16T12:48:16.124Z",
          "type": "sh.keptn.event.rollback.finished",
          "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
          "triggeredid": "c111f095-74dd-4a33-bfcb-b3d6440cd239"
        }
      }
    ],
    "data": {
      "configurationChange": {
        "values": {
          "image": "docker.io/keptnexamples/carts:0.12.3"
        }
      },
      "deployment": {
        "deploymentNames": [
          "canary"
        ],
        "deploymentURIsLocal": [
          "http://carts.sockshop-staging:80"
        ],
        "deploymentURIsPublic": [
          "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
        ],
        "deploymentstrategy": "duplicate",
        "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
      },
      "evaluation": {
        "comparedEvents": [
          "8ad74029-5ceb-4680-b697-d1701078faff"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 11.663352564204615,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 600,
                "violated": false
              }
            ],
            "score": 0.5,
            "status": "warning",
            "value": {
              "metric": "response_time_p95",
              "success": true,
              "value": 331.4225933914111
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 800,
                "violated": false
              }
            ]
          }
        ],
        "result": "fail",
        "score": 50,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-16T12:46:58Z",
        "timeStart": "2021-03-16T12:44:54Z"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "message": "",
      "project": "sockshop",
      "result": "",
      "rollback": null,
      "service": "carts",
      "stage": "staging",
      "status": "",
      "test": {
        "end": "2021-03-16T12:46:58Z",
        "gitCommit": "",
        "start": "2021-03-16T12:44:54Z"
      }
    },
    "id": "c111f095-74dd-4a33-bfcb-b3d6440cd239",
    "source": "shipyard-controller",
    "specversion": "1.0",
    "time": "2021-03-16T12:48:06.481Z",
    "type": "sh.keptn.event.rollback.triggered",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "shkeptnspecversion": "0.2.0",
    "plainEvent": {
      "data": {
        "configurationChange": {
          "values": {
            "image": "docker.io/keptnexamples/carts:0.12.3"
          }
        },
        "deployment": {
          "deploymentNames": [
            "canary"
          ],
          "deploymentURIsLocal": [
            "http://carts.sockshop-staging:80"
          ],
          "deploymentURIsPublic": [
            "http://carts.sockshop-staging.34.72.118.255.nip.io:80"
          ],
          "deploymentstrategy": "duplicate",
          "gitCommit": "cbb81a2a64047bef051f39d9153d14d524465442"
        },
        "evaluation": {
          "comparedEvents": [
            "8ad74029-5ceb-4680-b697-d1701078faff"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 11.663352564204615,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 0.5,
              "status": "warning",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 331.4225933914111
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "fail",
          "score": 50,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-16T12:46:58Z",
          "timeStart": "2021-03-16T12:44:54Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "message": "",
        "project": "sockshop",
        "result": "",
        "rollback": null,
        "service": "carts",
        "stage": "staging",
        "status": "",
        "test": {
          "end": "2021-03-16T12:46:58Z",
          "gitCommit": "",
          "start": "2021-03-16T12:44:54Z"
        }
      },
      "id": "c111f095-74dd-4a33-bfcb-b3d6440cd239",
      "source": "shipyard-controller",
      "specversion": "1.0",
      "time": "2021-03-16T12:48:06.481Z",
      "type": "sh.keptn.event.rollback.triggered",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "shkeptnspecversion": "0.2.0"
    },
    "finished": true,
    "started": true,
    "label": "rollback"
  }
];

const Traces = traceData.map(trace => {
  function tracesMapper(trace) {
    trace.traces.forEach(t => {
      tracesMapper(t);
    });
    trace.traces = Trace.traceMapper(trace.traces);
  }

  tracesMapper(trace);
  return Trace.fromJSON(trace);
});

export {Traces};
