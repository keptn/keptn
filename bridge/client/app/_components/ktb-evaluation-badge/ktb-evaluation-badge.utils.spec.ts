import { ResultTypes } from '../../../../shared/models/result-types';
import {
  EvaluationBadgeVariant,
  getEvaluationBadgeState,
  getEvaluationResultBadgeState,
} from './ktb-evaluation-badge.utils';
import { Trace } from '../../_models/trace';

describe('KtbEvaluationBadgeUtils', () => {
  describe('getEvaluationBadgeState', () => {
    it('should correctly set success status', () => {
      // given
      const evaluationResult = getEvaluationTrace(false, false, true);

      // when
      const state = getEvaluationBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: false,
        isWarning: false,
        isSuccess: true,
        score: undefined,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });

    it('should correctly set warning status', () => {
      // given
      const evaluationResult = getEvaluationTrace(false, true, false);

      // when
      const state = getEvaluationBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: false,
        isWarning: true,
        isSuccess: false,
        score: undefined,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });
    it('should correctly set error status', () => {
      // given
      const evaluationResult = getEvaluationTrace(true, false, false);

      // when
      const state = getEvaluationBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: true,
        isWarning: false,
        isSuccess: false,
        score: undefined,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });

    it('should correctly trim score', () => {
      // given
      const params = [
        { value: 0.56, expected: 0 },
        { value: 0.4, expected: 0 },
        { value: 0, expected: 0 },
        { value: 0.9, expected: 0 },
        { value: 1, expected: 1 },
        { value: undefined, expected: undefined },
      ];
      for (const param of params) {
        const evaluationResult = getEvaluationTrace(true, false, false, param.value);
        // when
        const state = getEvaluationBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

        // then
        expect(state).toEqual({
          isError: true,
          isWarning: false,
          isSuccess: false,
          score: param.expected,
          fillState: EvaluationBadgeVariant.NONE,
        });
      }
    });

    it('should correctly set fillState', () => {
      // given
      const params = [EvaluationBadgeVariant.FILL, EvaluationBadgeVariant.BORDER, EvaluationBadgeVariant.NONE];
      for (const param of params) {
        const evaluationResult = getEvaluationTrace(true, false, false, 0);

        // when
        const state = getEvaluationBadgeState(evaluationResult, param);

        // then
        expect(state).toEqual({
          isError: true,
          isWarning: false,
          isSuccess: false,
          score: 0,
          fillState: param,
        });
      }
    });

    function getEvaluationTrace(isError: boolean, isWarning: boolean, isSuccess: boolean, score?: number): Trace {
      const evaluation = Trace.fromJSON({});
      jest.spyOn(evaluation, 'isFaulty').mockReturnValue(isError);
      jest.spyOn(evaluation, 'isWarning').mockReturnValue(isWarning);
      jest.spyOn(evaluation, 'isSuccessful').mockReturnValue(isSuccess);
      jest
        .spyOn(evaluation, 'getEvaluationFinishedEvent')
        .mockReturnValue({ data: { evaluation: { score } } } as Trace);
      return evaluation;
    }
  });

  describe('getEvaluationResultBadgeState', () => {
    it('should correctly set success status', () => {
      // given
      const evaluationResult = {
        result: ResultTypes.PASSED,
        score: 0,
      };

      // when
      const state = getEvaluationResultBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: false,
        isWarning: false,
        isSuccess: true,
        score: 0,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });

    it('should correctly set warning status', () => {
      // given
      const evaluationResult = {
        result: ResultTypes.WARNING,
        score: 0,
      };

      // when
      const state = getEvaluationResultBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: false,
        isWarning: true,
        isSuccess: false,
        score: 0,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });
    it('should correctly set error status', () => {
      //  given
      const evaluationResult = {
        result: ResultTypes.FAILED,
        score: 0,
      };

      // when
      const state = getEvaluationResultBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

      // then
      expect(state).toEqual({
        isError: true,
        isWarning: false,
        isSuccess: false,
        score: 0,
        fillState: EvaluationBadgeVariant.NONE,
      });
    });

    it('should correctly trim score', () => {
      // given
      const params = [
        { value: 0.56, expected: 0 },
        { value: 0.4, expected: 0 },
        { value: 0, expected: 0 },
        { value: 0.9, expected: 0 },
        { value: 1, expected: 1 },
      ];
      for (const param of params) {
        const evaluationResult = {
          result: ResultTypes.FAILED,
          score: param.value,
        };
        // when
        const state = getEvaluationResultBadgeState(evaluationResult, EvaluationBadgeVariant.NONE);

        // then
        expect(state).toEqual({
          isError: true,
          isWarning: false,
          isSuccess: false,
          score: param.expected,
          fillState: EvaluationBadgeVariant.NONE,
        });
      }
    });

    it('should correctly set fillState', () => {
      // given
      const params = [EvaluationBadgeVariant.FILL, EvaluationBadgeVariant.BORDER, EvaluationBadgeVariant.NONE];
      for (const param of params) {
        const evaluationResult = {
          result: ResultTypes.FAILED,
          score: 0,
        };

        // when
        const state = getEvaluationResultBadgeState(evaluationResult, param);

        // then
        expect(state).toEqual({
          isError: true,
          isWarning: false,
          isSuccess: false,
          score: 0,
          fillState: param,
        });
      }
    });
  });
});
