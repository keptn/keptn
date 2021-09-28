export interface Shipyard {
  apiVersion: string;
  kind: string;
  metadata: {
    name: string;
  };
  spec: {
    stages: [
      {
        name: string;
        sequences?: [
          {
            name: string;
            triggeredOn?: [
              {
                event: string;
                selector?: {
                  match: unknown;
                };
              }
            ];
            tasks: [
              {
                name: string;
                properties: {
                  deploymentstrategy?: 'direct' | 'blue_green_service';
                  teststrategy?: 'performance' | 'functional';
                };
              }
            ];
          }
        ];
      }
    ];
  };
}
