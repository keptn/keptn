import { TestUtils } from '../../_utils/test.utils';
import { Trace } from '../../_models/trace';

const rootTraces = [
  {
    contenttype: 'application/json',
    data: {
      deploymentStrategies: {},
      eventContext: null,
      helmChart:
        'H4sIFAAAAAAA/ykAK2FIUjBjSE02THk5NWIzVjBkUzVpWlM5Nk9WVjZNV2xqYW5keVRRbz1IZWxtAOxXS2/iSBDOGYn/UMplTjFtCElkaQ8RsDPsBrB4RFqtVqhjV6A3bXdvdxsNGvHfV35hYzyTkSYPReK7gLuqPle1q9qfPaqMbvXWVBlrSwN+9goghJDrbjf5JYRUf4ndvjqzO1edztV1u33dPSN2u9PungF5jWSqiLSh6oz88r2qxX0QUMnuUWkmQgc2dsNH7SkmTXJ9C1+QB+DFzQGPQoFGtWEeQtI0jZAG6GT/NzkHsWyLNN67qhN+Fun8byiPUL/WAfDc/JPry+r8d69O8/8mYAFdoQO+8J5QWUy0nlCaEL/SQHLUraQ9HGLZxLKbDYWSM4/2RBQaB+z3zv2EX0c6/wYDyalB3fJRcrENMHxBOfDM/MdjX5n/7mW3c5r/t8DFxQU0G2UVQKXUrY3dbDyx0Hegv2+IZiNAQ31qqNNsAJRe/82Gluglq9kRoR349g2s+/S9Uj43YLeL/bRR1OBqmwQBKME5C1cL6VOD2RpAQL8uQrqhjNMHjg6Q1GC2Eh2YlkMSSuToGaGy8IAab31HH5DrPSGVcp8yQN71eUBeHeTu/DC6Gg+wLzuGJ0JDWYiqiLg43KWcJjt0z0tblCzBbnde8XIjzl3Bmbd1YPg4FsZVqJOHkbtJER/RxXV+z7UxslgFkEoY4QnuwLznlg37vF2hjAM35IYUZgw3Ndz9+bK3mM0no6U7nRyQJUrCgfPkNbKUSvyLnvmtVGhisDID7HaQembS8tgz15yFp6GrOj+TbmDmVZxjx66F7XDD8+rcSX85vh0Njur6XYnAKa8CPDLk/hQfK8uZwaVm7cB53llWzF9zw/7AvZv8NRqM569037SR//5UVP7pn5o8epPx/HY4HkyXw9Ht5+M8nm3ZnOjPgTsfx73xx6A3/x5N+tE5psEPSWbzukxe5EloSb26x5HdeDC9H/ZqNyEZ55rAxfhucDv7koQOpsvF9K4uOh5Mp9WKQo5Ur63s98LHTYtKVqLlbIMhau0q8YAHZcUUn9FUSpVJja01Um7WFVPNaAOwkBlGeR853c7QE6GvHbg6cJGomPD3RvvAaFiAIjKFtVtYFVKffeDstYiUh/ogSc4CZvRRg3kyirkJCaqWAAOhtg60yeXNiJWtCv+LUH+PrPsDLpu0L2Oul3n/V/Vfdty+6Kfgc99/HbtT0X+X121y0n9vgWP9V0i/WdoLJd0H9cIvWU5VWY9H2qAaJqIglSaJ9ViV5EOdXR2rE0PVCs2hKNlrvEyjlQXZe2/lCSeccMKHwv8BAAD//7KJCEUAGgAA',
      project: 'sockshop',
      service: 'carts',
    },
    id: 'ade14d1a-338d-4a88-ad88-7b03ec9b5d8e',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T12:45:01.643Z',
    type: 'sh.keptn.event.service.create.started',
    shkeptncontext: '56255df2-22a2-45ef-b2c4-cf4882096b3f',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.1' },
      },
    },
    id: 'f95d2c20-2d89-4f52-8838-36743cf8835f',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T12:47:05.910Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '218ddbfa-ed09-4cf9-887a-167a334a76d0',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.2' },
      },
    },
    id: 'd72b0a9d-24d3-48ec-832f-e36696c466d4',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T13:48:48.749Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '6f8bf547-d588-4866-b5f5-1fe43a4b0c65',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.2' },
      },
    },
    id: '21ccce4c-e77e-479c-91fa-7052551bdf48',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T15:09:02.740Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '934b6d92-7605-4eaa-b471-d6022f3e6d72',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.2' },
      },
    },
    id: '84235fd5-d009-4e42-96c2-4cfdcfb48b4d',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T15:54:52.249Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '38bdaef6-60a2-48f3-b474-626a683d175c',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.1' },
      },
    },
    id: '6eacae7f-7446-49a8-b351-0d95bb6be76b',
    source: 'https://github.com/keptn/keptn/cli#configuration-change',
    specversion: '0.2',
    time: '2020-05-28T07:49:29.742Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: 'fea3dc8c-5a85-435a-a86d-cee0b62f248e',
  },
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: '',
      end: '2020-05-28T08:59:00.000Z',
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: 'staging',
      start: '2020-05-28T08:54:00.000Z',
      teststrategy: 'manual',
    },
    id: '0de64842-b277-4629-a5b5-c87829ff1151',
    source: 'https://github.com/keptn/keptn/cli#configuration-change',
    specversion: '0.2',
    time: '2020-05-28T07:53:06.860Z',
    type: 'sh.keptn.event.start-evaluation',
    shkeptncontext: '3302455e-ffe6-4ae3-a514-31150557dd09',
  },
  {
    contenttype: 'application/json',
    data: {
      canary: { action: 'set', value: 100 },
      eventContext: null,
      labels: null,
      project: 'sockshop',
      service: 'carts',
      stage: '',
      configurationChange: {
        values: { image: 'docker.io/keptnexamples/carts:0.10.3' },
      },
    },
    id: 'd9949721-29de-47bf-9b04-b890eeb21e4d',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '0.2',
    time: '2020-03-16T16:14:29.785Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '42e8e409-5afc-4ee5-abdb-f41926ab2583',
  },
  {
    data: {
      configurationChange: {
        values: {},
      },
      deployment: {
        deploymentstrategy: '',
      },
      project: 'keptn',
      service: 'control-plane',
      stage: 'dev',
    },
    id: '3b209b06-597c-413e-9401-b80e4855a313',
    source: 'https://github.com/keptn/keptn/cli#configuration-change',
    specversion: '1.0',
    time: '2021-02-02T08:52:39.186Z',
    type: 'sh.keptn.event.dev.artifact-delivery.triggered',
    shkeptncontext: '0ede19b7-dc65-4f04-9882-ddadf3703019',
  },
];
const evaluationTraces = [
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: 'blue_green_service',
      evaluation: {
        indicatorResults: null,
        result: 'no evaluation performed by lighthouse because no SLO found for service carts',
        score: 0,
        sloFileContent: '',
        timeEnd: '2020-03-16T12:52:15Z',
        timeStart: '2020-03-16T12:50:44Z',
      },
      labels: null,
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      teststrategy: 'performance',
    },
    id: 'e8ace82f-3a6b-42a3-bc6b-4dcb93c741d7',
    source: 'lighthouse-service',
    specversion: '0.2',
    time: '2020-03-16T12:52:15.818Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '218ddbfa-ed09-4cf9-887a-167a334a76d0',
  },
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: 'blue_green_service',
      evaluation: {
        indicatorResults: [
          {
            score: 0,
            status: 'fail',
            targets: [
              { criteria: '<=900', targetValue: 900, violated: true },
              {
                criteria: '<600',
                targetValue: 600,
                violated: true,
              },
            ],
            value: { metric: 'response_time_p95', success: true, value: 2013.11088577307 },
          },
          {
            score: 0,
            status: 'fail',
            targets: [
              { criteria: '<=800', targetValue: 800, violated: true },
              {
                criteria: '<300',
                targetValue: 300,
                violated: true,
              },
            ],
            value: { metric: 'response_time_p50', success: true, value: 2010.8712390903988 },
          },
          {
            score: 0,
            status: 'fail',
            targets: [{ criteria: '=0', targetValue: 0, violated: true }],
            value: { metric: 'error_rate', success: true, value: 139 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: { metric: 'throughput', success: true, value: 3 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: { metric: 'cpu_usage', success: true, value: 8.50832800771676 },
          },
        ],
        result: 'fail',
        score: 0,
        sloFileContent:
          'Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=',
        timeEnd: '2020-03-16T16:11:54Z',
        timeStart: '2020-03-16T15:55:03Z',
      },
      labels: null,
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      teststrategy: 'performance',
    },
    id: 'cdcf2cf8-f2f9-4d7d-b62a-d27bba980389',
    source: 'lighthouse-service',
    specversion: '0.2',
    time: '2020-03-16T16:13:56.288Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '38bdaef6-60a2-48f3-b474-626a683d175c',
  },
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: 'blue_green_service',
      evaluation: {
        indicatorResults: [
          {
            score: 1,
            status: 'pass',
            targets: [{ criteria: '<600', targetValue: 600, violated: false }],
            value: { metric: 'response_time_p95', success: true, value: 339.15595320978224 },
          },
          {
            score: 1,
            status: 'pass',
            targets: [{ criteria: '<300', targetValue: 300, violated: false }],
            value: { metric: 'response_time_p50', success: true, value: 158.36151562776496 },
          },
          {
            score: 0,
            status: 'fail',
            targets: [{ criteria: '=0', targetValue: 0, violated: true }],
            value: { metric: 'error_rate', success: true, value: 3.3333333333333335 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: { metric: 'throughput', success: true, value: 3 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: { metric: 'cpu_usage', success: true, value: 15.433689541286892 },
          },
        ],
        result: 'fail',
        score: 66.66666666666666,
        sloFileContent:
          'Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=',
        timeEnd: '2020-03-16T16:18:48Z',
        timeStart: '2020-03-16T16:16:56Z',
      },
      labels: null,
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      teststrategy: 'performance',
    },
    id: '725206de-f7bf-4461-b758-070b77402dc3',
    source: 'lighthouse-service',
    specversion: '0.2',
    time: '2020-03-16T16:20:50.233Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '42e8e409-5afc-4ee5-abdb-f41926ab2583',
  },
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: '',
      evaluation: {
        indicatorResults: [
          {
            score: 0,
            status: 'fail',
            targets: [
              { criteria: '<=900', targetValue: 0, violated: true },
              {
                criteria: '<600',
                targetValue: 0,
                violated: true,
              },
            ],
            value: {
              message: 'end time must not be in the future',
              metric: 'response_time_p95',
              success: false,
              value: 0,
            },
          },
          {
            score: 0,
            status: 'fail',
            targets: [
              { criteria: '<=800', targetValue: 0, violated: true },
              {
                criteria: '<300',
                targetValue: 0,
                violated: true,
              },
            ],
            value: {
              message: 'end time must not be in the future',
              metric: 'response_time_p50',
              success: false,
              value: 0,
            },
          },
          {
            score: 0,
            status: 'fail',
            targets: [{ criteria: '=0', targetValue: 0, violated: true }],
            value: {
              message: 'end time must not be in the future',
              metric: 'error_rate',
              success: false,
              value: 0,
            },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: {
              message: 'end time must not be in the future',
              metric: 'throughput',
              success: false,
              value: 0,
            },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: {
              message: 'end time must not be in the future',
              metric: 'cpu_usage',
              success: false,
              value: 0,
            },
          },
        ],
        result: 'fail',
        score: 0,
        sloFileContent:
          'Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=',
        timeEnd: '2020-05-28T08:59:00.000Z',
        timeStart: '2020-05-28T08:54:00.000Z',
      },
      labels: null,
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      teststrategy: 'manual',
    },
    id: '1f2c7acc-c09e-424c-b4a0-1f8748657ee8',
    source: 'lighthouse-service',
    specversion: '0.2',
    time: '2020-05-28T07:53:07.193Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '3302455e-ffe6-4ae3-a514-31150557dd09',
  },
  {
    contenttype: 'application/json',
    data: {
      deploymentstrategy: 'blue_green_service',
      evaluation: {
        indicatorResults: [
          {
            score: 1,
            status: 'pass',
            targets: [{ criteria: '<600', targetValue: 600, violated: false }],
            value: { metric: 'response_time_p95', success: true, value: 345.8670349893803 },
          },
          {
            score: 1,
            status: 'pass',
            targets: [{ criteria: '<300', targetValue: 300, violated: false }],
            value: { metric: 'response_time_p50', success: true, value: 186.38390148907558 },
          },
          {
            score: 0,
            status: 'fail',
            targets: [{ criteria: '=0', targetValue: 0, violated: true }],
            value: { metric: 'error_rate', success: true, value: 5 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: { metric: 'throughput', success: true, value: 4 },
          },
          {
            score: 0,
            status: 'info',
            targets: null,
            value: {
              message:
                'Dynatrace API returned status code 403: The query involves 1794 metrics, but the limit for REST queries is 10. Consider splitting your query into multiple smaller queries.',
              metric: 'cpu_usage',
              success: false,
              value: 0,
            },
          },
        ],
        result: 'fail',
        score: 66.66666666666666,
        sloFileContent:
          'Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczogbnVsbAogIHNsaTogY3B1X3VzYWdlCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=',
        timeEnd: '2020-05-28T07:54:12Z',
        timeStart: '2020-05-28T07:52:07Z',
      },
      labels: null,
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      teststrategy: 'performance',
    },
    id: '214ca172-4080-4165-9a4d-39f399b17a45',
    source: 'lighthouse-service',
    specversion: '0.2',
    time: '2020-05-28T07:56:15.840Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: 'fea3dc8c-5a85-435a-a86d-cee0b62f248e',
  },
];
const multipleEvaluationsTraces = [
  {
    traces: [
      {
        traces: [
          {
            traces: [],
            data: {
              project: 'dynatrace',
              service: 'items',
              stage: 'quality-gate',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'b6d7982e-d604-43bc-8c2f-022d0cc9be34',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:37.414Z',
            type: 'sh.keptn.event.evaluation.started',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            shkeptnspecversion: '0.2.3',
            triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            plainEvent: {
              data: {
                project: 'dynatrace',
                service: 'items',
                stage: 'quality-gate',
                status: 'succeeded',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: 'b6d7982e-d604-43bc-8c2f-022d0cc9be34',
              source: 'lighthouse-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:37.414Z',
              type: 'sh.keptn.event.evaluation.started',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              shkeptnspecversion: '0.2.3',
              triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            },
          },
          {
            traces: [],
            data: {
              project: 'dynatrace',
              service: 'items',
              stage: 'quality-gate',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'bcc9d50c-2002-452a-9871-9e6c4650d35a',
            source: 'webhook-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:42.933Z',
            type: 'sh.keptn.event.evaluation.started',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            plainEvent: {
              data: {
                project: 'dynatrace',
                service: 'items',
                stage: 'quality-gate',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: 'bcc9d50c-2002-452a-9871-9e6c4650d35a',
              source: 'webhook-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:42.933Z',
              type: 'sh.keptn.event.evaluation.started',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            },
          },
          {
            traces: [],
            data: {
              labels: null,
              message: 'could not retrieve Webhook config: no webhook config found',
              project: 'dynatrace',
              result: 'fail',
              service: 'items',
              stage: 'quality-gate',
              status: 'errored',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '7b0ff004-58c4-4c11-afb4-dc0f70d6adb9',
            source: 'webhook-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:42.939Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            plainEvent: {
              data: {
                labels: null,
                message: 'could not retrieve Webhook config: no webhook config found',
                project: 'dynatrace',
                result: 'fail',
                service: 'items',
                stage: 'quality-gate',
                status: 'errored',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: '7b0ff004-58c4-4c11-afb4-dc0f70d6adb9',
              source: 'webhook-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:42.939Z',
              type: 'sh.keptn.event.evaluation.finished',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            },
          },
          {
            traces: [],
            data: {
              evaluation: {
                comparedEvents: ['5c563cbe-e862-4a42-829c-d319683d3654'],
                indicatorResults: [
                  {
                    displayName: '',
                    keySli: false,
                    passTargets: [
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
                      value: 10.219262258461157,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: false,
                      },
                    ],
                  },
                  {
                    displayName: '',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<5',
                        targetValue: 5,
                        violated: false,
                      },
                    ],
                    score: 1,
                    status: 'pass',
                    value: {
                      metric: 'error_rate',
                      success: true,
                      value: 0,
                    },
                    warningTargets: null,
                  },
                  {
                    displayName: '',
                    keySli: false,
                    passTargets: null,
                    score: 0,
                    status: 'info',
                    value: {
                      metric: 'throughput',
                      success: true,
                      value: 933,
                    },
                    warningTargets: null,
                  },
                ],
                result: 'pass',
                score: 100,
                sloFileContent:
                  'c3BlY192ZXJzaW9uOiAiMS4wIgpmaWx0ZXI6IHt9CmNvbXBhcmlzb246CiAgY29tcGFyZV93aXRoOiBzaW5nbGVfcmVzdWx0CiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDEKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwpvYmplY3RpdmVzOgotIHNsaTogcmVzcG9uc2VfdGltZV9wOTUKICBkaXNwbGF5TmFtZTogIiIKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICB3YXJuaW5nOgogIC0gY3JpdGVyaWE6CiAgICAtIDw9ODAwCiAgd2VpZ2h0OiAxCiAga2V5X3NsaTogZmFsc2UKLSBzbGk6IGVycm9yX3JhdGUKICBkaXNwbGF5TmFtZTogIiIKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw1CiAgd2FybmluZzogW10KICB3ZWlnaHQ6IDEKICBrZXlfc2xpOiBmYWxzZQotIHNsaTogdGhyb3VnaHB1dAogIGRpc3BsYXlOYW1lOiAiIgogIHBhc3M6IFtdCiAgd2FybmluZzogW10KICB3ZWlnaHQ6IDAKICBrZXlfc2xpOiBmYWxzZQp0b3RhbF9zY29yZToKICBwYXNzOiA5MCUKICB3YXJuaW5nOiA3NSUK',
                timeEnd: '2021-10-11T08:55:36.469Z',
                timeStart: '2021-10-11T08:50:36.469Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'dynatrace',
              result: 'pass',
              service: 'items',
              stage: 'quality-gate',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'e859b7ba-c82a-44fd-bd69-313aed264f7e',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:51.316Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            shkeptnspecversion: '0.2.3',
            triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            plainEvent: {
              data: {
                evaluation: {
                  comparedEvents: ['5c563cbe-e862-4a42-829c-d319683d3654'],
                  indicatorResults: [
                    {
                      displayName: '',
                      keySli: false,
                      passTargets: [
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
                        value: 10.219262258461157,
                      },
                      warningTargets: [
                        {
                          criteria: '<=800',
                          targetValue: 800,
                          violated: false,
                        },
                      ],
                    },
                    {
                      displayName: '',
                      keySli: false,
                      passTargets: [
                        {
                          criteria: '<5',
                          targetValue: 5,
                          violated: false,
                        },
                      ],
                      score: 1,
                      status: 'pass',
                      value: {
                        metric: 'error_rate',
                        success: true,
                        value: 0,
                      },
                      warningTargets: null,
                    },
                    {
                      displayName: '',
                      keySli: false,
                      passTargets: null,
                      score: 0,
                      status: 'info',
                      value: {
                        metric: 'throughput',
                        success: true,
                        value: 933,
                      },
                      warningTargets: null,
                    },
                  ],
                  result: 'pass',
                  score: 100,
                  sloFileContent:
                    'c3BlY192ZXJzaW9uOiAiMS4wIgpmaWx0ZXI6IHt9CmNvbXBhcmlzb246CiAgY29tcGFyZV93aXRoOiBzaW5nbGVfcmVzdWx0CiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDEKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwpvYmplY3RpdmVzOgotIHNsaTogcmVzcG9uc2VfdGltZV9wOTUKICBkaXNwbGF5TmFtZTogIiIKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICB3YXJuaW5nOgogIC0gY3JpdGVyaWE6CiAgICAtIDw9ODAwCiAgd2VpZ2h0OiAxCiAga2V5X3NsaTogZmFsc2UKLSBzbGk6IGVycm9yX3JhdGUKICBkaXNwbGF5TmFtZTogIiIKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw1CiAgd2FybmluZzogW10KICB3ZWlnaHQ6IDEKICBrZXlfc2xpOiBmYWxzZQotIHNsaTogdGhyb3VnaHB1dAogIGRpc3BsYXlOYW1lOiAiIgogIHBhc3M6IFtdCiAgd2FybmluZzogW10KICB3ZWlnaHQ6IDAKICBrZXlfc2xpOiBmYWxzZQp0b3RhbF9zY29yZToKICBwYXNzOiA5MCUKICB3YXJuaW5nOiA3NSUK',
                  timeEnd: '2021-10-11T08:55:36.469Z',
                  timeStart: '2021-10-11T08:50:36.469Z',
                },
                labels: {
                  DtCreds: 'dynatrace',
                },
                project: 'dynatrace',
                result: 'pass',
                service: 'items',
                stage: 'quality-gate',
                status: 'succeeded',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: 'e859b7ba-c82a-44fd-bd69-313aed264f7e',
              source: 'lighthouse-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:51.316Z',
              type: 'sh.keptn.event.evaluation.finished',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              shkeptnspecversion: '0.2.3',
              triggeredid: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
            },
          },
        ],
        data: {
          deployment: {
            deploymentNames: null,
          },
          evaluation: {
            end: '2021-10-11T08:55:36.469Z',
            start: '2021-10-11T08:50:36.469Z',
            timeframe: '',
          },
          message: '',
          project: 'dynatrace',
          result: '',
          service: 'items',
          stage: 'quality-gate',
          status: '',
          temporaryData: {
            distributor: {
              subscriptionID: '',
            },
          },
          test: {
            end: '',
            start: '',
          },
        },
        id: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
        source: 'shipyard-controller',
        specversion: '1.0',
        time: '2021-10-11T08:55:37.408Z',
        type: 'sh.keptn.event.evaluation.triggered',
        shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
        shkeptnspecversion: '0.2.3',
        plainEvent: {
          data: {
            deployment: {
              deploymentNames: null,
            },
            evaluation: {
              end: '2021-10-11T08:55:36.469Z',
              start: '2021-10-11T08:50:36.469Z',
              timeframe: '',
            },
            message: '',
            project: 'dynatrace',
            result: '',
            service: 'items',
            stage: 'quality-gate',
            status: '',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
            test: {
              end: '',
              start: '',
            },
          },
          id: '38d85901-1d46-43a8-8eae-c872d49ddd7c',
          source: 'shipyard-controller',
          specversion: '1.0',
          time: '2021-10-11T08:55:37.408Z',
          type: 'sh.keptn.event.evaluation.triggered',
          shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
          shkeptnspecversion: '0.2.3',
        },
        finished: false,
      },
      {
        traces: [
          {
            traces: [],
            data: {
              project: 'dynatrace',
              result: 'pass',
              service: 'items',
              stage: 'quality-gate',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '1d3e1634-791e-4fa1-8991-87e129f8dfc6',
            source: 'dynatrace-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:42.141Z',
            type: 'sh.keptn.event.get-sli.started',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
            plainEvent: {
              data: {
                project: 'dynatrace',
                result: 'pass',
                service: 'items',
                stage: 'quality-gate',
                status: 'succeeded',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: '1d3e1634-791e-4fa1-8991-87e129f8dfc6',
              source: 'dynatrace-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:42.141Z',
              type: 'sh.keptn.event.get-sli.started',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              shkeptnspecversion: '0.2.3',
              triggeredid: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
            },
          },
          {
            traces: [],
            data: {
              'get-sli': {
                end: '2021-10-11T08:55:36.469Z',
                indicatorValues: [
                  {
                    metric: 'response_time_p95',
                    success: true,
                    value: 10.219262258461157,
                  },
                  {
                    metric: 'error_rate',
                    success: true,
                    value: 0,
                  },
                  {
                    metric: 'throughput',
                    success: true,
                    value: 933,
                  },
                ],
                start: '2021-10-11T08:50:36.469Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'dynatrace',
              result: 'pass',
              service: 'items',
              stage: 'quality-gate',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '8f2a2ba2-7d50-4013-a0ec-0921591eb139',
            source: 'dynatrace-service',
            specversion: '1.0',
            time: '2021-10-11T08:55:48.996Z',
            type: 'sh.keptn.event.get-sli.finished',
            shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
            plainEvent: {
              data: {
                'get-sli': {
                  end: '2021-10-11T08:55:36.469Z',
                  indicatorValues: [
                    {
                      metric: 'response_time_p95',
                      success: true,
                      value: 10.219262258461157,
                    },
                    {
                      metric: 'error_rate',
                      success: true,
                      value: 0,
                    },
                    {
                      metric: 'throughput',
                      success: true,
                      value: 933,
                    },
                  ],
                  start: '2021-10-11T08:50:36.469Z',
                },
                labels: {
                  DtCreds: 'dynatrace',
                },
                project: 'dynatrace',
                result: 'pass',
                service: 'items',
                stage: 'quality-gate',
                status: 'succeeded',
                temporaryData: {
                  distributor: {
                    subscriptionID: '',
                  },
                },
              },
              id: '8f2a2ba2-7d50-4013-a0ec-0921591eb139',
              source: 'dynatrace-service',
              specversion: '1.0',
              time: '2021-10-11T08:55:48.996Z',
              type: 'sh.keptn.event.get-sli.finished',
              shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
              shkeptnspecversion: '0.2.3',
              triggeredid: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
            },
          },
        ],
        data: {
          deployment: '',
          'get-sli': {
            end: '2021-10-11T08:55:36.469Z',
            indicators: ['response_time_p95', 'error_rate', 'throughput'],
            sliProvider: 'dynatrace',
            start: '2021-10-11T08:50:36.469Z',
          },
          project: 'dynatrace',
          service: 'items',
          stage: 'quality-gate',
          temporaryData: {
            distributor: {
              subscriptionID: '',
            },
          },
        },
        id: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
        source: 'lighthouse-service',
        specversion: '1.0',
        time: '2021-10-11T08:55:39.642Z',
        type: 'sh.keptn.event.get-sli.triggered',
        shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
        shkeptnspecversion: '0.2.3',
        plainEvent: {
          data: {
            deployment: '',
            'get-sli': {
              end: '2021-10-11T08:55:36.469Z',
              indicators: ['response_time_p95', 'error_rate', 'throughput'],
              sliProvider: 'dynatrace',
              start: '2021-10-11T08:50:36.469Z',
            },
            project: 'dynatrace',
            service: 'items',
            stage: 'quality-gate',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'e143780a-85bb-47f8-8140-e1d7aa4922fc',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-11T08:55:39.642Z',
          type: 'sh.keptn.event.get-sli.triggered',
          shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
          shkeptnspecversion: '0.2.3',
        },
        finished: false,
      },
      {
        traces: [],
        data: {
          message: 'could not retrieve Webhook config: no webhook config found',
          project: 'dynatrace',
          result: 'fail',
          service: 'items',
          stage: 'quality-gate',
          status: 'errored',
          temporaryData: {
            distributor: {
              subscriptionID: '',
            },
          },
        },
        id: '0c50be08-98f4-4f79-9e80-9d62522a86b8',
        source: 'shipyard-controller',
        specversion: '1.0',
        time: '2021-10-11T08:55:43.019Z',
        type: 'sh.keptn.event.quality-gate.evaluation.finished',
        shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
        shkeptnspecversion: '0.2.3',
        triggeredid: 'd3501830-cf74-4704-93fa-de9e468c5106',
        plainEvent: {
          data: {
            message: 'could not retrieve Webhook config: no webhook config found',
            project: 'dynatrace',
            result: 'fail',
            service: 'items',
            stage: 'quality-gate',
            status: 'errored',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '0c50be08-98f4-4f79-9e80-9d62522a86b8',
          source: 'shipyard-controller',
          specversion: '1.0',
          time: '2021-10-11T08:55:43.019Z',
          type: 'sh.keptn.event.quality-gate.evaluation.finished',
          shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'd3501830-cf74-4704-93fa-de9e468c5106',
        },
      },
    ],
    data: {
      deployment: {
        deploymentNames: null,
      },
      evaluation: {
        end: '2021-10-11T08:55:36.469Z',
        start: '2021-10-11T08:50:36.469Z',
        timeframe: '',
      },
      project: 'dynatrace',
      service: 'items',
      stage: 'quality-gate',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        end: '',
        start: '',
      },
    },
    id: 'd3501830-cf74-4704-93fa-de9e468c5106',
    source: 'https://github.com/keptn/keptn/api',
    specversion: '1.0',
    time: '2021-10-11T08:55:36.471Z',
    type: 'sh.keptn.event.quality-gate.evaluation.triggered',
    shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
    shkeptnspecversion: '0.2.3',
    plainEvent: {
      data: {
        deployment: {
          deploymentNames: null,
        },
        evaluation: {
          end: '2021-10-11T08:55:36.469Z',
          start: '2021-10-11T08:50:36.469Z',
          timeframe: '',
        },
        project: 'dynatrace',
        service: 'items',
        stage: 'quality-gate',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
        test: {
          end: '',
          start: '',
        },
      },
      id: 'd3501830-cf74-4704-93fa-de9e468c5106',
      source: 'https://github.com/keptn/keptn/api',
      specversion: '1.0',
      time: '2021-10-11T08:55:36.471Z',
      type: 'sh.keptn.event.quality-gate.evaluation.triggered',
      shkeptncontext: 'cdc4cbaf-9fc8-45a8-8e8c-7b51e0f056e8',
      shkeptnspecversion: '0.2.3',
    },
    finished: false,
  },
];

const rootTracesMock: Trace[] = TestUtils.mapTraces(rootTraces);
const evaluationTracesMock: Trace[] = TestUtils.mapTraces(evaluationTraces);
const multipleEvaluationsTracesMock: Trace = TestUtils.mapTraces(multipleEvaluationsTraces)[0];

export {
  rootTracesMock as RootTracesMock,
  evaluationTracesMock as EvaluationTracesMock,
  multipleEvaluationsTracesMock as MultipleEvaluationTracesMock,
};
