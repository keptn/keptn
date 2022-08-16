import { SequenceExecutionResult } from '../../../../shared/interfaces/sequence-execution-result';

const sequenceExecutionResult = {
  sequenceExecutions: [
    {
      scope: {
        keptnContext: '99a20ef4-d822-4185-bbee-0d7a364c213a',
        stage: 'dev',
      },
    },
  ],
} as SequenceExecutionResult;

export { sequenceExecutionResult as SequenceExecutionResultMock };
