import { Resource } from '../../shared/interfaces/resource';

const response: Resource = {
  metadata: {
    branch: 'master',
    upstreamURL: 'https://github.com/ermin-muratovic/test-keptn.git',
    version: '08442ea52b1a4587383f73f68ac92bb6f606a681',
  },
  resourceContent:
    'YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjIiDQpraW5kOiAiU2hpcHlhcmQiDQptZXRhZGF0YToNCiAgbmFtZTogInNoaXB5YXJkLXF1YWxpdHktZ2F0ZXMiDQpzcGVjOg0KICBzdGFnZXM6DQogIC0gbmFtZTogInF1YWxpdHktZ2F0ZSINCiAgICBzZXF1ZW5jZXM6DQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24NCiAgICAgICAgdGFza3M6DQogICAgICAgIC0gbmFtZTogbW9uYWNvDQogICAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg0KICAgICAgLSBuYW1lOiBzaW1wbGVfZXZhbHVhdGlvbg0KICAgICAgICB0YXNrczoNCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg0KICAgICAgICB0cmlnZ2VyZWRBZnRlcjogIjZtIg0KICAgICAgICBwcm9wZXJ0aWVzOg0KICAgICAgICAgIHRpbWVmcmFtZTogIjZtIiAgICAgIA0K',
  resourceURI: 'shipyard.yaml',
};

export { response as ShipyardInvalidResponse };
