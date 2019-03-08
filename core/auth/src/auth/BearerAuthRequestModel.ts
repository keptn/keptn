import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';
import { stringify } from 'querystring';

@ApiModel({
  description: '',
  name: 'BearerAuthRequestModel',
})
export class BearerAuthRequestModel {
  @ApiModelProperty({
    description: 'Token value',
    example: ['atokenvalue'],
    required: true,
  })
  token: string;

  static isBearerAuthRequestModel(authRequest: BearerAuthRequestModel): boolean {
    return authRequest.token !== undefined &&
      typeof(authRequest.token) === 'string';
  }
}
