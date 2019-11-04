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
    currentEventId: '',
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
      commit('setCurrentEventId', '');
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
    activateEvent({ commit }, eventId) {
      commit('setCurrentEventId', eventId);
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
    setCurrentEventId(state, eventId) {
      if (state.currentEventId === eventId) {
        state.currentEventId = '';
      } else {
        state.currentEventId = eventId;
      }
    },
  },
});
