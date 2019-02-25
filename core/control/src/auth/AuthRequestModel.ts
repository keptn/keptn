import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';

@ApiModel({
  description: '',
  name: 'AuthRequestModel',
})
export class ConfigRequestModel {
  @ApiModelProperty({
    description: 'Arbitrary JSON payload',
    example: [{
    }],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  public data: any;
}
