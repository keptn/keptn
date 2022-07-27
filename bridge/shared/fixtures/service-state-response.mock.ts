const serviceStateResponseMock = [
  {
    deploymentInformation: [
      {
        stages: [
          {
            name: 'dev',
            time: '2021-11-05T12:21:33.991Z',
          },
        ],
        name: 'carts',
        image: 'carts',
        version: '0.12.3',
        keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      },
      {
        stages: [
          {
            name: 'staging',
            time: '2021-11-05T10:49:01.288Z',
          },
        ],
        name: 'carts',
        image: 'carts',
        version: '0.12.3',
        keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      },
      {
        stages: [
          {
            name: 'production',
            time: '2021-10-13T11:01:18.567Z',
          },
        ],
        name: 'carts',
        image: 'carts',
        version: '0.12.1',
        keptnContext: '2c0e568b-8bd3-4726-a188-e528423813ed',
      },
    ],
    name: 'carts',
  },
  {
    deploymentInformation: [
      {
        stages: [
          {
            name: 'dev',
            time: '2021-10-12T11:13:18.563Z',
          },
          {
            name: 'staging',
            time: '2021-10-12T11:13:59.263Z',
          },
          {
            name: 'production',
            time: '2021-10-12T11:14:57.469Z',
          },
        ],
        name: 'carts-db',
        image: 'mongo',
        version: '4.2.2',
        keptnContext: '0cc574e9-3d47-4a29-81b7-84faf33bdc9c',
      },
    ],
    name: 'carts-db',
  },
];

const serviceStateQualityGatesOnlyResponse = [
  {
    deploymentInformation: [
      {
        stages: [
          {
            name: 'dev',
            time: '2021-11-05T12:24:08.667Z',
          },
        ],
        name: 'carts',
        keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      },
      {
        stages: [
          {
            name: 'staging',
            time: '2021-11-05T10:52:31.167Z',
          },
        ],
        name: 'carts',
        keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      },
    ],
    name: 'carts',
  },
  {
    deploymentInformation: [],
    name: 'carts-db',
  },
];

export { serviceStateResponseMock as ServiceStateResponse };
export { serviceStateQualityGatesOnlyResponse as ServiceStateQualityGatesOnlyResponse };
