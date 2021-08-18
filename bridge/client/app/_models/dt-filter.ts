export type DtAutoComplete = { name: string, autocomplete: { name: string }[] };
export type DtFilter = DtAutoComplete[] | { name: string };
export type DtFilterArray = [DtAutoComplete, { name: string }];
