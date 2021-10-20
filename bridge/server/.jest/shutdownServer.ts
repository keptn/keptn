const shutdown = async () => {
  await global.server?.close();
  process.exit();
}

export default shutdown();
