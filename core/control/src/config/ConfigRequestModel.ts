import { KeptnGithubCredentials } from '../lib/types/KeptnGithubCredentials';
import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';

@ApiModel({
  description: '',
  name: 'ConfigRequestModel',
})
export class ConfigRequestModel {
  @ApiModelProperty({
    description: 'Object containing the required GitHub credentials',
    example: [{
      org: 'my_github_org',
      user: 'my_github_user',
      token: 'my_github_token',
    }],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  public data: KeptnGithubCredentials;
}
