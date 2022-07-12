type JsonTypes = number | string | Record<string, unknown>;

export type JsonSerializable = undefined | null | JsonTypes | Array<JsonTypes>;
