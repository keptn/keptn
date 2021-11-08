const response = {
  nextPageKey: '0',
  stages: [
    {
      services: [
        {
          creationDate: '1633418377048040308',
          deployedImage: 'docker.io/keptnexamples/carts:0.12.1',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: 'a8703786-8413-4826-9ca2-d09552bca16c',
              keptnContext: 'dd6204af-aaa4-41e7-baaa-97f1d075fc03',
              time: '1634219591556383033',
            },
            'sh.keptn.event.deployment.started': {
              eventId: '6d8b2906-e539-4b19-a445-986b130fad33',
              keptnContext: 'dd6204af-aaa4-41e7-baaa-97f1d075fc03',
              time: '1634219498351380701',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: 'd2440fff-2afd-40b1-a1f0-2f2976dfb894',
              keptnContext: 'dd6204af-aaa4-41e7-baaa-97f1d075fc03',
              time: '1634219498250616333',
            },
            'sh.keptn.event.dev.delivery.finished': {
              eventId: '27a14ff2-5a45-4a0d-b451-479c69835cfa',
              keptnContext: 'dd6204af-aaa4-41e7-baaa-97f1d075fc03',
              time: '1634219519760285655',
            },
            'sh.keptn.event.dev.delivery.triggered': {
              eventId: '2af0543a-b25f-450d-9898-074092419414',
              keptnContext: 'dd6204af-aaa4-41e7-baaa-97f1d075fc03',
              time: '1634219496348491714',
            },
            'sh.keptn.event.dev.evaluation.finished': {
              eventId: '80c3dfc2-6fb0-49d7-ad4d-ddbe040ff1ec',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290208545331838',
            },
            'sh.keptn.event.dev.evaluation.triggered': {
              eventId: '7632a141-378d-4aaa-b15a-f462978d5d8d',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290168843331004',
            },
            'sh.keptn.event.evaluation.finished': {
              eventId: 'e312754e-dcc5-4c68-beb3-e213476d2358',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290208352123386',
            },
            'sh.keptn.event.evaluation.started': {
              eventId: 'fdf43c58-6b77-4cdb-8ad8-70aa47762553',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290181059580334',
            },
            'sh.keptn.event.evaluation.triggered': {
              eventId: '821a1adc-001e-4769-b662-79691f330eed',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290170653051109',
            },
            'sh.keptn.event.get-sli.finished': {
              eventId: '258ec098-ad53-4c39-84f9-a0ac073e68c6',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290206550509437',
            },
            'sh.keptn.event.get-sli.started': {
              eventId: '9e350dad-918e-42da-b0e0-fcebdcaf5d0f',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290172179310515',
            },
            'sh.keptn.event.get-sli.triggered': {
              eventId: '74888f99-a9df-422c-a27a-165da3299616',
              keptnContext: '275eca1e-5c3e-4373-86ab-1b1c271503f7',
              time: '1634290172173795694',
            },
            'sh.keptn.event.release.finished': {
              eventId: 'f6847e26-f831-4f05-8ca3-dbf3d3fdaa9b',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115339748455992',
            },
            'sh.keptn.event.release.started': {
              eventId: '44e46383-9224-4b8a-a391-255d0ec69ddb',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115339847069791',
            },
            'sh.keptn.event.release.triggered': {
              eventId: '3890f4cd-0387-4ae6-ab2b-fee146a03c1d',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115339743540831',
            },
            'sh.keptn.event.test.finished': {
              eventId: 'b2f775c7-3363-42ef-9016-ab53a88bff9c',
              keptnContext: '766d0224-de3f-4dc3-9a6f-2fc182dd98c8',
              time: '1634212301152663652',
            },
            'sh.keptn.event.test.started': {
              eventId: 'd7245d2b-2550-4225-8786-7d7b119ff407',
              keptnContext: '766d0224-de3f-4dc3-9a6f-2fc182dd98c8',
              time: '1634212266547128150',
            },
            'sh.keptn.event.test.triggered': {
              eventId: '13fa1404-727b-4c58-b207-77ab37688c58',
              keptnContext: '766d0224-de3f-4dc3-9a6f-2fc182dd98c8',
              time: '1634212266543111780',
            },
          },
          openRemediations: null,
          serviceName: 'carts',
        },
        {
          creationDate: '1633418401151119585',
          deployedImage: 'docker.io/mongo:4.2.2',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: 'b2c0ef8f-b241-41c8-904d-02be2d5a7671',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024971267804221',
            },
            'sh.keptn.event.deployment.started': {
              eventId: 'a0cf756e-1c80-4b41-bfdb-16abaf26df26',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024916567531337',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: 'd9b150c7-f773-4b8d-abf3-66cb4f027237',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024901859996636',
            },
            'sh.keptn.event.dev.delivery-direct.finished': {
              eventId: '615dac0b-938e-405f-a62d-a40454f84a2f',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024973564679058',
            },
            'sh.keptn.event.dev.delivery-direct.triggered': {
              eventId: '511ff8b2-5430-46c5-96e7-2f65880d9924',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024900121181962',
            },
            'sh.keptn.event.release.finished': {
              eventId: '8d287b2a-b7a6-47b0-800a-cadeb502d429',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025117359042118',
            },
            'sh.keptn.event.release.started': {
              eventId: 'd97c209f-eec3-455f-b57a-b239d688b201',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025042652215823',
            },
            'sh.keptn.event.release.triggered': {
              eventId: 'e1abb15d-2c85-46d5-83c3-dea456af9633',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024971552270760',
            },
          },
          openRemediations: null,
          serviceName: 'carts-db',
        },
      ],
      stageName: 'dev',
    },
    {
      services: [
        {
          creationDate: '1633418377949679695',
          deployedImage: 'docker.io/keptnexamples/carts:0.12.1',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: 'f92f51dc-f769-44a6-b2c6-507664c4157e',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026181357351242',
            },
            'sh.keptn.event.deployment.started': {
              eventId: 'f57697cf-7347-4dc5-8fee-061e6e1f0a06',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026039158764716',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: '71a2e4b4-5d92-4288-95f1-e3dd614625fa',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026013650480599',
            },
            'sh.keptn.event.production.delivery.finished': {
              eventId: '7b198133-abf0-4915-b33a-15f21d373764',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026274650753296',
            },
            'sh.keptn.event.production.delivery.triggered': {
              eventId: '9d9057ea-23c4-415d-8294-29abf2f7af15',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026006950177824',
            },
            'sh.keptn.event.release.finished': {
              eventId: '4398b5ed-74a0-4951-a03d-fa81176300fd',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026274516800840',
            },
            'sh.keptn.event.release.started': {
              eventId: 'cf075d00-9028-41dc-bb2c-33780ce7fc10',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026223954549463',
            },
            'sh.keptn.event.release.triggered': {
              eventId: '6db843c1-e394-4a9a-9e39-ecdb8a7467a7',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026181651357481',
            },
          },
          openRemediations: null,
          serviceName: 'carts',
        },
        {
          creationDate: '1633418401854937845',
          deployedImage: 'docker.io/mongo:4.2.2',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: '21a9020a-e946-49c1-802e-4943b41bb5b3',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025360859296290',
            },
            'sh.keptn.event.deployment.started': {
              eventId: '87179ebb-7f94-423a-bcc6-71a5d3c49399',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025290760115022',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: '7f60f49d-122e-4baa-aea3-648b96024463',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025199644800439',
            },
            'sh.keptn.event.production.delivery-direct.finished': {
              eventId: 'dab521c9-c77e-47ce-9166-b79d1842774f',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025363342434248',
            },
            'sh.keptn.event.production.delivery-direct.triggered': {
              eventId: '4d1dbd08-910e-4ed6-a44e-c9bb63aea998',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025164350668076',
            },
            'sh.keptn.event.release.finished': {
              eventId: '29bf8f4a-180b-47c0-8c65-b462572b0404',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025562456054396',
            },
            'sh.keptn.event.release.started': {
              eventId: '8dd9126e-73f7-4540-b26b-240e3a48b08a',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025463457689022',
            },
            'sh.keptn.event.release.triggered': {
              eventId: '574cff9e-d1e7-49bc-b95b-2cbc3852f437',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025361150140460',
            },
          },
          openRemediations: null,
          serviceName: 'carts-db',
        },
      ],
      stageName: 'production',
      parentStages: ['staging'],
    },
    {
      services: [
        {
          creationDate: '1633418377354372259',
          deployedImage: 'docker.io/keptnexamples/carts:0.12.1',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: 'f770326d-2674-43e0-a99f-a6b358c7b863',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115449949038124',
            },
            'sh.keptn.event.deployment.started': {
              eventId: '080e857f-4dd6-4f4d-be56-cd83d82d7e62',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115345943219783',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: 'febfe371-56ad-44e8-a73c-d6103dd52aa4',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115345848666108',
            },
            'sh.keptn.event.evaluation.finished': {
              eventId: '3e047ce6-209b-4c25-948a-1de830d01285',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115654259970403',
            },
            'sh.keptn.event.evaluation.started': {
              eventId: 'ed1870aa-892b-482c-8897-8886044bdac4',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115578461133447',
            },
            'sh.keptn.event.evaluation.triggered': {
              eventId: '68def93d-aeb9-48f3-afbb-7af70e103e22',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115574144106880',
            },
            'sh.keptn.event.get-sli.finished': {
              eventId: 'a875ce72-7f0c-4e89-9463-cb9c7ee4b313',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115650423158286',
            },
            'sh.keptn.event.get-sli.started': {
              eventId: 'e0bb6cc2-6f0a-4a34-9039-2649883407a2',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115576874778953',
            },
            'sh.keptn.event.get-sli.triggered': {
              eventId: '8e4d63c0-5c85-4228-a9a5-4750edd89ccf',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115576871066957',
            },
            'sh.keptn.event.release.finished': {
              eventId: '2f5864af-45ae-4795-875b-04ea3d07dce3',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634026006649470338',
            },
            'sh.keptn.event.release.started': {
              eventId: '0a0826bd-cb4f-46c0-b34e-03e3ae753a8f',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634025937957237564',
            },
            'sh.keptn.event.release.triggered': {
              eventId: 'ec88bf9b-7121-496d-81d4-5210ca115b13',
              keptnContext: '9ad16886-a979-4be3-8180-31b8b3ae1b22',
              time: '1634025890344789468',
            },
            'sh.keptn.event.rollback.finished': {
              eventId: '73270021-8874-494f-9ded-a44df691e01e',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115685864608993',
            },
            'sh.keptn.event.rollback.started': {
              eventId: 'ef24b9b6-9a15-4503-8cbb-e188b32e09bd',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115658348172019',
            },
            'sh.keptn.event.rollback.triggered': {
              eventId: 'f8a77af6-de01-4024-9de3-fd72f84adc13',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115658344561699',
            },
            'sh.keptn.event.staging.delivery.finished': {
              eventId: 'c2327d02-f306-4a95-b5e2-3f9bd16cbc86',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115654746159564',
            },
            'sh.keptn.event.staging.delivery.triggered': {
              eventId: '24a9b374-1218-47da-a687-6f7804c08e6f',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115342144538405',
            },
            'sh.keptn.event.staging.evaluation.finished': {
              eventId: 'd5a60670-b658-4dde-b99c-d8912b894f65',
              keptnContext: 'c0d1875d-6a65-4b15-8cc8-47e8d85315ce',
              time: '1634113565343859164',
            },
            'sh.keptn.event.staging.evaluation.triggered': {
              eventId: '95a68650-3588-4f85-a430-3c46d05c9831',
              keptnContext: 'c0d1875d-6a65-4b15-8cc8-47e8d85315ce',
              time: '1634113533967218772',
            },
            'sh.keptn.event.staging.rollback.finished': {
              eventId: 'dc5d67d1-79b6-421b-81c9-cbcd64b47ae4',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115686444113590',
            },
            'sh.keptn.event.staging.rollback.triggered': {
              eventId: '7ef762aa-e52c-4d99-a297-0c4b7df009d3',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115655344330890',
            },
            'sh.keptn.event.test.finished': {
              eventId: 'cdce6905-88f0-4f96-93db-d008f83afbe3',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115573647783353',
            },
            'sh.keptn.event.test.started': {
              eventId: '8539be6e-ea69-4a97-9a27-facceee54851',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115450346439912',
            },
            'sh.keptn.event.test.triggered': {
              eventId: '2b729a53-44ab-4637-acbd-a637c6bfbe21',
              keptnContext: '4c37cb99-bd53-455c-be73-1937bb0d8c36',
              time: '1634115450243029924',
            },
          },
          openRemediations: null,
          serviceName: 'carts',
        },
        {
          creationDate: '1633418401552780431',
          deployedImage: 'docker.io/mongo:4.2.2',
          lastEventTypes: {
            'sh.keptn.event.deployment.finished': {
              eventId: '08adeff0-3429-4a67-b0cf-3029e7ccb133',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025161957567607',
            },
            'sh.keptn.event.deployment.started': {
              eventId: '446760ee-253c-476d-994f-93e4992f6548',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025075849909748',
            },
            'sh.keptn.event.deployment.triggered': {
              eventId: '07e95add-37a9-419d-9c5b-7e4783e2937e',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024996248154850',
            },
            'sh.keptn.event.release.finished': {
              eventId: 'a74c94b0-423d-4d3b-9fec-b280c188a28b',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025329856454710',
            },
            'sh.keptn.event.release.started': {
              eventId: '3577c0b8-18cf-49eb-aae8-a8bc968aa741',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025251060062509',
            },
            'sh.keptn.event.release.triggered': {
              eventId: '9f7aa659-da6d-4397-9b20-c4ca4340b3b8',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025162245017189',
            },
            'sh.keptn.event.staging.delivery-direct.finished': {
              eventId: '6ad2b4ba-7334-4d72-8349-552c44cfaa40',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634025164262597768',
            },
            'sh.keptn.event.staging.delivery-direct.triggered': {
              eventId: 'd31b8ba3-7b7b-49e8-915b-b66f95ed104e',
              keptnContext: 'd55e27e7-53fb-4ba3-ac43-1617cf1f0644',
              time: '1634024973648393363',
            },
          },
          openRemediations: null,
          serviceName: 'carts-db',
        },
      ],
      stageName: 'staging',
      parentStages: ['dev'],
    },
  ],
  totalCount: 3,
};

export { response as StagesResponse };
