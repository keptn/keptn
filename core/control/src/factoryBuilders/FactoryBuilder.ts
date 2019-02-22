import { CreateRequest } from '../types/createRequest';

export interface FactoryBuilder {
  createRepository(gitHubOrgName : string, payload : CreateRequest): Promise<void>;

  setHook(gitHubOrgName : string, payload : CreateRequest) : Promise<any>;
  
  initialCommit(gitHubOrgName : string, payload : CreateRequest) : Promise<any>;
  
  createBranchesForEachStages(gitHubOrgName : string, payload : CreateRequest) : Promise<any>;
  
  addShipyardToMaster(gitHubOrgName : string, payload : CreateRequest) : Promise<any>;
}