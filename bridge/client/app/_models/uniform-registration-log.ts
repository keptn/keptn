export class UniformRegistrationLog {
  integrationid: string;
  message: string;
  time: Date;
  shkeptncontext?: string;
  task?: string;
  triggeredid?: string;
}

export class UniformRegistrationLogResponse {
  logs: UniformRegistrationLog[];
}
