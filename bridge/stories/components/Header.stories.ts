import { HttpClientModule } from '@angular/common/http';
import '@angular/localize/init';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterTestingModule } from '@angular/router/testing';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbAppHeaderComponent } from '../../client/app/_components/ktb-app-header/ktb-app-header.component';
import { KtbAppHeaderModule } from '../../client/app/_components/ktb-app-header/ktb-app-header.module';
import { Project } from '../../client/app/_models/project';

export default {
  title: 'Components/App Header',
  decorators: [
    moduleMetadata({
      imports: [KtbAppHeaderModule, HttpClientModule, RouterTestingModule, BrowserAnimationsModule],
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
} as Meta;

const template: Story<KtbAppHeaderComponent> = (args: KtbAppHeaderComponent) => ({
  props: args,
  template: `<ktb-header
    [info]="info"
    [projects]="projects"
    [metadata]="metadata"
    [selectedProject]="selectedProject"
  ></ktb-header>`,
});

export const standard = template.bind({});
standard.args = {
  projects: [{ projectName: 'pod-tato-head' } as Project, { projectName: 'sockshop' } as Project],
  info: {
    bridgeInfo: {
      featureFlags: {
        RESOURCE_SERVICE_ENABLED: true,
        D3_ENABLED: true,
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
