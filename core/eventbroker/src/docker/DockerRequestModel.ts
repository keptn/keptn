import { ApiModel, ApiModelProperty } from 'swagger-express-ts';

interface DockerEvent {
  id: string;
  timestamp: string;
  action: string;
  target: Target;
  request: Request;
  actor: Actor;
  source: Source;
}

interface Source {
  addr: string;
  instanceID: string;
}

interface Actor {
}

interface Request {
  id: string;
  addr: string;
  host: string;
  method: string;
  useragent: string;
}

interface Target {
  mediaType: string;
  size: number;
  digest: string;
  length: number;
  repository: string;
  url: string;
  tag: string;
}

@ApiModel({
  description: '',
  name: 'KeptnRequestModel',
})
export class DockerRequestModel {
  @ApiModelProperty({
    description: 'Events',
    example: [{
      events: [
        {
          id: 'a24e1fe3-efc9-42e3-b274-3c736f015552',
          timestamp: '2019-03-05T14:52:38.292839945Z',
          action: 'push',
          target: {
            mediaType: 'application/vnd.docker.distribution.manifest.v2+json',
            size: 2223,
            digest: 'sha256:d1b654481b04da5f1f69dc4e3bd72f4b592a60c6fb5618a9096eeac870cd3fe6',
            length: 2223,
            repository: 'keptn/keptn-event-broker-ext',
            url: 'http://docker-registry.keptn.svc.cluster.local:5000/',
            tag: 'latest',
          },
          request: {
            id: '164cafa5-ec03-4445-83b3-3ba6f7f28772',
            addr: '127.0.0.1:37402',
            host: 'docker-registry.keptn.svc.cluster.local:5000',
            method: 'PUT',
            useragent: 'kaniko/unset',
          },
          actor: {

          },
          source: {
            addr: 'docker-registry-55bd8d967c-hztw9:5000',
            instanceID: '2871c1e7-78b9-4fa5-b749-dc5e5fff8f9c',
          },
        },
      ],
    }],
    type: 'Object[]',
    required: true,
  })
  public events: DockerEvent[];
}
