import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbExpandableTileComponent } from '../../client/app/_components/ktb-expandable-tile/ktb-expandable-tile.component';
import { KtbExpandableTileModule } from '../../client/app/_components/ktb-expandable-tile/ktb-expandable-tile.module';

export default {
  title: 'Components/Expandable Tile',
  component: KtbExpandableTileComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbExpandableTileModule, BrowserAnimationsModule, HttpClientModule],
    }),
  ],
} as Meta;

const template: Story = (args) => ({
  props: args,
  template: `<ktb-expandable-tile [alignment]="alignment" [expanded]="expanded">
        <span ktb-expandable-tile-header>Expandable tile header</span>
        This is the content of the tile
    </ktb-expandable-tile>`,
});

export const expanded = template.bind({});
expanded.args = {
  alignment: 'right',
  expanded: true,
};

export const collapsed = template.bind({});
collapsed.args = {
  alignment: 'right',
  expanded: false,
};
