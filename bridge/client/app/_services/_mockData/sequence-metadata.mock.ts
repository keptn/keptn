const sequenceMetadata = {
  deployments: [
    {
      stage: {
        name: 'dev',
        services: [
          {
            name: 'carts',
            image: 'carts:0.12.3',
          },
          {
            name: 'carts-db',
            image: 'mongo:4.2.2',
          },
        ],
      },
    },
    {
      stage: {
        name: 'staging',
        services: [
          {
            name: 'carts',
            image: 'carts:0.12.3',
          },
          {
            name: 'carts-db',
            image: 'mongo:4.2.2',
          },
        ],
      },
    },
    {
      stage: {
        name: 'production',
        services: [
          {
            name: 'carts',
            image: 'carts:0.12.3',
          },
          {
            name: 'carts-db',
            image: 'mongo:4.2.2',
          },
        ],
      },
    },
  ],
  filter: {
    stages: ['dev', 'production', 'staging'],
    services: ['carts-db', 'carts'],
  },
};

export { sequenceMetadata as SequenceMetadataMock };
