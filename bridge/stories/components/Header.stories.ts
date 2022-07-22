import '@angular/localize/init';
import { RouterTestingModule } from '@angular/router/testing';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbAppHeaderComponent } from '../../client/app/_components/ktb-app-header/ktb-app-header.component';
import { KtbAppHeaderModule } from '../../client/app/_components/ktb-app-header/ktb-app-header.module';
import { Project } from '../../client/app/_models/project';

export default {
  title: 'Components/App Header',
  component: KtbAppHeaderComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbAppHeaderModule, RouterTestingModule],
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
} as Meta;

const template: Story<KtbAppHeaderComponent> = (args: KtbAppHeaderComponent) => ({
  props: args,
});

export const standard = template.bind({});
standard.args = {
  projects: [{ projectName: 'pod-tato-head' } as Project, { projectName: 'sockshop' } as Project],
  info: {
    bridgeInfo: {
      featureFlags: {
        RESOURCE_SERVICE_ENABLED: true,
        D3_HEATMAP_ENABLED: true,
      },
      cliDownloadLink: '',
      enableVersionCheckFeature: true,
      showApiToken: true,
      authType: '',
    },
  },
  metadata: {
    namespace: 'keptn',
    keptnversion: '0.18.0',
    keptnlabel: '',
    bridgeversion: '0.18.0',
    shipyardversion: '2',
  },
  selectedProject: 'sockshop',
  changeProject: (selectedProject: string | undefined): void => {
    (standard.args ?? {}).selectedProject = selectedProject;
  },
};
