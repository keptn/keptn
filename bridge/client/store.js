/* eslint-disable no-param-reassign */
import Vue from 'vue';
import Vuex from 'vuex';

import api from './lib/api';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    roots: [],
    traces: [],
    remotecall: false,
    currentContextId: '',
  },

  actions: {
    remotecallStart({ commit }) {
      commit('remotecall', true);
    },

    remotecallEnd({ commit }) {
      commit('remotecall', false);
    },

    async fetchRoots({ commit, dispatch }) {
      dispatch('remotecallStart');
      try {
        const roots = await api.fetchRoots();
        commit('setRoots', roots);
      } finally {
        dispatch('remotecallEnd');
      }
    },

    reset({ commit }) {
      commit('setCurrentContextId', '');
      commit('setTraces', []);
      commit('setRoots', []);
    },

    async fetchTraces({ commit, dispatch }, contextId) {
      dispatch('remotecallStart');
      try {
        const traces = await api.fetchTraces(contextId);
        commit('setTraces', traces);
        commit('setCurrentContextId', contextId);
      } finally {
        dispatch('remotecallEnd');
      }
    },
    async findRoots({ commit, dispatch }, contextId) {
      dispatch('remotecallStart');
      try {
        const roots = await api.findRoots(contextId);
        commit('setRoots', roots);
      } finally {
        dispatch('remotecallEnd');
      }
    },
  },
  mutations: {
    setRoots(state, roots) {
      state.roots = roots;
    },
    setTraces(state, traces) {
      state.traces = traces;
    },
    remotecall(state, status) {
      state.remotecall = status;
    },
    setCurrentContextId(state, contextId) {
      state.currentContextId = contextId;
    },
  },
});
