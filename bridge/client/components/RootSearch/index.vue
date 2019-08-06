<template>
  <div>
    <b-card class="m-2" bg-variant="dark" text-variant="light">
      <b-form @submit="onSubmit" @reset="onReset">
        <b-input-group>
          <b-form-input
            placeholder="keptn context"
            required
            minlength="6"
            v-model="form.keptnContext"
          ></b-form-input>
          <b-input-group-append>
            <b-button variant="info" type="submit">
              <font-awesome-icon icon="search"/>
            </b-button>
            <b-button variant="secondary" type="reset">
              <font-awesome-icon icon="undo"/>
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </b-form>
    </b-card>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'RootSearch',
  props: {},
  data() {
    return {
      form: {
        keptnContext: '',
      },
    };
  },

  methods: {
    async onReset() {
      this.form.keptnContext = '';
      this.$store.dispatch('reset');
      this.$store.dispatch('fetchRoots');
      this.$router.push({ name: 'tracelist' });
    },
    async onSubmit() {
      this.$store.dispatch('findRoots', this.form.keptnContext);
    },
    isActive(contextId) {
      return this.$store.state.currentContextId === contextId;
    },
  },

  computed: mapState({
    roots: state => state.roots,
  }),
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
