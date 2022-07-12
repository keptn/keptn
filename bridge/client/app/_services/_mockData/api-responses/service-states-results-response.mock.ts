const serviceStates = [
  {
    deploymentInformation: [
      {
        stages: [
          {
            name: 'dev',
            time: '2021-11-11T15:24:07.170Z',
          },
          {
            name: 'staging',
            time: '2021-11-11T15:28:22.406Z',
          },
        ],
        name: 'carts',
        image: 'carts',
        version: '0.12.3',
        keptnContext: '08547346-c845-4f49-acab-9f0b9301067e',
      },
      {
        stages: [
          {
            name: 'production',
            time: '2021-11-11T13:21:21.359Z',
          },
        ],
        name: 'carts',
        image: 'carts',
        version: '0.12.3',
        keptnContext: '76f0b0af-0290-458e-82da-56bec6ec5868',
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
            time: '2021-11-09T15:15:24.800Z',
          },
          {
            name: 'staging',
            time: '2021-11-09T15:15:54.741Z',
          },
          {
            name: 'production',
            time: '2021-11-09T15:16:18.314Z',
          },
        ],
        name: 'carts-db',
        image: 'mongo',
        version: '4.2.2',
        keptnContext: '6fd2e002-a732-463e-a449-178ef2f183a7',
      },
    ],
    name: 'carts-db',
  },
];

export { serviceStates as ServiceStatesResultResponseMock };
