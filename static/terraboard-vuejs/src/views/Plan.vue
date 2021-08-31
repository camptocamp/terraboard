<template>
  <div id="mainrow" class="row">
    <div id="leftcol" class="col-xl-4 col-xxl-3">
      <div class="mr-4">
        <router-link class="ms-2" :to="`/lineage/${url.lineage}`"><i class="fas fa-arrow-left"></i> Back to workspace</router-link>
        <div id="nodes" class="card mt-2">
          <h5 class="card-header">Plans</h5>
              <ul id="nodeslist" class="list-group m-3">
                <li
                  v-for="plan in plans"
                  v-bind:key="plan"
                  v-bind:class="{ selected: plan == selectedPlan }"
                  @click="setPlanSelected(plan)"
                  class="list-group-item plan"
                >
                  {{ this.formatDate(plan.CreatedAt) }}
                </li>
              </ul>
        </div>
      </div>
    </div>
    <div id="node" class="col-xl-8 col-xxl-9">
      <PlanContent
        v-if="selectedPlan.parsed_plan !== undefined"
        v-bind:plan="selectedPlan"
        v-bind:key="selectedPlan"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import router from "../router";
import axios from "axios";
import hljs from "highlight.js";

import PlanContent from "../components/PlanContent.vue";

@Options({
  title: "States",
  components: {
    PlanContent,
  },
  emits: ["refresh"],
  data() {
    return {
      selectedPlan: {},
      url: {
        lineage: "",
        planid: "",
      },
      state: {
        details: {},
        path: {},
        outputs: false,
      },
      plans: [],
    };
  },
  methods: {
    formatDate(date: string): string {
        return new Date(date).toUTCString();
    },
    fetchLatestPlans(limit: number): void {
      const url = `/api/plans?limit=`+limit+`&lineage=`+this.url.lineage;
      axios
        .get(url)
        .then((response) => {
          this.plans = response.data.plans;
          if (router.currentRoute.value.query.planid !== undefined && router.currentRoute.value.query.planid !== "") {
            this.url.planid = router.currentRoute.value.query.planid;

            let planFinded = false;
            this.plans.forEach((plan: any) => {
              if (plan.ID == this.url.planid) {
                planFinded = true;
                this.setPlanSelected(plan); 
              }
            });
            if (planFinded === false) {
              const url = `/api/plans?lineage=`+this.url.lineage;
              axios
                .get(url)
                .then((response) => {
                  response.data.plans.forEach((plan: any) => {
                    if (plan.ID == this.url.planid) {
                      this.setPlanSelected(plan); 
                    }
                  });
                })
                .catch(function(err) {
                  if (err.response) {
                    console.log("Server Error:", err);
                  } else if (err.request) {
                    console.log("Network Error:", err);
                  } else {
                    console.log("Client Error:", err);
                  }
                })
                .then(function() {
                  // always executed
                });
            }
          } else {
            this.url.planid = "";
            this.setPlanSelected(this.plans[0]); 
          }
        })
        .catch(function(err) {
          if (err.response) {
            console.log("Server Error:", err);
          } else if (err.request) {
            console.log("Network Error:", err);
          } else {
            console.log("Client Error:", err);
          }
        })
        .then(function() {
          // always executed
        });
    },
    setPlanSelected(plan: any): void {
      this.selectedPlan = plan;
      router.replace({
        path: `/lineage/${this.url.lineage}/plans`,
        query: { 
          planid: plan.ID,
        },
      });
    },
  },
  created() {
    this.updateTitle();
    this.url.lineage = this.$route.params.lineage;
    this.url.planid = router.currentRoute.value.query.planid;
    this.fetchLatestPlans(10);
  },
  updated() {
    hljs.highlightAll();
  }
})
export default class Plan extends Vue {}
</script>

<style lang="scss">
#nodeslist .list-group-item {
  cursor: pointer;
}

#nodeslist .list-group-item .resource:hover,
#nodeslist .list-group-item.resource:hover,
#nodeslist .list-group-item.plan:hover,
#only-in-old .list-group-item:hover,
#only-in-new .list-group-item:hover {
  background-color: #d9edf7;
  background-image: none;
  color: #337ab7;
}

#nodeslist .list-group-item.selected {
  background-color: #d9edf7;
  color: #337ab7;
}

#nodeslist .list-group-item.selected {
  background-color: #d9edf7;
  color: #337ab7;
}

#nodeslist .list-group-item .fa-star,
.panel-title .fa-star {
  float: right;
  color: orange;
  display: none;
}

#nodeslist .list-group-item.starred .fa-star,
.panel-title.starred .fa-star {
  display: block;
}

#node .list-group {
  padding: 0 10px 0 10px;
}

#node .list-group.active {
  border-left: solid #4d90f0;
  padding-left: 7px;
}

#node .diff-stats .progress {
  text-align: center;
  width: 5em;
  display: inline;
  float: right;
}

#nodeslist .node-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

#nodeslist .progress {
  text-align: center;
  display: inline;
  float: right;
  margin-left: 5px;
}

.btn-file {
  position: relative;
  overflow: hidden;
}

.btn-file input[type=file] {
  position: absolute;
  top: 0;
  right: 0;
  min-width: 100%;
  min-height: 100%;
  font-size: 100px;
  text-align: right;
  filter: alpha(opacity=0);
  opacity: 0;
  outline: none;
  background: white;
  cursor: inherit;
  display: block;
}

.btn-checkbox {
  top: -2px;
  color: #777;
  padding: 15px 0 15px 20px;
}

#nodeslist .progress-bar {
  text-indent: -9999px; // Hide text to the left
}

#nodeslist .progress-bar:hover {
  text-indent: 0; // Reveal text
}

#nodeslist .progress-bar {
  float: right;
}

.resource-title {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.resource-title:hover {
  text-overflow: clip;
  overflow: auto;
}
</style>
