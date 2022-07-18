import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbDangerZoneComponent } from '../../client/app/_components/ktb-danger-zone/ktb-danger-zone.component';
import { moduleMetadata } from '@storybook/angular';
import { KtbDangerZoneModule } from '../../client/app/_components/ktb-danger-zone/ktb-danger-zone.module';
import { DeleteType } from '../../client/app/_interfaces/delete';

export default {
  title: 'Components/Danger Zone',
  component: KtbDangerZoneComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbDangerZoneModule],
    }),
  ],
} as Meta;

const template: Story<KtbDangerZoneComponent> = (args: KtbDangerZoneComponent) => ({
  props: args,
});

export const project = template.bind({});
project.args = {
  data: {
    name: 'project-1',
    type: DeleteType.PROJECT,
  },
  openDeletionDialog: (): void => alert('Deleting project'),
};

export const service = template.bind({});
service.args = {
  data: {
    name: 'service-1',
    type: DeleteType.SERVICE,
  },
  openDeletionDialog: (): void => alert('Deleting service'),
};
