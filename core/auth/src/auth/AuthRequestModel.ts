import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';
import { stringify } from 'querystring';

@ApiModel({
  description: '',
  name: 'AuthRequestModel',
})
export class AuthRequestModel {
  @ApiModelProperty({
    description: 'Arbitrary JSON payload',
    example: ['{}'],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  payload: string;
  @ApiModelProperty({
    description: 'Signature',
    example: ['sha1=134hfdjslkfds'],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  signature: string;

  static isAuthRequestModel(authRequest: AuthRequestModel): boolean {
    console.log(authRequest.payload !== undefined);
    return authRequest.payload !== undefined &&
      typeof(authRequest.payload) === 'string' &&
      authRequest.signature !== undefined &&
      typeof(authRequest.signature) === 'string';
  }
}
