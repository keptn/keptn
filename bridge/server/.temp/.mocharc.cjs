module.exports = {
  extension: ['ts'],
  spec: ["tests/**.ts"],
  // loader: 'ts-node/esm',
  // require: 'esm',
  // require: 'ts-node/register',
  'node-option': ['loader=ts-node/esm','experimental-specifier-resolution=node']
};
