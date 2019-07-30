<template>
  <div>
    <div>
      <b-table
        class="mt-2"
        striped
        hover
        :fields="fields"
        :items="traces"
        bordered
        v-if="traces && traces.length"
      >
        <template slot="row-details" slot-scope="row">
          <small>{{row.item.message}}</small>
        </template>
      </b-table>
      <div class="placeholder" v-else>
        <font-awesome-icon class="icon" icon="arrow-left"/>
        <div class="placeholder-text">Please select an entry point!</div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import moment from 'moment';

export default {
  name: 'TraceList',
  props: {},
  data() {
    return { fields: ['timestamp', 'service'] };
  },
  filters: {
    moment: function formatDate(date) {
      return moment(date).format('YYYY-MM-DD, hh:mm:ss');
    },
  },
  computed: mapState({
    traces: state =>
      state.traces.map(trace => {
        return {
          timestamp: moment(trace._source['@timestamp']).format(
            'YYYY-MM-DD, hh:mm:ss',
          ),
          service: trace._source.keptnService,
          message: trace._source.message,
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
