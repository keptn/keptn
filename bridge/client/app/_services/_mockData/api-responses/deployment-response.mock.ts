import { Deployment } from '../../../_models/deployment';

const deploymentData = {
  state: 'finished',
  keptnContext: '77baf26f-f64d-4a68-9ab5-efde9276ee73',
  service: 'carts',
  stages: [
    {
      name: 'dev',
      lastTimeUpdated: '2021-11-09T15:16:36.095Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.dev.delivery.triggered',
          result: 'pass',
          time: '2021-11-09T15:15:14.274Z',
          state: 'finished',
          id: 'a16a727f-2d04-42cd-acf0-2b99920ff7be',
          message: 'Finished release',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: true,
      latestEvaluation: {
        traces: [],
        data: {
          evaluation: {
            indicatorResults: null,
            result: 'pass',
            score: 0,
            sloFileContent: '',
            timeEnd: '2021-11-09T15:16:33Z',
            timeStart: '2021-11-09T15:16:27Z',
          },
          message: 'no evaluation performed by lighthouse because no SLI-provider configured for project sockshop',
          project: 'sockshop',
          result: 'pass',
          service: 'carts',
          stage: 'dev',
          status: 'succeeded',
        },
        id: '427123bf-5336-40d0-9ac6-4eabadb66867',
        source: 'lighthouse-service',
        specversion: '1.0',
        time: '2021-11-09T15:16:33.707Z',
        type: 'sh.keptn.event.evaluation.finished',
        shkeptncontext: '77baf26f-f64d-4a68-9ab5-efde9276ee73',
        shkeptnspecversion: '0.2.3',
        triggeredid: '2e85f33b-bb9d-4e05-b95e-441da9236f77',
      },
    },
    {
      name: 'staging',
      lastTimeUpdated: '2021-11-10T08:11:34.992Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.staging.delivery.triggered',
          result: 'pass',
          time: '2021-11-09T15:16:36.190Z',
          state: 'finished',
          id: '0bfbbdd3-cc75-417b-8aa2-dc8ca15d18ed',
          message: 'Finished release',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: true,
      latestEvaluation: {
        traces: [],
        data: {
          evaluation: {
            indicatorResults: [
              {
                displayName: 'Response time P95',
                keySli: false,
                passTargets: [
                  {
                    criteria: '<=+10%',
                    targetValue: 0,
                    violated: false,
                  },
                  {
                    criteria: '<600',
                    targetValue: 600,
                    violated: false,
                  },
                ],
                score: 1,
                status: 'pass',
                value: {
                  metric: 'response_time_p95',
                  success: true,
                  value: 379.64400029590837,
                },
                warningTargets: [
                  {
                    criteria: '<=800',
                    targetValue: 800,
                    violated: false,
                  },
                ],
              },
            ],
            result: 'pass',
            score: 100,
            sloFileContent:
              'LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiUmVzcG9uc2UgdGltZSBQOTUiDQogICAga2V5X3NsaTogZmFsc2UNCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3INCiAgICB3YXJuaW5nOiAgICAgICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQogICAgd2VpZ2h0OiAxDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
            timeEnd: '2021-11-09T15:20:53Z',
            timeStart: '2021-11-09T15:18:54Z',
          },
          labels: {
            DtCreds: 'dynatrace',
          },
          project: 'sockshop',
          result: 'pass',
          service: 'carts',
          stage: 'staging',
          status: 'succeeded',
        },
        id: '80da2703-0df4-4995-80f6-39cae6fec4e5',
        source: 'lighthouse-service',
        specversion: '1.0',
        time: '2021-11-09T15:22:56.501Z',
        type: 'sh.keptn.event.evaluation.finished',
        shkeptncontext: '77baf26f-f64d-4a68-9ab5-efde9276ee73',
        shkeptnspecversion: '0.2.3',
        triggeredid: 'b22b4575-248c-44fa-bc27-8167ea765db4',
      },
    },
    {
      name: 'production',
      lastTimeUpdated: '2021-11-10T08:13:57.796Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.production.delivery.triggered',
          result: 'pass',
          time: '2021-11-10T08:11:35.092Z',
          state: 'finished',
          id: '0f4ff2c8-5009-4101-90a1-941bbe0171ee',
          message: 'Finished release',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: false,
    },
  ],
  labels: {
    DtCreds: 'dynatrace',
  },
  image: 'carts:0.12.3',
};

const deployment = Deployment.fromJSON(deploymentData);

export { deployment as DeploymentResponseMock };
