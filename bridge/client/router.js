import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'tracelist',
      component: () => import(/* webpackChunkName: "traceslist" */ './components/TraceList'),
    },
    {
      path: '/view-context/:keptnContext',
      name: 'viewcontext',
      component: () => import(/* webpackChunkName: "traceslist" */ './components/TraceList'),
    },
  ],
});
