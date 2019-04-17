export interface DynatraceRequestModel {
  specversion: string;
  type: string;
  source: string;
  id: string;
  time: string;
  contenttype: string;
  data: Data;
  shkeptncontext: string;
}

interface Data {
  State: string;
  ProblemID: string;
  PID: string;
  ProblemTitle: string;
  ProblemDetails: ProblemDetails;
  ImpactedEntities: ImpactedEntity[];
  ImpactedEntity: string;
}

interface ImpactedEntity {
  type: string;
  name: string;
  entity: string;
}

interface ProblemDetails {
  id: string;
}
