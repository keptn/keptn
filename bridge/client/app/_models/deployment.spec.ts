import { Deployment, StageDeployment } from './deployment';
import {
  ServiceDeploymentMock,
  ServiceDeploymentWithApprovalMock,
} from '../../../shared/fixtures/service-deployment-response.mock';
import { ServiceRemediationInformation } from './service-remediation-information';
import {
  ExpectedDeploymentMock,
  MergedSubSequencesDeliveryRollback,
  ServiceRemediationInformationDevWithRemediation,
  ServiceRemediationInformationProductionWithRemediation,
  StageDeploymentDeliveryFinishedPass,
  StageDeploymentEmpty,
  StageDeploymentRollbackFinishedPass,
  SubSequencesFailedAndPassed,
  SubSequencesPassed,
  SubSequencesPassedLoading,
  SubSequencesWarning,
  SubSequencesWarningFailed,
  UpdatedDeploymentMock,
} from '../_services/_mockData/deployments.mock';

describe('Deployment', () => {
  it('should correctly create new class', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    expect(deployment).toBeInstanceOf(Deployment);
    expect(deployment.stages[0]).toBeInstanceOf(StageDeployment);
    expect(deployment.stages[1]).toBeInstanceOf(StageDeployment);
  });

  it('should correctly update', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    const newDeployment = Deployment.fromJSON(UpdatedDeploymentMock);
    const expectedDeployment = Deployment.fromJSON(ExpectedDeploymentMock);
    deployment.update(newDeployment);

    expect(deployment).toEqual(expectedDeployment);
  });

  it('should assign subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON(StageDeploymentEmpty);
    const newStageDeployment = StageDeployment.fromJSON(StageDeploymentDeliveryFinishedPass);
    stageDeployment.update(newStageDeployment);
    expect(stageDeployment.subSequences).toEqual(newStageDeployment.subSequences);
  });

  it('should add subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON(StageDeploymentDeliveryFinishedPass);
    const newStageDeployment = StageDeployment.fromJSON(StageDeploymentRollbackFinishedPass);
    const expectedSubSequences = MergedSubSequencesDeliveryRollback;
    stageDeployment.update(newStageDeployment);
    expect(stageDeployment.subSequences).toEqual(expectedSubSequences);
  });

  it('should return latest time updated', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    expect(deployment.latestTimeUpdated).toEqual(new Date('2021-10-13T10:54:43.315Z'));
  });

  it('should remove open remediations', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    const serviceRemediationInformation = ServiceRemediationInformation.fromJSON(
      ServiceRemediationInformationDevWithRemediation
    );
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].remediationConfig).toBeUndefined();
    expect(deployment.stages[0].openRemediations).toEqual([]);

    expect(deployment.stages[1].remediationConfig).toBeUndefined();
    expect(deployment.stages[1].openRemediations).toEqual([]);

    expect(deployment.stages[2].remediationConfig).toBeUndefined();
    expect(deployment.stages[2].openRemediations).toEqual([]);
  });

  it('should update open remediations', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    deployment.stages[2].remediationConfig = undefined;
    deployment.stages[2].openRemediations = [];
    const serviceRemediationInformation = ServiceRemediationInformation.fromJSON(
      ServiceRemediationInformationProductionWithRemediation
    );
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].remediationConfig).toBeUndefined();
    expect(deployment.stages[0].openRemediations).toEqual([]);

    expect(deployment.stages[1].remediationConfig).toBeUndefined();
    expect(deployment.stages[1].openRemediations).toEqual([]);

    expect(deployment.stages[2].remediationConfig).toEqual(serviceRemediationInformation.stages[0].config);
    expect(deployment.stages[2].openRemediations).toEqual(serviceRemediationInformation.stages[0].remediations);
  });

  it('should remove approval', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.removeApproval();
    expect(stageDeployment.subSequences[0].hasPendingApproval).toBe(false);
    expect(stageDeployment.approvalInformation).toBeUndefined();
  });

  it('should be faulty', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesFailedAndPassed;
    expect(stageDeployment.isFaulty()).toBe(true);
  });

  it('should not be faulty', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassedLoading;
    expect(stageDeployment.isFaulty()).toBe(false);
  });

  it('should not be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassedLoading;
    expect(stageDeployment.isSuccessful()).toBe(false);
  });

  it('should be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassed;
    expect(stageDeployment.isSuccessful()).toBe(true);
  });

  it('should not be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesWarningFailed;
    expect(stageDeployment.isWarning()).toBe(false);
  });

  it('should be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesWarning;
    expect(stageDeployment.isWarning()).toBe(true);
  });

  function getStageDeployment(): StageDeployment {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return StageDeployment.fromJSON(ServiceDeploymentWithApprovalMock.stages[1]);
  }
});
