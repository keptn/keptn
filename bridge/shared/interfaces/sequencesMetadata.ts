export interface ISequencesMetadata {
  deployments: SequenceMetadataDeployment[];
  filter: {
    stages: string[];
    services: string[];
  };
}

export type SequenceMetadataDeployment = { service: string; stage: string; image: string };
