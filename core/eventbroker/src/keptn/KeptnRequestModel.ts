import { ApiModel, ApiModelProperty, SwaggerDefinitionConstant } from 'swagger-express-ts';
import moment from 'moment';

@ApiModel({
  description: '',
  name: 'KeptnRequestModel',
})
export class KeptnRequestModel {
  @ApiModelProperty({
    description: 'CE SpecVersion',
    example: ['0.2'],
    type: 'string',
    required: true,
  })
  public specversion: string;

  @ApiModelProperty({
    description: 'CE Type',
    example: [
      'sh.keptn.events.new-artefact',
      'sh.keptn.deployment-finished',
      'sh.keptn.tests-finished',
      'sh.keptn.evaluation-finished',
    ],
    type: 'string',
    required: true,
  })
  public type: string;

  @ApiModelProperty({
    description: 'CE Source',
    example: ['https://github-operator.svc.cluster.local'],
    type: 'string',
    required: true,
  })
  public source: string;

  @ApiModelProperty({
    description: 'CE Id',
    example: ['A234-1234-1234'],
    type: 'string',
    required: true,
  })
  public id: string;

  @ApiModelProperty({
    description: 'CE Time',
    example: ['2018-04-05T17:31:00Z'],
    type: 'string',
    required: true,
  })
  public time: string;

  @ApiModelProperty({
    description: 'CE Data-ContentType',
    example: ['application/json'],
    type: 'string',
    required: true,
  })
  public datacontenttype: string;

  @ApiModelProperty({
    description: 'Object containing the event payload',
    example: [{}],
    type: SwaggerDefinitionConstant.Model.Type.OBJECT,
    required: true,
  })
  public data: any;

  constructor() {
    this.specversion = '0.2';
    this.time = moment().format();
    this.datacontenttype = 'application/json';
  }
}
