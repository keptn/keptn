<template>
  <div>
    <div>
      <b-table
        class="mt-2"
        striped
        hover
        dark
        :fields="fields"
        :items="traces"
        bordered
        v-if="traces && traces.length"
      ></b-table>
      <div class="placeholder" v-else>
        <font-awesome-icon class="icon" icon="arrow-left" />
        <div class="placeholder-text">Please select an entry point!</div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import moment from 'moment';

function formatDate(date) {
  return moment(date).format('YYYY-MM-DD, hh:mm:ss');
}

export default {
  name: 'TraceList',
  props: {},
  data() {
    return {
      fields: [
        { key: 'type', sortable: true },
        { key: 'data.project', label: 'Project', sortable: true },
        { key: 'data.service', label: 'Service', sortable: true },
        { key: 'data.stage', label: 'Stage', sortable: true },
        { key: 'data.tag', label: 'Tag', sortable: true },
        {
          key: 'timestamp',
          sortable: true,
          formatter: value => {
            return moment(value).format('YYYY-MM-DD, hh:mm:ss');
          },
        },
      ],
    };
  },
  computed: mapState({
    traces: state =>
      state.traces.map(trace => {
        return {
          ...trace,
          _showDetails: true,
        };
      }),
  }),
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="less">
a {
  text-decoration: none !important;
  .card {
    color: #000000;
    text-decoration: none !important;
  }
}
.placeholder {
  width: 100px;
  height: 100px;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 320px;
  right: 0;
  margin: auto;

  .icon {
    font-size: 130px;
    color: #aaaaaa;
    text-align: center;
    width: 300px;
  }

  .placeholder-text {
    color: #aaaaaa;
    text-align: center;
    width: 300px;
  }
}
</style>
