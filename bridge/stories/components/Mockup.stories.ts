import { FlexLayoutModule } from '@angular/flex-layout';
import '@angular/localize/init';
import { RouterTestingModule } from '@angular/router/testing';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbAppHeaderComponent } from '../../client/app/_components/ktb-app-header/ktb-app-header.component';
import { KtbAppHeaderModule } from '../../client/app/_components/ktb-app-header/ktb-app-header.module';
import { KtbSelectableTileModule } from '../../client/app/_components/ktb-selectable-tile/ktb-selectable-tile.module';
import { Project } from '../../client/app/_models/project';

export default {
  title: 'Components/Mockup',
  decorators: [
    moduleMetadata({
      imports: [KtbAppHeaderModule, KtbSelectableTileModule, FlexLayoutModule, RouterTestingModule],
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
} as Meta;

const toTemplate = (a: unknown): string => JSON.stringify(a).replace(/"/g, "'");

const template: Story<KtbAppHeaderComponent> = (args: KtbAppHeaderComponent) => ({
  props: args,
  template: `
    <div fxFlexFill fxLayout="column">
      <ktb-header
          [projects]="${toTemplate(args.projects)}"
          [info]="${toTemplate(args.info)}"
          [metadata]="${toTemplate(args.metadata)}"
      ></ktb-header>
      <div id="page-content" style="padding: 16px; margin: 0 auto; width: 800px; background-color: white">
          <div fxFlex fxLayoutGap="16px" fxLayout="row wrap">
            <ktb-selectable-tile
              *ngFor="let project of ${toTemplate(args.projects)}"
              style="width: 240px; padding-top: 16px"
            >
                <ktb-selectable-tile-header style="margin: 0 16px">{{project.projectName}}</ktb-selectable-tile-header>
                <div>
                  Stages:
                  <a *ngFor="let stage of project.stages" href="/" style="margin-right: 8px">{{stage.stageName}}</a>
                </div>
                <div>
                  Services:
                  <a *ngFor="let service of project.stages[0].services" href="/" style="margin-right: 8px">{{service.serviceName}}</a>
                </div>
            </ktb-selectable-tile>
          </div>
      </div>
    </div>
  `,
});

export const standard = template.bind({});
standard.args = {
  projects: [
    {
      projectName: 'pod-tato-head',
      stages: [
        { stageName: 'develop', services: [{ serviceName: 'pod-tato-svc' }] },
        { stageName: 'prod', services: [{ serviceName: 'pod-tato-svc' }] },
      ],
    } as Project,
    {
      projectName: 'sockshop',
      stages: [
        { stageName: 'develop', services: [{ serviceName: 'sockshop-svc' }] },
        { stageName: 'prod', services: [{ serviceName: 'sockshop-svc' }] },
      ],
    } as Project,
    {
      projectName: 'sockmock',
      stages: [
        { stageName: 'develop', services: [{ serviceName: 'sockmock-svc' }] },
        { stageName: 'prod', services: [{ serviceName: 'sockmock-svc' }] },
      ],
    } as Project,
    {
      projectName: 'lockmock',
      stages: [{ stageName: 'develop', services: [{ serviceName: 'lockmock-svc' }, { serviceName: 'lockshock-svc' }] }],
    } as Project,
  ],
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

export const single = template.bind({});
single.args = {
  projects: [
    {
      projectName: 'pod-tato-head',
      stages: [
        { stageName: 'develop', services: [{ serviceName: 'pod-tato-svc' }] },
        { stageName: 'prod', services: [{ serviceName: 'pod-tato-svc' }] },
      ],
    } as Project,
  ],
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
