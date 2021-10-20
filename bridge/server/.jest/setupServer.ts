import { init } from '../app';

const setup = async () => {
  global.baseUrl = 'http://localhost/api/';

  const app = await init();
  app.set('port', 80);
  global.app = app;
  global.server = app.listen(80, '0.0.0.0');
};

export default await setup();
