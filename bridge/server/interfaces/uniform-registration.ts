export interface UniformRegistration {
  id: string;
  metadata: {
    deplyomentname: string,
    distributorversion: string,
    hostname: string,
    integrationversion: string,
    kubernetesmetadata: {
      deploymentname: string,
      namespace: string,
      podname: string
    },
    location: string,
    status: string
  };
  name: string;
  subscription: {
    filter: {
      project: string,
      service: string,
      stage: string
    },
    status: string,
    topics: string []
  };
  unreadEvents?: number;
}

