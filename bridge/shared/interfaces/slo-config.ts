type scoreType = `${number}%`;
type constraint = `${'<' | '<=' | '=' | '' | '>' | '>='}${number}${'%' | ''}`;
interface Target {
  criteria: constraint[];
}

export interface ISloObjectives {
  sli?: string;
  displayName?: string;
  key_sli?: boolean;
  pass?: Target[];
  warning?: Target[];
  weight?: number;
}

export interface SloConfig {
  spec_version: string;
  comparison: {
    aggregate_function: 'avg' | 'p90' | 'p95';
    compare_with?: 'single_result' | 'several_results';
    include_result_with_score: 'pass' | 'pass_or_warn' | 'all';
    number_of_comparison_results: number;
  };
  filter: unknown;
  objectives: ISloObjectives[];
  total_score?: {
    pass?: scoreType;
    warning?: scoreType;
  };
}
