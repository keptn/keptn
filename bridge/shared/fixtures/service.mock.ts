const serviceMock = {
  creationDate: '1632999819358104940',
  deployedImage: 'docker.io/keptnexamples/carts:0.12.3',
  lastEventTypes: {
    'sh.keptn.event.approval.finished': {
      eventId: 'a7b4db51-c004-4a16-b4b9-eb290ac64ae8',
      keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      time: '1636108936755163267',
    },
    'sh.keptn.event.approval.started': {
      eventId: 'b1079092-4c45-4583-abcc-c248528f7dd4',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636115048858051676',
    },
    'sh.keptn.event.approval.triggered': {
      eventId: 'c218f363-4219-45d2-863f-f3fc2329fb16',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636115048772239285',
    },
    'sh.keptn.event.deployment.finished': {
      eventId: '88c5b44b-9b97-4ffb-9349-f0e431efe519',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114893991052283',
    },
    'sh.keptn.event.deployment.started': {
      eventId: '06c24857-a010-4bdd-bcac-639e6ae7d3ed',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114806860574185',
    },
    'sh.keptn.event.deployment.triggered': {
      eventId: 'ff94e27d-4309-4fa2-8cf3-20f266f9e244',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114806657239737',
    },
    'sh.keptn.event.dev.delivery.finished': {
      eventId: '0f44f059-d2ef-4cbe-875f-7345f90711fb',
      keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      time: '1636108939659404279',
    },
    'sh.keptn.event.dev.delivery.triggered': {
      eventId: 'a9b2c7dd-f9c0-4f5a-94b4-b0ce60ac97d2',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114805394927785',
    },
    'sh.keptn.event.dev.evaluation.finished': {
      eventId: '16338541-30a0-4a6b-b990-d4e81d0dd9fd',
      keptnContext: '79a80508-b9c5-4b71-a8c7-98b39bb4b5dd',
      time: '1635866557202816616',
    },
    'sh.keptn.event.dev.evaluation.triggered': {
      eventId: 'c4b192e5-e33f-441a-bae8-e0922a80aaa4',
      keptnContext: '79a80508-b9c5-4b71-a8c7-98b39bb4b5dd',
      time: '1635866530701681563',
    },
    'sh.keptn.event.dev.remediation.finished': {
      eventId: '2e607d76-cbcf-4749-9ddf-db55d725f8f3',
      keptnContext: '2d343938-3632-4432-b232-323530393438',
      time: '1636001467063433773',
    },
    'sh.keptn.event.dev.remediation.triggered': {
      eventId: '3a41c9a6-c124-4ea4-8ac6-52db7af92e79',
      keptnContext: '2d343938-3632-4432-b232-323530393438',
      time: '1636001465183051256',
    },
    'sh.keptn.event.evaluation.finished': {
      eventId: '1064035a-a908-4165-9368-f70297716e34',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636115048667099257',
    },
    'sh.keptn.event.evaluation.started': {
      eventId: '9fa91160-b91d-4980-9b39-9e2aae13466c',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114917956875836',
    },
    'sh.keptn.event.evaluation.triggered': {
      eventId: '6a456cc7-dcaa-4fb5-8c19-f74f0d999ab6',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114917876402210',
    },
    'sh.keptn.event.get-sli.finished': {
      eventId: 'b37c2058-1c19-4a27-996d-3d2256aa40fc',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636115047365792075',
    },
    'sh.keptn.event.get-sli.started': {
      eventId: 'e4e03769-94e7-4ecc-9728-bab6f795de58',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114920461529175',
    },
    'sh.keptn.event.get-sli.triggered': {
      eventId: '06a169b2-cd9e-4e46-a9fc-1b132b42d8de',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114920456530146',
    },
    'sh.keptn.event.release.finished': {
      eventId: 'c23e32de-10d6-4202-b8ad-178234d61628',
      keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      time: '1636108937459912381',
    },
    'sh.keptn.event.release.started': {
      eventId: 'fe73d937-c4ae-41a1-802d-487cc693aba2',
      keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      time: '1636108937455992116',
    },
    'sh.keptn.event.release.triggered': {
      eventId: 'b5488157-c683-4ae3-a5a2-fa6773cd4611',
      keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      time: '1636108937354477796',
    },
    'sh.keptn.event.test.finished': {
      eventId: '47013a83-9bb2-485a-b55c-43db355d6117',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114917763904168',
    },
    'sh.keptn.event.test.started': {
      eventId: 'c822f8a3-1ef0-4768-b525-9c87f932f251',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114894258451746',
    },
    'sh.keptn.event.test.triggered': {
      eventId: '9f98bfaf-35bf-44e1-90c2-f5c7fb740a39',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636114894156898359',
    },
    'sh.keptn.events.problem': {
      eventId: '0548ed6a-bc6d-4df7-9623-6f911c94ee7a',
      keptnContext: '2d343938-3632-4432-b232-323530393438',
      time: '1636002113431434060',
    },
  },
  serviceName: 'carts',
};

export { serviceMock as ServiceMock };
