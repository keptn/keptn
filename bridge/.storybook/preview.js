import {setCompodocJson} from "@storybook/addon-docs/angular";
import docJson from "../documentation.json";
import '!style-loader!css-loader!sass-loader!./scss-loader.scss';
import {componentWrapperDecorator} from "@storybook/angular";

setCompodocJson(docJson);

export const parameters = {
  layout: 'centered',
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
  docs: { inlineStories: true },
  storySort: {
    order: ['Introduction']
  },
};

export const decorators = [
  componentWrapperDecorator((story) =>
    `<link href="/assets/default-branding/keptn-theme.css" rel="stylesheet"/>${story}`),
];
