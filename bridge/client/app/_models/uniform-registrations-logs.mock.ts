import { UniformRegistrationLog } from '../../../server/interfaces/uniform-registration-log';

const logs = [
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error on evaluation - failed',
    time: '2021-05-10T03:04:05.000Z',
    shkeptncontext: 'a6c9ec8b-2021-4797-875c-2693005312b8',
    task: 'evaluation',
    triggeredid: '96aeda8b-570b-49af-9a3d-d86dbc07678d'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error on evaluation - failed - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, ' +
      'quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ' +
      'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
    time: '2021-05-10T04:04:05.000Z',
    shkeptncontext: 'a6c9ec8b-2021-4797-875c-2693005312b8',
    task: 'evaluation',
    triggeredid: '96aeda8b-570b-49af-9a3d-d86dbc07678d'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Very long Error line 1 - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, ' +
      'quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ' +
      'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
    time: '2021-05-10T05:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 2',
    time: '2021-05-10T06:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 3',
    time: '2021-05-10T07:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 4',
    time: '2021-05-10T08:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 5',
    time: '2021-05-10T09:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 6',
    time: '2021-05-10T09:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 7',
    time: '2021-05-10T10:04:05.000Z'
  },
  {
    integrationid: '6d3190bed8866ebd90ec3d12875e890802d08d47',
    message: 'Error line 8',
    time: '2021-05-10T11:04:05.000Z'
  }
];
const UniformRegistrationLogsMock: UniformRegistrationLog[] = JSON.parse(JSON.stringify(logs));
export {UniformRegistrationLogsMock};
