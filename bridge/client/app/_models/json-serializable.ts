type JsonTypes = number | string | Record<string, unknown>;

export type JsonSerializable = undefined | JsonTypes | Array<JsonTypes>;
