import { Resource } from '../../shared/interfaces/resource';

const response: Resource = {
  resourceURI: 'shipyard.yaml',
  metadata: {
    branch: 'master',
    upstreamURL: 'https://github.com/Kirdock/keptn-dynatrace-v1',
    version: 'cbba5536042af2b25cef77dc7c27cf735a93015f',
  },
  resourceContent:
    'YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjIiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXF1YWxpdHktZ2F0ZXMiCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogImRldiIKICAgIHNlcXVlbmNlczoKICAgICAgLSBuYW1lOiByb2xsYmFjawogIC0gbmFtZTogInN0YWdpbmci',
};

export { response as ShipyardEmptySequenceResponse };
