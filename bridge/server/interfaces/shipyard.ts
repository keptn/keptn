export interface Shipyard {
  apiVersion: string;
  kind: string;
  metadata: {
    name: string;
  };
  spec: {
    stages: IShipyardStage[];
  };
}

export interface IShipyardStage {
  name: string;
  sequences?: IShipyardSequence[];
}

export interface IShipyardSequence {
  name: string;
  triggeredOn?: [
    {
      event: string;
      selector?: {
        match: unknown;
      };
    }
  ];
  tasks: IShipyardTask[];
}

export interface IShipyardTask {
  name: string;
  properties: {
    deploymentstrategy?: 'direct' | 'blue_green_service';
    teststrategy?: 'performance' | 'functional';
  };
}
