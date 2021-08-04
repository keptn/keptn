import { DtIconType } from '@dynatrace/barista-icons';

export const EVENT_ICONS: {[key: string]: DtIconType} = {
  'artifact-delivery': 'duplicate',
  delivery: 'duplicate',
  deployment: 'deploy',
  test: 'perfromance-health',
  evaluation: 'traffic-light',
  problem: 'criticalevent',
  remediation: 'criticalevent',
  release: 'hops',
  approval: 'unknown',
  waiting: 'idle',
  default: 'information'
};
