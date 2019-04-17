export interface KeptnOrgToRepoConfigMapModel {
  kind: string;
  apiVersion: string;
  metadata: Metadata;
  data: Data;
}

interface Data {
  orgsToRepos: string;
}

interface Metadata {
  name: string;
  namespace: string;
  selfLink: string;
  uid: string;
  resourceVersion: string;
  creationTimestamp: string;
  annotations: Annotations;
}

interface Annotations {
  'kubectl.kubernetes.io/last-applied-configuration': string;
}
