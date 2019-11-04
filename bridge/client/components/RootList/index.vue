<template>
  <div>
    <a
      href="#"
      v-for="root in roots"
      v-bind:key="root.keptnContext"
      v-bind:class="{ active: isActive(root.keptnContext) }"
      class="cardlink"
      @click="loadTraces(root.keptnContext)"
    >
      <b-card class="m-2">
        <b-card-text>
          <small>
            {{root.eventTypeHeadline}}
            <hr>
            <p class="mb-1">
              <b>Project:</b> {{ root.data.project }}
              <br>
              <b>Service:</b> {{ root.data.service }}
              <br>
              <b>Timestamp:</b> {{root.timestamp | moment}}
              <br>
              <b>KeptnContext:</b> {{root.keptnContext}}
            </p>
          </small>
          {{root.message}}
        </b-card-text>
      </b-card>
    </a>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import moment from 'moment';

export default {
  name: 'RootList',
  props: {},
  mounted() {
    this.$store.dispatch('reset');
    this.$store.dispatch('fetchRoots');
  },
  filters: {
    moment: function formatDate(date) {
      return moment(date).format('YYYY-MM-DD, hh:mm:ss');
    },
  },

  methods: {
    isActive(contextId) {
      return this.$store.state.currentContextId === contextId;
    },
    loadTraces(contextId) {
      return this.$store.dispatch('fetchTraces', contextId);
    },
  },

  computed: mapState({
    roots: state => state.roots,
    keptnContext: function keptnContext() {
      return this.$route.params.keptnContext;
    },
  }),

  watch: {
    keptnContext: function watchContext(keptnContext) {
      if (keptnContext) {
        this.$store.dispatch('findRoots', keptnContext);
        this.$store.dispatch('fetchTraces', keptnContext);
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="less">
a.cardlink {
  text-decoration: none !important;
  .card {
    color: #000000;
    text-decoration: none !important;
  }

  &:hover {
    .card {
      color: #000000;
      background-color: #eeeeee;
    }
  }

  &.active {
    .card {
      color: #ffffff;
      background-color: #006bb8;
    }
  }
}
</style>
