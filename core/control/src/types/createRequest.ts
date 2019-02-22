export interface CreateRequest {
  data : Data
}

export interface Data {
  application: string,
  stages: Stage[]
}

export interface Stage {
  name: string,
  deployment_strategy: string,
}