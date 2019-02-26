import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';

@ApiModel({
  description: '',
  name: 'ServiceRequestModel',
})
export class ServiceRequestModel {
  @ApiModelProperty({
    description: 'Object containing service information',
    example: [{
      data : {
        project: 'sockshop',
        file : 'deployment and service definition in YAML format',
      },
    }],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  public data: any;
}
