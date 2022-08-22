import { IUniformSubscription, IUniformSubscriptionFilter } from '../../../shared/interfaces/uniform-subscription';
import { EventTypes } from '../../../shared/interfaces/event-types';

export function getEventContent(us: IUniformSubscription): string {
  return us.event.replace(EventTypes.PREFIX, '');
}

export function getPrefix(us: IUniformSubscription): string {
  return getEventContent(us).substring(0, getEventContent(us).lastIndexOf('.'));
}

export function getSuffix(us: IUniformSubscription): string {
  return getEventContent(us).split('.').pop() ?? '';
}

export function getFormattedEvent(us: IUniformSubscription): string {
  return us.event.replace('>', '*');
}

export function isGlobal(uf: IUniformSubscriptionFilter): boolean {
  return !uf.projects?.length;
}

export function hasProject(uf: IUniformSubscriptionFilter, projectName: string, includeEmpty = false): boolean {
  return uf.projects?.includes(projectName) || (includeEmpty && !uf.projects?.length);
}

export function getFirstStage(uf: IUniformSubscriptionFilter): string | undefined {
  return uf.stages?.find(() => true);
}

export function getFirstService(uf: IUniformSubscriptionFilter): string | undefined {
  return uf.services?.find(() => true);
}

export function getGlobalProjects(uf: IUniformSubscriptionFilter, status: boolean, projectName: string): string[] {
  if (status) {
    return [];
  }

  if (hasProject(uf, projectName)) {
    return uf.projects ? [...uf.projects] : [];
  }

  return uf.projects ? [...uf.projects, projectName] : [projectName];
}

export function formatFilter(uf: IUniformSubscriptionFilter, key: 'services' | 'stages' | 'projects'): string {
  return uf[key]?.join(', ') || 'all';
}

export function hasFilter(uf: IUniformSubscriptionFilter): boolean {
  return !!(uf.stages?.length || uf.services?.length);
}
