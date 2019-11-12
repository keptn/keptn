<template>
  <div>
    <div>
      <div v-if="traces && traces.length">
        <div class="traceHeader">
          <h2 v-if="traces.length > 1">{{traces[1].data.project + ':' + traces[1].data.service}}</h2>
          <p
            v-if="traces[0].type === 'sh.keptn.event.configuration.change'"
          >
            <b>New artifact:</b> {{traces[0].data.valuesCanary.image}}
          </p>
          <p
            v-if="traces[0].type === 'sh.keptn.event.problem.open'"
          >
            <b>Problem detected:</b> {{traces[0].data.problemtitle}}
          </p>
        </div>

        <div v-for="stage in getTracesPerStage(traces)" v-bind:key="stage.stageName" class="traceHeader">
          <h2 v-if="stage.stageName !== ''">Stage: {{stage.stageName}}</h2>
          <b-list-group >
            <div v-for="event in stage.events" v-bind:key="event.id" class="event-item">
              <b-list-group-item
                href="#"
                v-bind:class="{ active: isActive(event.id), error: isError(event), success: isSuccess(event), warning: isWarning(event) }"

                @click="activateEvent(event.id)"
              >
                <div class="d-flex w-100 justify-content-between">
                  <h5 class="mb-1">{{event.eventTypeHeadline}}</h5>
                  <small>{{ event.timestamp | moment }}</small>
                </div>
                <small>
                  <p class="mb-1">
                    <b>Project:</b> {{ event.data.project }}
                    <br>
                    <b>Service:</b> {{ event.data.service }}
                    <br>
                    <b>Stage:</b> {{ event.data.stage }}
                    <br>
                    <b>Source: </b> {{ event.source }}
                  </p>
                </small>

                <!-- EVENT SPECIFIC DETAILS -->
                <div v-if="event.type === 'sh.keptn.events.evaluation-done'">
                  <hr>
                  <small>
                    <b>Evaluation result: </b> {{ event.data.result }}
                    <br>
                    <div v-if="event.source === 'lighthouse-service'">
                      <b>Total score: </b> {{ event.data.evaluationdetails | totalScore }}
                    </div>
                  </small>
                  <div
                    v-if="isError(event) && event.source === 'lighthouse-service'">
                    <!--
                    <b>Violations:</b>
                    <div v-for="violation in getViolations(event.data.evaluationdetails)" :key="violation.indicatorId">
                      <div v-if="violation.type === 'upperSevere'">
                        <small><b>{{violation.indicatorId}}: </b> Measured value of <b>{{violation.actualValue}}</b> exceeded threshold of <b>{{violation.expectedValue}}</b></small>
                      </div>
                      <div v-if="violation.type === 'generic'">
                        <small><b>{{violation.indicatorId}}: </b> {{violation.reason}}</small>
                      </div>
                    </div>
                    -->
                  </div>
                </div>

                <div v-if="event.type === 'sh.keptn.events.tests-finished'">
                  <hr>
                  <small>
                    <b>Test strategy: </b> {{ event.data.teststrategy }}
                    <br>
                    <b>Duration: </b> {{ getDuration(event.timestamp, event.data.startedat) }}
                  </small>
                </div>

                <div v-if="event.type === 'sh.keptn.event.configuration.change' && !event.source.includes('remediation-service')">
                  <hr>
                  <small>
                    <b>Action: </b> {{ event.data.canary | canaryAction }}
                  </small>
                </div>

                <div v-if="event.type === 'sh.keptn.event.configuration.change' && event.source.includes('remediation-service')">
                  <hr>
                  <small>
                    <b>Action: </b> {{ event.data.deploymentChanges | remediationAction }}
                  </small>
                </div>
                <!-- EVENT SPECIFIC DETAILS -->
              </b-list-group-item>
              <div v-if="isActive(event.id)" class="event-detail">
                <small>
                  <b>Event payload:</b>
                  <br>
                  <vue-json-pretty
                    :path="'res'"
                    :data="JSON.parse(event.plainEvent)">
                  </vue-json-pretty>
                </small>
              </div>
            </div>
          </b-list-group>
        </div>
      </div>

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
import VueJsonPretty from 'vue-json-pretty';

export default {
  components: { VueJsonPretty },
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
          formatter: value => moment(value).format('YYYY-MM-DD, hh:mm:ss'),
        },
      ],
    };
  },

  filters: {
    moment: function formatDate(date) {
      return moment(date).format('YYYY-MM-DD, hh:mm:ss');
    },
    totalScore: function getTotalScore(evaluationDetails) {
      let totalScoreItem;
      if (evaluationDetails !== undefined) {
        if (evaluationDetails.hasOwnProperty('score') && evaluationDetails.indicatorResults !== null) {
          return evaluationDetails.score;
        }
      }
      return 'n/a (no evaluation performed by lighthouse service)';
    },
    canaryAction: function getCanaryAction(canary) {
      if (canary === undefined) {
        return 'n/a';
      }
      if (canary.hasOwnProperty('action')) {
        if (canary.action === 'promote') {
          return 'Promote to next stage';
        }
        if (canary.action === 'set') {
          if (canary.value !== undefined) {
            return `Settting traffic percentage for canary deployment to ${canary.value}`;
          }
          return 'Settting traffic percentage for canary deploymen';
        }
        if (canary.action === 'discard') {
          return 'Discarding deployment and reverting back to latest stable version';
        }
      }
    },
    remediationAction: function getRemediationAction(deploymentChanges) {
      if (deploymentChanges === undefined) {
        return 'n/a';
      }
      if (deploymentChanges.length > 0 && deploymentChanges[0].length > 0) {
        if (deploymentChanges[0].hasOwnProperty('propertyPath')) {
          if (deploymentChanges[0].propertyPath === undefined) {
            return 'n/a';
          }
          if (deploymentChanges[0].value === undefined) {
            return 'n/a';
          }

          return `Set property ${deploymentChanges[0].propertyPath} to ${deploymentChanges[0].value}`;
        }
      }
      return 'n/a';
    },
  },

  methods: {
    getTracesPerStage(traces) {
      const stages = [];

      traces.forEach((traceEvent) => {
        if (traceEvent.data !== undefined && traceEvent.data.stage !== undefined) {
          let stage = stages.find(stage => stage.stageName === traceEvent.data.stage);
          if (stage === undefined) {
            const newStage = {
              stageName: traceEvent.data.stage,
              events: [],
            };
            stages.push(newStage);
            stage = newStage;
          }
          stage.events.push(traceEvent);
        }
      });
      return stages;
    },
    isActive(contextId) {
      return this.$store.state.currentEventId === contextId;
    },
    isError(event) {
      return event.type === 'sh.keptn.events.evaluation-done' && event.data.result === 'fail';
    },
    isSuccess(event) {
      return event.type === 'sh.keptn.events.evaluation-done' && event.data.result === 'pass';
    },
    isWarning(event) {
      return event.type === 'sh.keptn.events.evaluation-done' && event.data.result === 'warning';
    },
    activateEvent(contextId) {
      return this.$store.dispatch('activateEvent', contextId);
    },
    getDuration(endTime, startTime) {
      return moment.utc(moment(endTime).diff(moment(startTime))).format('HH:mm:ss');
    },
    getViolations(evaluationDetails) {
      const violationsResult = [];
      /*
      if (evaluationDetails !== undefined) {
        if (evaluationDetails.indicatorResults !== undefined) {
          if (evaluationDetails.indicatorResults.length === 0) return violationsResult;

          evaluationDetails.indicatorResults.forEach((result) => {

            const { score } = result;
            if (score > 0) return;

            const indicatorId = result.id;
            const { violations } = result;

            if (violations.length < 1) return;

            for (let i = 0; i < violations.length; i += 1) {
              const violation = violations[i];
              console.log(violation);
              const newViolation = {};

              newViolation.indicatorId = indicatorId;

              if (violation.hasOwnProperty('breach') && violation.hasOwnProperty('upperSevere')) {
                const upperSevereBreach = violation;

                if (upperSevereBreach !== undefined) {
                  const expectedValue = upperSevereBreach.threshold;
                  const actualValue = upperSevereBreach.value;
                  newViolation.type = 'upperSevere';
                  newViolation.expectedValue = expectedValue;
                  newViolation.actualValue = actualValue;
                } else {
                  newViolation.type = 'generic';
                  newViolation.reason = violation.breach;
                }
                violationsResult.push(newViolation);
              }
            }

          });
        }
      }
       */
      return violationsResult;
    },
  },

  computed: mapState({
    traces: state => state.traces.map(trace => ({
      ...trace,
      _showDetails: true,
    })),
  }),
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="less">
  &.active {
    color: #ffffff;
    background-color: #006bb8;
  }

  a {
    text-decoration: none !important;
  .card {
    color: #000000;
    text-decoration: none !important;
  }
  }


  .error {
    background-color: #cd5c5c;
    color: #ffffff;
  }

  .success {
    background-color: #8fbc8f;
    color: #ffffff;
  }

  .warning {
    background-color: #bcac1a;
    color: #ffffff;
  }

  .traceHeader {
    margin-top: 10px;
    padding: 20px;
  }

  .event-item {
    margin-bottom: 10px;
  }

  .event-detail {
    padding: 10px;

    border-radius: 3px;
    background-color: #eeeeee;
    color: #000000;
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
