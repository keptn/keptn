export interface Resource {
  resourceURI: string;
  resourceContent: string;
  metadata: {
    branch: string,
    upstreamURL: string,
    version: string
  };
}
