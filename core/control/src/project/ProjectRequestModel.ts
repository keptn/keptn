import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';

@ApiModel({
  description: '',
  name: 'ProjectRequestModel',
})
export class ProjectRequestModel {
  @ApiModelProperty({
    description: 'Object containing the project information',
    example: [{
      data: {
        project: 'sockshop',
        stages: [
          {
            name: 'dev',
            deployment_strategy: 'direct',
          },
          {
            name: 'staging',
            deployment_strategy: 'blue_green_service',
          },
          {
            name: 'production',
            deployment_strategy: 'blue_green_service',
          },
        ],
      },
    }],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  public data: any;
}
