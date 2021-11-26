const versions = {
  cli: {
    stable: [
      '0.7.0',
      '0.7.1',
      '0.7.2',
      '0.7.3',
      '0.8.0',
      '0.8.1',
      '0.8.2',
      '0.8.3',
      '0.8.4',
      '0.8.5',
      '0.8.6',
      '0.8.7',
      '0.9.0',
      '0.9.1',
      '0.9.2',
      '0.10.0',
    ],
    prerelease: [],
  },
  bridge: {
    stable: [
      '0.7.0',
      '0.7.1',
      '0.7.2',
      '0.7.3',
      '0.8.0',
      '0.8.1',
      '0.8.2',
      '0.8.3',
      '0.8.4',
      '0.8.5',
      '0.8.6',
      '0.8.7',
      '0.9.0',
      '0.9.1',
      '0.9.2',
      '0.10.0',
    ],
    prerelease: [],
  },
  keptn: {
    stable: [
      {
        version: '0.10.0',
        upgradableVersions: ['0.9.1', '0.9.0', '0.9.2'],
      },
      {
        version: '0.9.2',
        upgradableVersions: ['0.9.1', '0.9.0'],
      },
      {
        version: '0.9.1',
        upgradableVersions: ['0.9.0'],
      },
      {
        version: '0.9.0',
        upgradableVersions: ['0.8.7', '0.8.6', '0.8.5', '0.8.4', '0.8.3', '0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.7',
        upgradableVersions: ['0.8.6', '0.8.5', '0.8.4', '0.8.3', '0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.6',
        upgradableVersions: ['0.8.5', '0.8.4', '0.8.3', '0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.5',
        upgradableVersions: ['0.8.4', '0.8.3', '0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.4',
        upgradableVersions: ['0.8.3', '0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.3',
        upgradableVersions: ['0.8.2', '0.8.1', '0.8.0'],
      },
      {
        version: '0.8.2',
        upgradableVersions: ['0.8.1', '0.8.0'],
      },
      {
        version: '0.8.1',
        upgradableVersions: ['0.8.0'],
      },
      {
        version: '0.8.0',
        upgradableVersions: ['0.7.1', '0.7.2', '0.7.3'],
      },
      {
        version: '0.7.3',
        upgradableVersions: ['0.7.0', '0.7.1', '0.7.2'],
      },
      {
        version: '0.7.2',
        upgradableVersions: ['0.7.0', '0.7.1'],
      },
      {
        version: '0.7.1',
        upgradableVersions: ['0.7.0'],
      },
    ],
  },
};

export { versions as VersionResponseMock };
