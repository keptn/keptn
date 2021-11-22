import { Root } from '../../../_models/root';

const rootsData = [
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.3',
        },
      },
      deployment: {
        deploymentURIsLocal: null,
        deploymentstrategy: '',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'dev',
    },
    id: 'a16a727f-2d04-42cd-acf0-2b99920ff7be',
    source: 'https://github.com/keptn/keptn/cli#configuration-change',
    specversion: '1.0',
    time: '2021-11-09T15:15:14.274Z',
    type: 'sh.keptn.event.dev.delivery.triggered',
    shkeptncontext: '77baf26f-f64d-4a68-9ab5-efde9276ee73',
    shkeptnspecversion: '0.2.3',
  },
];

const root = [Root.fromJSON(rootsData[0])];
export { root as rootResultMock };
