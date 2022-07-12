import { ServiceState } from './service-state';
import { ServiceStateResponse } from '../../../shared/fixtures/service-state-response.mock';
import { ServiceDeploymentWithApprovalMock } from '../../../shared/fixtures/service-deployment-response.mock';
import { Deployment } from './deployment';

describe('ServiceState', () => {
  it('should correctly create new class', () => {
    const serviceStateBasic = {
      name: 'carts',
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
      ],
    };
    const serviceState = ServiceState.fromJSON(serviceStateBasic);
    expect(serviceState).toBeInstanceOf(ServiceState);
  });

  it('should correctly update', () => {
    const serviceStates = ServiceStateResponse.map((sr) => ServiceState.fromJSON(sr));
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    serviceStates[0].deploymentInformation[0].deployment = deployment;
    const newServiceStatesBasic = [
      {
        deploymentInformation: [
          {
            // updated stage
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
            // updated stages
            stages: [
              {
                name: 'staging',
                time: '2021-11-05T10:49:01.288Z',
              },
              {
                name: 'production',
                time: '2021-10-13T11:01:18.567Z',
              },
            ],
            name: 'carts',
            image: 'carts',
            version: '0.12.3',
            keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
          },
          {
            // new deploymentInformation
            stages: [
              {
                name: 'production-A',
                time: '2021-11-10T12:21:33.991Z',
              },
            ],
            name: 'carts',
            image: 'carts',
            version: '0.12.1',
            keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd1',
          },
          // removed one deploymentInformation
        ],
        name: 'carts',
      },
      // removed serviceState
      // added serviceState
      {
        deploymentInformation: [],
        name: 'newService',
      },
    ];
    const newServiceStates = newServiceStatesBasic.map((st) => ServiceState.fromJSON(st));
    ServiceState.update(serviceStates, newServiceStates);
    newServiceStates[0].deploymentInformation[0].deployment = deployment; // cached deployment should not be overwritten
    expect(serviceStates).toEqual(newServiceStates);
  });

  it('should return latest image', () => {
    const serviceStateBasic = {
      name: 'carts',
      deploymentInformation: [
        {
          stages: [
            {
              name: 'dev',
              time: '2021-11-05T12:21:33.991Z',
            },
          ],
          name: 'carts',
          image: 'cartsImage1',
          version: '0.12.3',
          keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
        },
        {
          stages: [
            {
              name: 'staging',
              time: '2021-11-05T10:49:01.288Z',
            },
            {
              name: 'production',
              time: '2021-10-13T11:01:18.567Z',
            },
          ],
          name: 'cartsImage2',
          image: 'carts',
          version: '0.12.3',
          keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
        },
      ],
    };
    const serviceState = ServiceState.fromJSON(serviceStateBasic);
    expect(serviceState.getLatestImage()).toBe('cartsImage1:0.12.3');
  });

  it('should return "unknown" image', () => {
    const serviceStateBasic = {
      name: 'carts',
      deploymentInformation: [
        {
          stages: [
            {
              name: 'dev',
              hasOpenRemediations: false,
              time: '2021-11-05T12:21:33.991Z',
            },
          ],
          name: 'carts',
          version: '0.12.3',
          keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
        },
        {
          stages: [
            {
              name: 'staging',
              time: '2021-11-05T10:49:01.288Z',
            },
            {
              name: 'production',
              time: '2021-10-13T11:01:18.567Z',
            },
          ],
          name: 'cartsImage2',
          image: 'carts',
          version: '0.12.3',
          keptnContext: '29af69cc-ea85-4358-b169-ce29034d9c81',
        },
      ],
    };
    const serviceState = ServiceState.fromJSON(serviceStateBasic);
    expect(serviceState.getLatestImage()).toBe('unknown');
  });
});
