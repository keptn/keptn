import {Project} from '../../_models/project';
import {Service} from '../../_models/service';
import {Stage} from '../../_models/stage';

// tslint:disable-next-line:no-any
const projectsData: any = [
  {
    "creationDate": "1614603785739925270",
    "gitRemoteURI": "https://github.com/Kirdock/keptn-dynatrace",
    "gitUser": "Kirdock",
    "projectName": "sockshop",
    "shipyard": "apiVersion: spec.keptn.sh/0.2.0\nkind: Shipyard\nmetadata:\n  name: shipyard-sockshop\nspec:\n  stages:\n  - name: dev\n    sequences:\n    - name: delivery\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: test\n        properties:\n          teststrategy: functional\n      - name: evaluation\n        properties: null\n      - name: release\n        properties: null\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n  - name: staging\n    sequences:\n    - name: delivery\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: blue_green_service\n      - name: test\n        properties:\n          teststrategy: performance\n      - name: evaluation\n        properties: null\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: dev.delivery.finished\n        selector:\n          match: null\n    - name: rollback\n      tasks:\n      - name: rollback\n        properties: null\n      triggeredOn:\n      - event: staging.delivery.finished\n        selector:\n          match:\n            result: fail\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: dev.delivery-direct.finished\n        selector:\n          match: null\n  - name: production\n    sequences:\n    - name: delivery\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: blue_green_service\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: staging.delivery.finished\n        selector:\n          match: null\n    - name: rollback\n      tasks:\n      - name: rollback\n        properties: null\n      triggeredOn:\n      - event: production.delivery.finished\n        selector:\n          match:\n            result: fail\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: staging.delivery-direct.finished\n        selector:\n          match: null\n",
    "shipyardVersion": "spec.keptn.sh/0.2.0",
    "stages": [
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614603919173353566",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "202c36a2-30b5-4dac-ae95-bb3cef1d6842",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545130257293411"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "4fd46017-8103-470b-9d52-1567ee41c684",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545122587658908"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "3be85957-bd5a-425f-ab68-3a220ff545d6",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545122571635732"
              },
              "sh.keptn.event.dev.delivery-direct.finished": {
                "eventId": "d9e23c66-f371-4b23-9b19-11c20b663285",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545132058591548"
              },
              "sh.keptn.event.dev.delivery-direct.triggered": {
                "eventId": "2f2e2a12-4c4f-4611-b20b-d2d21000a718",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545122108717957"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "f80d06ab-84c8-443e-b58a-156d9168d183",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545130983519777"
              },
              "sh.keptn.event.release.started": {
                "eventId": "e7475922-7487-4c05-87f2-d7ce1c9efc42",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545130971097045"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "ed27cf71-f6e0-40fe-af87-ee3bead34cfc",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545130961599182"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "dev"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614683118180068607",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.3",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "fde6a8c8-a2f0-4bbb-9abf-1d36208d6ed9",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898472045437872"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "fdc4c95c-163c-4c3b-94a2-6e57e314fca5",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898394237735710"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "a030c58a-9446-45f6-a6cd-ba00e9fc3729",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898394225809044"
              },
              "sh.keptn.event.dev.delivery.finished": {
                "eventId": "afeaf817-6600-4925-a89b-1086b9baf18b",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898613559503742"
              },
              "sh.keptn.event.dev.delivery.triggered": {
                "eventId": "cd9015ff-fad0-43d3-abc0-d50455aa8a82",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898393627257537"
              },
              "sh.keptn.event.evaluation.finished": {
                "eventId": "ccee2d2d-61ae-40ab-ac60-325d55f38e69",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898610375438224"
              },
              "sh.keptn.event.evaluation.started": {
                "eventId": "37f838ec-ebe8-440e-aeb9-f8153a6cb462",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898485292418963"
              },
              "sh.keptn.event.evaluation.triggered": {
                "eventId": "e8f12220-b0f7-4e2f-898a-b6b7e699f12a",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898485255518680"
              },
              "sh.keptn.event.get-sli.finished": {
                "eventId": "0480badf-c3ab-4356-a132-b6f75b757eb3",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898609964450887"
              },
              "sh.keptn.event.get-sli.started": {
                "eventId": "ed677716-7d7a-44c9-ba5c-7baff85fd98e",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898486208132629"
              },
              "sh.keptn.event.get-sli.triggered": {
                "eventId": "640157d9-7c33-4429-94b5-1532b6d06c33",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898486134484401"
              },
              "sh.keptn.event.problem.open": {
                "eventId": "31a4c5a4-97a8-421d-b8d5-0c33087824cc",
                "keptnContext": "35393135-3337-4034-b036-313832383631",
                "time": "1614889276800269532"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "78694385-3ab3-4c91-94fc-9bdf2c87f519",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898611130968020"
              },
              "sh.keptn.event.release.started": {
                "eventId": "19ea144d-3c08-429f-8338-48388c99cefb",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898611125181870"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "52bab387-1738-40a5-83d5-aebbe630bb4b",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898611112643920"
              },
              "sh.keptn.event.remediation.finished": {
                "eventId": "db32c59f-aefa-46f5-a11a-22d9f79469da",
                "keptnContext": "35393135-3337-4034-b036-313832383631",
                "time": "1614889277272328434"
              },
              "sh.keptn.event.test.finished": {
                "eventId": "40ada5cd-a5ef-491b-8a6a-21d3aff6d730",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898484396438377"
              },
              "sh.keptn.event.test.started": {
                "eventId": "19dad0eb-166e-4bc9-b972-9488c7df6e9a",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898472866520338"
              },
              "sh.keptn.event.test.triggered": {
                "eventId": "111654f4-ee87-4650-99d2-dbf1b612209b",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898472844048602"
              },
              "sh.keptn.events.problem": {
                "eventId": "c31cfbe0-9bff-47e3-bb9f-8039cf48a127",
                "keptnContext": "35393135-3337-4034-b036-313832383631",
                "time": "1614923573857764232"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "dev"
          }
        ],
        "stageName": "dev"
      },
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614603922172957849",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "70cdc742-bdb8-4c10-ada1-460dd6152688",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545140633525288"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "e8667067-27b1-485c-a589-0a64f5e2e029",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545132083550050"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "2e2e88ed-57e7-4a47-9f53-7090a70a94dc",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545132076584253"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "21d996e2-9c49-4a82-b3d8-c0b8df43bcd7",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545141381207222"
              },
              "sh.keptn.event.release.started": {
                "eventId": "71002a41-98a7-4cee-8a43-fba175cbcf25",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545141376856199"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "ae23b030-a425-4ac7-817d-cb517e3ec703",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545141369935465"
              },
              "sh.keptn.event.staging.delivery-direct.finished": {
                "eventId": "4955b152-0574-4b64-b601-6846cdc0f921",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545143728033338"
              },
              "sh.keptn.event.staging.delivery-direct.triggered": {
                "eventId": "a0b2b3fe-63de-4cde-80a8-46526aa462cd",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545132066189728"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "staging"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614683121241067984",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.3",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "6cd8e6bd-30c9-4523-8876-c5e2f0ce7b85",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898693990663363"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "93e14eb4-7854-4fe7-8610-4a093451aa94",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898613573692102"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "4d8fd810-1417-496e-bc15-0e90992cadd4",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898613569486946"
              },
              "sh.keptn.event.evaluation.finished": {
                "eventId": "843a4a6b-4dec-4328-8579-6af3d1a6c7cb",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898885741220213"
              },
              "sh.keptn.event.evaluation.started": {
                "eventId": "63134434-6f6e-4743-abb7-e2b03177da6f",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898818957131635"
              },
              "sh.keptn.event.evaluation.triggered": {
                "eventId": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898818943282666"
              },
              "sh.keptn.event.get-sli.finished": {
                "eventId": "11618479-0c21-447d-aef2-03fef0a5ee04",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898884878751901"
              },
              "sh.keptn.event.get-sli.started": {
                "eventId": "06b0b784-70bf-494d-a4be-8339cac1d273",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898819680818518"
              },
              "sh.keptn.event.get-sli.triggered": {
                "eventId": "f4a6c143-d433-4451-b2b6-b67686775e11",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898819665203549"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "22b7b774-9128-45e4-a8e8-2c4528e967d0",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474626169359095"
              },
              "sh.keptn.event.release.started": {
                "eventId": "f3af31c8-113d-464b-ac67-72d887e2b7b4",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474541420933551"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "779a0bf9-8818-4e45-a769-eff6de522072",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474541399138902"
              },
              "sh.keptn.event.rollback.finished": {
                "eventId": "f24967fb-786d-48dd-bb2f-956923e89788",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898896131859300"
              },
              "sh.keptn.event.rollback.started": {
                "eventId": "c0fc1417-667b-4ca2-aaca-4ef9ec126e00",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898886515677581"
              },
              "sh.keptn.event.rollback.triggered": {
                "eventId": "c111f095-74dd-4a33-bfcb-b3d6440cd239",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898886487700557"
              },
              "sh.keptn.event.staging.delivery.finished": {
                "eventId": "2b63e4c4-9278-4231-98bc-eb9f20e5ebbd",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898886474698013"
              },
              "sh.keptn.event.staging.delivery.triggered": {
                "eventId": "3c55bace-67c6-40e3-83db-cf5255481450",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898613565062073"
              },
              "sh.keptn.event.staging.evaluation.finished": {
                "eventId": "c022c751-f0ef-4298-8072-efb05273cc78",
                "keptnContext": "8e452c9b-b8ac-4310-9f8e-bb62acfe4ca4",
                "time": "1615543816735591333"
              },
              "sh.keptn.event.staging.evaluation.triggered": {
                "eventId": "4b868d6c-e1fb-46e1-ba0e-29f416bbe88a",
                "keptnContext": "8e452c9b-b8ac-4310-9f8e-bb62acfe4ca4",
                "time": "1615543813048505365"
              },
              "sh.keptn.event.staging.rollback.finished": {
                "eventId": "9589b4d7-fa63-4dfe-a15e-a9528ed6e76c",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898896576530346"
              },
              "sh.keptn.event.staging.rollback.triggered": {
                "eventId": "247b2d4b-a178-4dc6-aaef-e60b2ea2eda7",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898886479794324"
              },
              "sh.keptn.event.test.finished": {
                "eventId": "c6d4ad34-fb77-4d73-b695-44e17d4e58ae",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898818099504962"
              },
              "sh.keptn.event.test.started": {
                "eventId": "652e89da-d54e-4f2f-95fd-774446d07be3",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898694823527680"
              },
              "sh.keptn.event.test.triggered": {
                "eventId": "92418029-b9fd-4dde-80be-ddcc07d03bf6",
                "keptnContext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
                "time": "1615898694804880166"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "staging"
          }
        ],
        "stageName": "staging"
      },
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614603925140444462",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "6af68341-0189-4f14-9a17-ebaee102ae20",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545151267591359"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "b83c3145-52b5-434d-974f-3224b495620b",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545143781219481"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "7ebcec2b-571d-484e-a2f9-baeed495e847",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545143777266660"
              },
              "sh.keptn.event.production.delivery-direct.finished": {
                "eventId": "b82eb79a-2b2f-46c9-bf22-cbd88f072d23",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545154372228594"
              },
              "sh.keptn.event.production.delivery-direct.triggered": {
                "eventId": "898e3a3a-bb3a-452c-ba7a-6dc95c94c6f7",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545143736044107"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "4aedbd21-2c90-43bd-8680-9f7b560495ec",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545152024553783"
              },
              "sh.keptn.event.release.started": {
                "eventId": "c5f42ecc-62f7-4ad4-a18b-e3931e6a2c28",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545152033261026"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "b966ed4a-d618-445b-b39a-28f471f19ad1",
                "keptnContext": "9bff030c-1647-41ba-90fc-608577cbeb3a",
                "time": "1615545152016404557"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "production"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614683124274548599",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.3",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "304e1f24-4d11-4e0e-9c7f-91b96f5d0be4",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474722389854904"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "debf5976-28d4-4e99-8836-fa977777f612",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474626580191638"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "6ab5b350-1fd3-42bb-b990-09afaec5a59f",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474626571723822"
              },
              "sh.keptn.event.production.delivery.finished": {
                "eventId": "e9bed2dd-2c79-4834-9469-4277e59fc264",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474809222114700"
              },
              "sh.keptn.event.production.delivery.triggered": {
                "eventId": "de98c80e-da89-4852-83a5-701844b43914",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474626562627333"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "7c9a8b80-524f-4520-8f5d-7928d7372998",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474808774861076"
              },
              "sh.keptn.event.release.started": {
                "eventId": "fe49d542-6302-4198-a592-79b30784eb30",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474723127771535"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "f03ba9a8-b720-445c-ab67-5f161ca17fa7",
                "keptnContext": "3061342b-03eb-485c-8586-de1b50c16451",
                "time": "1615474723121107065"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "production"
          }
        ],
        "stageName": "production"
      }
    ]
  },
  {
    "creationDate": "1614610006846830117",
    "gitRemoteURI": "https://github.com/Kirdock/keptnTestApprove",
    "gitUser": "Kirdock",
    "projectName": "sockshop-approve",
    "shipyard": "apiVersion: spec.keptn.sh/0.2.0\nkind: Shipyard\nmetadata:\n  name: shipyard-sockshop\nspec:\n  stages:\n  - name: dev\n    sequences:\n    - name: delivery\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: test\n        properties:\n          teststrategy: functional\n      - name: evaluation\n        properties: null\n      - name: release\n        properties: null\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n  - name: staging\n    sequences:\n    - name: delivery\n      tasks:\n      - name: approval\n        properties:\n          pass: manual\n          warning: manual\n      - name: deployment\n        properties:\n          deploymentstrategy: blue_green_service\n      - name: test\n        properties:\n          teststrategy: performance\n      - name: evaluation\n        properties: null\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: dev.delivery.finished\n        selector:\n          match: null\n    - name: rollback\n      tasks:\n      - name: rollback\n        properties: null\n      triggeredOn:\n      - event: staging.delivery.finished\n        selector:\n          match:\n            result: fail\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: dev.delivery-direct.finished\n        selector:\n          match: null\n  - name: production\n    sequences:\n    - name: delivery\n      tasks:\n      - name: approval\n        properties:\n          pass: manual\n          warning: manual\n      - name: deployment\n        properties:\n          deploymentstrategy: blue_green_service\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: staging.delivery.finished\n        selector:\n          match: null\n    - name: rollback\n      tasks:\n      - name: rollback\n        properties: null\n      triggeredOn:\n      - event: production.delivery.finished\n        selector:\n          match:\n            result: fail\n    - name: delivery-direct\n      tasks:\n      - name: deployment\n        properties:\n          deploymentstrategy: direct\n      - name: release\n        properties: null\n      triggeredOn:\n      - event: staging.delivery-direct.finished\n        selector:\n          match: null\n",
    "shipyardVersion": "spec.keptn.sh/0.2.0",
    "stages": [
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610067247515271",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "661009e8-bb55-4c41-be72-0375d34482c6",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683234888685129"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "097cbd7e-0527-4e89-87d5-64178374a899",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683154888980985"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "dfa8247b-1335-412b-a02b-420a2e773a57",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683154852092079"
              },
              "sh.keptn.event.dev.delivery.finished": {
                "eventId": "af78edbf-c05d-44c2-ac40-9e421c033568",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683252205585870"
              },
              "sh.keptn.event.dev.delivery.triggered": {
                "eventId": "dab515e2-809e-4543-b6ec-e62d585a35d1",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683154435958109"
              },
              "sh.keptn.event.evaluation.finished": {
                "eventId": "6896383a-e5d5-4ece-a6ae-86b79bc1ebbe",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683250033478466"
              },
              "sh.keptn.event.evaluation.started": {
                "eventId": "d50d383b-dd9d-448d-a58c-76d79d4dd30d",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683249304937254"
              },
              "sh.keptn.event.evaluation.triggered": {
                "eventId": "becbce99-caff-4e77-8636-174a700cb33d",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683249265683243"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "7730f83b-4719-40ee-b9c5-0870d3930732",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683251115775790"
              },
              "sh.keptn.event.release.started": {
                "eventId": "3f30c90e-ffa6-4f3f-a532-2555a931e92c",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683251109567402"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "02af0943-460d-4d64-a6e7-a484622a06b5",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683251103211502"
              },
              "sh.keptn.event.test.finished": {
                "eventId": "9b6b1d69-3a95-4485-847b-45d082db1ee5",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683248523488576"
              },
              "sh.keptn.event.test.started": {
                "eventId": "85a9a41a-5a91-46ed-81ad-ad4a0e346156",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683235711846947"
              },
              "sh.keptn.event.test.triggered": {
                "eventId": "b400e48d-9614-414f-8c18-366845c9a719",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683235676615049"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "dev"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610137361542762",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "a476bcf0-51d6-4e9c-8b2a-8b3edc52360a",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610198943708443"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "190c544a-a9d2-4d1c-ba31-0e9d5fc0dfee",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610164789559815"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "5665a33b-1795-4099-a763-46891fc116a0",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610164784117392"
              },
              "sh.keptn.event.dev.delivery-direct.finished": {
                "eventId": "5cbd388d-8196-4c63-8697-0cd641260d5b",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610202057754648"
              },
              "sh.keptn.event.dev.delivery-direct.triggered": {
                "eventId": "413f8ca0-bcde-4b98-abaa-d7e6ea9a9d5d",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610164369652732"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "4d589bb7-95a7-419f-8402-bd56f6cb97aa",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610199632019244"
              },
              "sh.keptn.event.release.started": {
                "eventId": "8a266d56-a6eb-430d-9315-cc459f1f5882",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610199635207083"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "f935b264-55d8-4b33-ba87-e9c5803e1796",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610199628683416"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "dev"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478044097626854",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "b5e210a0-2d1e-43de-820f-d33dd2ab0fc6",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479365266582820"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "e0b6b5fa-a10a-4d40-8ff5-f184b8cc70f8",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479290340192938"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "fc5ccd2f-7afc-489f-aa4a-e81109550c52",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479290325578724"
              },
              "sh.keptn.event.dev.delivery.finished": {
                "eventId": "b595a781-2eb1-427c-b72c-3101b84d5c88",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479382768764231"
              },
              "sh.keptn.event.dev.delivery.triggered": {
                "eventId": "fb787c3f-5f22-4038-973d-792e42aa67c4",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479289967989041"
              },
              "sh.keptn.event.evaluation.finished": {
                "eventId": "7a81718d-18cb-4284-9ab9-4f6b9084db0b",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479379329259291"
              },
              "sh.keptn.event.evaluation.started": {
                "eventId": "4260fb39-7b7c-4418-bf03-5402ba02e5f9",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479378559775741"
              },
              "sh.keptn.event.evaluation.triggered": {
                "eventId": "f7a07fef-66d8-4107-8864-a51f07492ea9",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479378541534147"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "17ba73aa-97f1-45b6-8544-978ded8e3b7d",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479380405979184"
              },
              "sh.keptn.event.release.started": {
                "eventId": "f087a5f3-0670-49ff-910e-f73a59b947aa",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479380411304426"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "6dde7544-1715-4c6a-8e81-aa8131df1dd3",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479380393383215"
              },
              "sh.keptn.event.test.finished": {
                "eventId": "cd1aa49c-0ffe-4f16-bb43-5632a50aaec5",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479377830518216"
              },
              "sh.keptn.event.test.started": {
                "eventId": "f566e51e-9465-45f1-a6cf-68c02d216d88",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479366053496890"
              },
              "sh.keptn.event.test.triggered": {
                "eventId": "0a35eec3-2023-46ad-a915-6ea7a27c977d",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479366034572304"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-s1",
            "stage": "dev"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478098597874154",
            "openRemediations": null,
            "serviceName": "carts-s2",
            "stage": "dev"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478144985066621",
            "openRemediations": null,
            "serviceName": "carts-s3",
            "stage": "dev"
          }
        ],
        "stageName": "dev"
      },
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610070364938604",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.2",
            "lastEventTypes": {
              "sh.keptn.event.approval.finished": {
                "eventId": "e5b917bd-8ada-4e5a-888e-0cc476484308",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683517555711525"
              },
              "sh.keptn.event.approval.started": {
                "eventId": "2f5faf00-6d0a-41f4-9ae9-ec9ac70e21f6",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683517247519004"
              },
              "sh.keptn.event.approval.triggered": {
                "eventId": "3a2d8e70-f86b-4b4d-b25e-aab8b7e5d36c",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683252271655959"
              },
              "sh.keptn.event.deployment.finished": {
                "eventId": "64c82587-e4e5-4d8d-a41a-73a8411a17b4",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683599755741951"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "db7f5bde-2de1-4204-8480-1817908b6c06",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683517991195984"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "e85cfbdb-8af7-4bef-8a47-466aa645ca86",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683517975726918"
              },
              "sh.keptn.event.evaluation.finished": {
                "eventId": "d2825612-c1c7-4aec-a5f9-6102dce2886b",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684133993047985"
              },
              "sh.keptn.event.evaluation.started": {
                "eventId": "b42d4689-a491-448a-be68-01214e77c230",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684133285362219"
              },
              "sh.keptn.event.evaluation.triggered": {
                "eventId": "c5900b89-9ab0-4352-a327-68a2588e153e",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684133272736635"
              },
              "sh.keptn.event.problem.open": {
                "eventId": "ec232311-80e8-4e71-a539-efa59c1e6fa9",
                "keptnContext": "2d343637-3330-4432-b533-353137343834",
                "time": "1615505656551327116"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "0cbdacb0-b77f-4529-a343-1df55ead1167",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684219831151446"
              },
              "sh.keptn.event.release.started": {
                "eventId": "c70e22b6-230f-4019-b0ac-65a89194d28e",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684135155064268"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "86e9ce5a-0c99-4e91-b9d5-ed775d616852",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684135140084209"
              },
              "sh.keptn.event.remediation.finished": {
                "eventId": "da32ad0d-4c50-4548-aaca-5c73e89648ac",
                "keptnContext": "2d343637-3330-4432-b533-353137343834",
                "time": "1615505657171483154"
              },
              "sh.keptn.event.rollback.finished": {
                "eventId": "4370ffd2-8b2e-431f-8d09-6b6908a7c3f2",
                "keptnContext": "8c152dc3-6777-4f08-85b9-5c05dc6379aa",
                "time": "1615819736764848293"
              },
              "sh.keptn.event.rollback.started": {
                "eventId": "8b20ec78-9acb-4789-bdbb-cd55b080b603",
                "keptnContext": "8c152dc3-6777-4f08-85b9-5c05dc6379aa",
                "time": "1615819727972636040"
              },
              "sh.keptn.event.rollback.triggered": {
                "eventId": "d1e2ccde-78f4-4673-b84c-84f005222377",
                "keptnContext": "8c152dc3-6777-4f08-85b9-5c05dc6379aa",
                "time": "1615819727913817349"
              },
              "sh.keptn.event.staging.delivery.finished": {
                "eventId": "e0f92c7e-2ccc-472e-a529-76a826eed4e6",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684220215799371"
              },
              "sh.keptn.event.staging.delivery.triggered": {
                "eventId": "e2d30b53-0774-4eb6-a5a7-8a78a2136f68",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683252211203416"
              },
              "sh.keptn.event.staging.rollback.finished": {
                "eventId": "d40b5f82-c8bd-4646-bb37-b3101a78764c",
                "keptnContext": "8c152dc3-6777-4f08-85b9-5c05dc6379aa",
                "time": "1615819737130172457"
              },
              "sh.keptn.event.staging.rollback.triggered": {
                "eventId": "d08f1973-c670-4da2-9234-8bcc94806364",
                "keptnContext": "8c152dc3-6777-4f08-85b9-5c05dc6379aa",
                "time": "1615819727905255180"
              },
              "sh.keptn.event.test.finished": {
                "eventId": "58139559-7afd-47cd-b633-fd92f4b9c8e0",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684132423921051"
              },
              "sh.keptn.event.test.started": {
                "eventId": "4a1ef311-0248-4ee7-ac53-57c4f7f6a036",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683600568734850"
              },
              "sh.keptn.event.test.triggered": {
                "eventId": "273ad612-5f51-4e7b-8219-256ba217a478",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616683600552439564"
              },
              "sh.keptn.events.problem": {
                "eventId": "1e3cf156-8269-4a68-aef5-de19b73b125d",
                "keptnContext": "2d343637-3330-4432-b533-353137343834",
                "time": "1615505706529080764"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "staging"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610140488691226",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "82bfbacd-b121-4bb2-a76f-37450f363790",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610228427679155"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "6210dc0e-488d-42ad-90da-7785f6d1349a",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610202096307951"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "81445372-061d-4bd6-9ea5-eebe16ee2954",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610202077772196"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "ff022bc2-59f7-423c-9ea7-ee00ad424ce9",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610229092997094"
              },
              "sh.keptn.event.release.started": {
                "eventId": "31687450-88dd-4def-bb9b-f336cad1360c",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610229089525857"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "fc77aeca-e62a-4972-850e-c32e8b59e0e2",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610229085309626"
              },
              "sh.keptn.event.staging.delivery-direct.finished": {
                "eventId": "1aad21b6-e27d-452d-9b27-d7dbe09e7bc3",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610230076578212"
              },
              "sh.keptn.event.staging.delivery-direct.triggered": {
                "eventId": "c493d51e-e7bb-4d09-b695-7274a2424778",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610202071163140"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "staging"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478046642275369",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.2",
            "lastEventTypes": {
              "sh.keptn.event.approval.finished": {
                "eventId": "83daad83-ac32-4b0d-ae20-fd4dd5f92b3e",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615537889940450225"
              },
              "sh.keptn.event.approval.started": {
                "eventId": "64a6403b-06fa-45dd-9d78-01ad8d3af603",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615537889684588384"
              },
              "sh.keptn.event.approval.triggered": {
                "eventId": "78f4eb93-09a0-4533-b306-e59c69041658",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479382781225754"
              },
              "sh.keptn.event.deployment.finished": {
                "eventId": "10073c1e-bcdb-45b5-841a-cb3438a26c49",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538150192518173"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "d5624c15-7f39-4c36-b76e-e601492845a1",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615537890436770975"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "dd3c6592-bbdd-4359-9a23-97475078a747",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615537890415216391"
              },
              "sh.keptn.event.rollback.finished": {
                "eventId": "e96de56c-686b-4a6f-9f3f-1d985c4a0c08",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538154479748633"
              },
              "sh.keptn.event.rollback.started": {
                "eventId": "d02a0622-b64d-4fdf-afa3-23124d45b681",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538150976779102"
              },
              "sh.keptn.event.rollback.triggered": {
                "eventId": "8677cc3d-0e63-4ae8-9f79-226cd95afca0",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538150972262751"
              },
              "sh.keptn.event.staging.delivery.finished": {
                "eventId": "6ce23557-fddd-4ad1-ab91-fb5e6d60ccd9",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538150913154913"
              },
              "sh.keptn.event.staging.delivery.triggered": {
                "eventId": "180bc57f-5636-42ca-a6c3-7be684d26a36",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615479382774586482"
              },
              "sh.keptn.event.staging.rollback.finished": {
                "eventId": "c7087fdf-1baf-4973-9848-72165d846401",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538154864148256"
              },
              "sh.keptn.event.staging.rollback.triggered": {
                "eventId": "f3337a2e-9972-4177-b258-2a74df489a37",
                "keptnContext": "384e3f4c-7587-4c3d-97ae-856435193130",
                "time": "1615538150928333958"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-s1",
            "stage": "staging"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478101146553024",
            "openRemediations": null,
            "serviceName": "carts-s2",
            "stage": "staging"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478147230298554",
            "openRemediations": null,
            "serviceName": "carts-s3",
            "stage": "staging"
          }
        ],
        "stageName": "staging"
      },
      {
        "services": [
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610073528643205",
            "deployedImage": "docker.io/keptnexamples/carts:0.12.3",
            "lastEventTypes": {
              "sh.keptn.event.approval.finished": {
                "eventId": "1a8de917-b6fe-41f7-9a44-527e627b82d5",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427877513703274"
              },
              "sh.keptn.event.approval.started": {
                "eventId": "7122bea6-ac22-4602-80e8-28bd1e888f55",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427877149416073"
              },
              "sh.keptn.event.approval.triggered": {
                "eventId": "2fdcf806-9385-475a-a467-1f9ae8013e96",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684220272801093"
              },
              "sh.keptn.event.deployment.finished": {
                "eventId": "25980aed-60d5-4ea2-86f9-96d6f1065941",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427952393692433"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "966e4611-8870-407a-b9e7-5769a403ca04",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427877901920229"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "28b27d7c-ee48-44b7-ac26-51dc0a032386",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427877895673847"
              },
              "sh.keptn.event.production.delivery.finished": {
                "eventId": "8c9df5d3-4726-49ce-aa55-08aee79981bd",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616428042425645745"
              },
              "sh.keptn.event.production.delivery.triggered": {
                "eventId": "390469a0-2a04-4de7-93ca-0ed118713e49",
                "keptnContext": "0bbaaa6b-fd89-4def-ad2c-975beda970cf",
                "time": "1616684220221781528"
              },
              "sh.keptn.event.production.rollback.finished": {
                "eventId": "84799fb0-19ae-4fe8-af7c-51356bf065a5",
                "keptnContext": "9b0c0913-30e0-4d52-a411-fe1ee4627ac6",
                "time": "1615819464464969095"
              },
              "sh.keptn.event.production.rollback.triggered": {
                "eventId": "c534bb02-2dae-4b59-91c3-d4942d0af107",
                "keptnContext": "9b0c0913-30e0-4d52-a411-fe1ee4627ac6",
                "time": "1615819454929076196"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "98f82e28-a3ad-4910-abac-86ac5a8b9d55",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616428042023056568"
              },
              "sh.keptn.event.release.started": {
                "eventId": "7b89d7a8-34a4-4fb2-a44b-9c895a5175d2",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427953169998851"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "3864580a-ffd4-4a38-867a-a9f07ae54f1f",
                "keptnContext": "31dee9a9-e1ed-4088-a651-ee54911afdf1",
                "time": "1616427953164578761"
              },
              "sh.keptn.event.rollback.finished": {
                "eventId": "c6e4dc10-e1f8-4764-8497-8cd905f6e239",
                "keptnContext": "9b0c0913-30e0-4d52-a411-fe1ee4627ac6",
                "time": "1615819464062083354"
              },
              "sh.keptn.event.rollback.started": {
                "eventId": "c1deb452-45fa-4467-b1f5-3689280fc72e",
                "keptnContext": "9b0c0913-30e0-4d52-a411-fe1ee4627ac6",
                "time": "1615819454955918611"
              },
              "sh.keptn.event.rollback.triggered": {
                "eventId": "32d82c0e-4415-4d35-ad15-690ae3df270c",
                "keptnContext": "9b0c0913-30e0-4d52-a411-fe1ee4627ac6",
                "time": "1615819454939348734"
              }
            },
            "openRemediations": null,
            "serviceName": "carts",
            "stage": "production"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1614610143775644760",
            "deployedImage": "docker.io/mongo:4.2.2",
            "lastEventTypes": {
              "sh.keptn.event.deployment.finished": {
                "eventId": "3fa68f4e-4a7c-46e1-900c-d0f2d3e57f95",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610257208520909"
              },
              "sh.keptn.event.deployment.started": {
                "eventId": "ba244aca-5a69-46d2-97bf-20155c5c3311",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610230171143927"
              },
              "sh.keptn.event.deployment.triggered": {
                "eventId": "5682b1ee-9d26-42b6-b6b3-4ea00f05c615",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610230090205588"
              },
              "sh.keptn.event.production.delivery-direct.finished": {
                "eventId": "6a4482ee-9384-41f0-a010-ec0285be43ad",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610260323282857"
              },
              "sh.keptn.event.production.delivery-direct.triggered": {
                "eventId": "a7fe9606-e769-402a-9e2d-d371555c10a9",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610230080067000"
              },
              "sh.keptn.event.release.finished": {
                "eventId": "4e476591-0b45-46ba-8d08-0cf0a114406c",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610257959436867"
              },
              "sh.keptn.event.release.started": {
                "eventId": "ddadd6e1-8845-4557-908e-fa17287b9144",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610257956807854"
              },
              "sh.keptn.event.release.triggered": {
                "eventId": "b4856315-3e5d-4c26-9fbb-7742f6f0892b",
                "keptnContext": "2d68a93e-cd6f-4f1f-8ba4-2db85e5f023a",
                "time": "1614610257951086764"
              }
            },
            "openRemediations": null,
            "serviceName": "carts-db",
            "stage": "production"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478049158525732",
            "openRemediations": null,
            "serviceName": "carts-s1",
            "stage": "production"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478104139026666",
            "openRemediations": null,
            "serviceName": "carts-s2",
            "stage": "production"
          },
          {
            "roots": [],
            "openApprovals": [],
            "creationDate": "1615478149495331147",
            "openRemediations": null,
            "serviceName": "carts-s3",
            "stage": "production"
          }
        ],
        "stageName": "production"
      }
    ]
  }
];

const Projects = projectsData.map((project: Project) => {
  project.stages = project.stages.map(stage => {
    stage.services = stage.services.map(service => {
      return Service.fromJSON(service);
    });
    return Stage.fromJSON(stage);
  });
  return Project.fromJSON(project);
});

export {Projects};
