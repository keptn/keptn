import { Deployment, StageDeployment } from './deployment';
import {
  ServiceDeploymentMock,
  ServiceDeploymentWithApprovalMock,
} from '../../../shared/fixtures/service-deployment-response.mock';
import { ServiceRemediationInformation } from './service-remediation-information';
import {
  ExpectedDeploymentMock,
  MergedSubSequencesDeliveryRollbackMock,
  ServiceRemediationInformationDevWithRemediationMock,
  ServiceRemediationInformationProductionWithRemediationMock,
  StageDeploymentDeliveryFinishedPassMock,
  StageDeploymentEmptyMock,
  StageDeploymentRollbackFinishedPassMock,
  SubSequencesFailedAndPassedMock,
  SubSequencesPassedLoadingMock,
  SubSequencesPassedMock,
  SubSequencesWarningFailedMock,
  SubSequencesWarningMock,
  UpdatedDeploymentMock,
} from '../_services/_mockData/deployments.mock';
import { SequenceState } from '../../../shared/models/sequence';

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

  it('should be finished', () => {
    const deployment = Deployment.fromJSON(ExpectedDeploymentMock);
    expect(deployment.isFinished()).toBe(true);
  });

  it('should not be finished', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    expect(deployment.isFinished()).toBe(false);
  });

  it('should assign subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON(StageDeploymentEmptyMock);
    const newStageDeployment = StageDeployment.fromJSON(StageDeploymentDeliveryFinishedPassMock);
    stageDeployment.update(newStageDeployment);
    expect(stageDeployment.subSequences).toEqual(newStageDeployment.subSequences);
  });

  it('should add subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON(StageDeploymentDeliveryFinishedPassMock);
    const newStageDeployment = StageDeployment.fromJSON(StageDeploymentRollbackFinishedPassMock);
    const expectedSubSequences = MergedSubSequencesDeliveryRollbackMock;
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
      ServiceRemediationInformationDevWithRemediationMock
    );
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].openRemediations).toEqual([]);
    expect(deployment.stages[1].openRemediations).toEqual([]);
    expect(deployment.stages[2].openRemediations).toEqual([]);
  });

  it('should update open remediations', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    deployment.stages[2].openRemediations = [];
    const serviceRemediationInformation = ServiceRemediationInformation.fromJSON(
      ServiceRemediationInformationProductionWithRemediationMock
    );
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].openRemediations).toEqual([]);
    expect(deployment.stages[1].openRemediations).toEqual([]);
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
    stageDeployment.subSequences = SubSequencesFailedAndPassedMock;
    expect(stageDeployment.isFaulty()).toBe(true);
  });

  it('should not be faulty', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassedLoadingMock;
    expect(stageDeployment.isFaulty()).toBe(false);
  });

  it('should not be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassedLoadingMock;
    expect(stageDeployment.isSuccessful()).toBe(false);
  });

  it('should be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesPassedMock;
    expect(stageDeployment.isSuccessful()).toBe(true);
  });

  it('should not be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesWarningFailedMock;
    expect(stageDeployment.isWarning()).toBe(false);
  });

  it('should be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = SubSequencesWarningMock;
    expect(stageDeployment.isWarning()).toBe(true);
  });

  it('should be aborted', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.state = SequenceState.ABORTED;
    expect(stageDeployment.isAborted()).toBe(true);
  });

  it('should not be aborted', () => {
    const stageDeployment = getStageDeployment();
    expect(stageDeployment.isAborted()).toBe(false);
  });

  function getStageDeployment(): StageDeployment {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return StageDeployment.fromJSON(ServiceDeploymentWithApprovalMock.stages[1]);
  }
});
