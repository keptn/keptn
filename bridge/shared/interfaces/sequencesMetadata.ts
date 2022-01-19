export interface ISequencesMetadata {
  deployments: SequenceMetadataDeployment[];
  filter: {
    stages: string[];
    services: string[];
  };
}

export type SequenceMetadataDeployment = {
  stage: {
    name: string;
    services: { name: string; image: string }[];
  };
};
