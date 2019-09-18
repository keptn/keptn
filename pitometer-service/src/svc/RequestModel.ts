import { ApiModel, ApiModelProperty } from 'swagger-express-ts';
import moment from 'moment';

const uuidv4 = require('uuid/v4');

interface Data {
  githuborg: string;
  project: string;
  teststrategy: string;
  deploymentstrategy: string;
  stage: string;
  service: string;
  image: string;
  tag: string;
  [x: string]: any;
}

@ApiModel({
  description: '',
  name: 'RequestModel',
})
export class RequestModel {

  static EVENT_TYPES = {
    NEW_ARTIFACT: 'sh.keptn.events.new-artifact',
    CONFIGURATION_CHANGED: 'sh.keptn.events.configuration-changed',
    DEPLOYMENT_FINISHED: 'sh.keptn.events.deployment-finished',
    TESTS_FINISHED: 'sh.keptn.events.tests-finished',
    EVALUATION_DONE: 'sh.keptn.events.evaluation-done',
  };

  @ApiModelProperty({
    description: 'specversion',
    example: ['0.2'],
    type: 'string',
    required: true,
  })
  specversion: string;

  @ApiModelProperty({
    description: 'type',
    example: ['sh.keptn.events.tests-finished'],
    type: 'string',
    required: true,
  })
  type: string;

  @ApiModelProperty({
    description: 'source',
    example: ['Keptn'],
    type: 'string',
    required: true,
  })
  source: string;

  @ApiModelProperty({
    description: 'id',
    example: ['1234'],
    type: 'string',
    required: true,
  })
  id: string;

  @ApiModelProperty({
    description: 'time',
    example: ['20190325-15:25:56.096'],
    type: 'string',
    required: true,
  })
  time: string;

  @ApiModelProperty({
    description: 'contenttype',
    example: ['application/json'],
    type: 'string',
    required: true,
  })
  contenttype: string;

  @ApiModelProperty({
    description: 'shkeptncontext',
    example: ['db51be80-4fee-41af-bb53-1b093d2b694c'],
    type: 'string',
    required: true,
  })
  shkeptncontext: string;

  @ApiModelProperty({
    description: 'data',
    example: ['db51be80-4fee-41af-bb53-1b093d2b694c'],
    type: 'string',
    required: true,
  })
  data: Data;

  constructor() {
    this.id = uuidv4();
    this.specversion = '0.2';
    this.time = moment().format();
    this.contenttype = 'application/json';
    this.shkeptncontext = uuidv4();
    this.data = {} as Data;
  }
}
