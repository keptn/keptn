import { KeptnGithubCredentials } from '../KeptnGithubCredentials';
import { ApiModel, ApiModelProperty } from 'swagger-express-ts';

@ApiModel({
  description: '',
})
export class ConfigRequestModel {
  @ApiModelProperty({
    description: 'Object containing the required GitHub credentials',
    example: [{
      org: 'bar',
      user: 'test',
      token: 'token',
    }],
  })
  data: KeptnGithubCredentials;
}
