export interface DynatraceRequestModel {
  State: string;
  ProblemID: string;
  ProblemTitle: string;
  ProblemDetails: ProblemDetails;
  ProblemImpact: string;
  ImpactedEntity: string;
  ImpactedEntities: ImpactedEntity[];
}

interface ImpactedEntity {
  type: string;
  name: string;
  entity: string;
}

interface ProblemDetails {
  id: string;
  startTime: number;
  endTime: number;
  displayName: string;
  impactLevel: string;
  status: string;
  severityLevel: string;
  commentCount: number;
  tagsOfAffectedEntities: TagsOfAffectedEntity[];
  rankedEvents: RankedEvent[];
  rankedImpacts: RankedImpact[];
  affectedCounts: AffectedCounts;
  recoveredCounts: AffectedCounts;
  hasRootCause: boolean;
}

interface AffectedCounts {
  INFRASTRUCTURE: number;
  SERVICE: number;
  APPLICATION: number;
  ENVIRONMENT: number;
}

interface RankedImpact {
  entityId: string;
  entityName: string;
  severityLevel: string;
  impactLevel: string;
  eventType: string;
}

interface RankedEvent {
  startTime: number;
  endTime: number;
  entityId: string;
  entityName: string;
  severityLevel: string;
  impactLevel: string;
  eventType: string;
  status: string;
  severities: Severity[];
  isRootCause: boolean;
  serviceMethodGroup: string;
  affectedRequestsPerMinute: number;
  userDefinedFailureRateThreshold: number;
  service: string;
}

interface Severity {
  context: string;
  value: number;
  unit: string;
}

interface TagsOfAffectedEntity {
  context: string;
  key: string;
  value: string;
}
