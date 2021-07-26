export interface UniformRegistrationLog {
  integrationid: string;
  message: string;
  time: Date;
  shkeptncontext?: string;
  task?: string;
  triggeredid?: string;
}

export interface UniformRegistrationLogResponse {
  logs: UniformRegistrationLog[];
}
