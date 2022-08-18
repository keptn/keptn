import semver from 'semver/preload';
import { IUniformRegistration } from '../../../shared/interfaces/uniform-registration';
import { hasProject } from './uniform-subscription';
import { IUniformSubscription } from '../../../shared/interfaces/uniform-subscription';

const preventSubscriptionUpdate = ['approval-service', 'remediation-service', 'lighthouse-service'];

export function getSubscriptions(ur: IUniformRegistration, projectName: string): IUniformSubscription[] {
  return ur.subscriptions.filter((subscription) => hasProject(subscription.filter, projectName, true));
}

export function hasSubscriptions(ur: IUniformRegistration, projectName: string): boolean {
  return ur.subscriptions.some((subscription) => hasProject(subscription.filter, projectName, true));
}

export function canEditSubscriptions(ur: IUniformRegistration): boolean {
  return !!(semver.valid(ur.metadata.distributorversion) && semver.gte(ur.metadata.distributorversion, '0.9.0'));
}

export function isChangeable(ur: IUniformRegistration): boolean {
  return !preventSubscriptionUpdate.includes(ur.name);
}
