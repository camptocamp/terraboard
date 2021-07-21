<template>
  <div id="mainrow" class="row">
    <div id="leftcol" class="col-xl-4 col-xxl-3">
      <div class="mr-4">
        <div class="card">
          <h5 class="card-header">
            General Information
            <span
              v-if="isLocked(state.path)"
              class="float-right fas fa-lock"
              title="Locked by {{locks[state.path].Who}} on {{locks[state.path].Created | date:'medium'}} ({{locks[state.path].Operation}})"
            ></span>
          </h5>
          <ul class="list-group">
            <li class="list-group-item">
              <h5>
                <span class="fas fa-thumbtack" aria-hidden="true"></span>
                Version
              </h5>
              <select class="form-control" v-model="selectedVersion">
                <option
                  v-for="version in versions"
                  v-bind:key="version"
                  v-bind:value="version.versionId"
                >
                  {{ version.date }}
                </option>
              </select>
              <ul class="mt-2">
                <li class="list-group-item no-border-item">
                  Terraform version: {{ state.details.terraform_version }}
                </li>
                <li class="list-group-item no-border-item">
                  Serial: {{ state.details.serial }}
                </li>
              </ul>
            </li>
            <li class="list-group-item">
              <h5>
                <span class="fas fa-exchange-alt" aria-hidden="true"></span>
                Compare with version
              </h5>
              <div class="row justify-content-around">
                <select
                  class="form-control col-8"
                  v-model="compareVersion"
                >
                  <option
                    v-for="version in versions"
                    v-bind:key="version.versionId"
                    v-bind:value="version.versionId"
                  >
                    {{ version.date }}
                  </option>
                </select>
                <button
                  type="button"
                  class="btn btn-warning mt-2 col-5"
                  @click="
                    compareVersion = undefined;
                    display.compare = false;
                  "
                >
                  Reset
                </button>
              </div>
              <ul class="mt-2" v-if="display.compare && compare.stats">
                <li class="list-group-item no-border-item">
                  Terraform version: {{ compare.stats.to.terraform_version }}
                </li>
                <li class="list-group-item no-border-item">
                  Serial: {{ compare.stats.to.serial }}
                </li>
              </ul>
            </li>
          </ul>
        </div>
        <div id="nodes" class="card mt-4" v-if="display.details">
          <h5 class="card-header">Modules</h5>
          <ul id="nodeslist" class="list-group">
            <li class="list-group-item">
              <input
                id="resFilterInput"
                class="form-control"
                type="search"
                v-model="resFilter"
                placeholder="Filter resources..."
              />
            </li>
            <li
              class="list-group-item"
              v-for="mod in state.details.modules"
              v-bind:key="mod"
            >
              <div
                @click="display.mod = display.mod != mod ? mod : selectedMod"
                class="node-name"
                v-bind:class="{ selected: mod == selectedMod }"
              >
                <h4>{{mod.path ? mod.path : &quot;root&quot;}}</h4>
                <span class="badge float-right"
                  >{{resFilter == &quot;&quot; ? &apos;&apos; : filteredRes.length+&apos;/&apos;



                  }}{{ mod.resources.keys().length }}</span
                >
              </div>
              <!-- | filter = {name:resFilter} as filteredRes -->
              <ul v-show="display.mod == mod" class="list-group">
                <li
                  v-for="r in mod.resources"
                  v-bind:key="r"
                  v-bind:class="{ selected: r == selectedRes }"
                  @click="setSelected(mod, r)"
                  class="list-group-item resource"
                >
                  {{ r.type }}.{{ r.name }}{{ r.index }}
                </li>
                <li
                  v-bind:class="{ selected: state.outputs }"
                  v-if="mod.outputs.length &gt; 0"
                  @click="setOutputs(mod)"
                  class="list-group-item resource"
                >
                  Outputs
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div id="node" class="col-xl-8 col-xxl-9">
      <div class="row">
        <h1>{{ state.path }}</h1>
      </div>
      <StateDetails
        v-if="display.details && !display.outputs && !display.compare"
        v-bind:resource="selectedRes"
      />
      <StateOutputs
        v-if="display.details && display.outputs"
        v-bind:module="selectedMod"
      />
      <StatesCompare
        v-if="!display.details && display.compare"
        v-bind:compare="compare"
        v-bind:compareDiff="compareDiff"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import router from "../router";
import axios from "axios";
import hljs from "highlight.js";

import StateDetails from "../components/StateDetails.vue";
import StateOutputs from "../components/StateOutputs.vue";
import StatesCompare from "../components/StatesCompare.vue";

@Options({
  title: "States",
  components: {
    StateDetails,
    StateOutputs,
    StatesCompare,
  },
  data() {
    return {
      locks: {},
      versions: [],
      selectedVersion: "",
      compareVersion: "",
      selectedRes: {},
      selectedMod: {},
      resFilter: "",
      compare: {},
      compareDiff: {},
      url: {
        lineage: "",
        versionid: "",
        compare: "",
      },
      display: {
        details: true,
        compare: false,
        outputs: false,
        mod: {},
      },
      state: {
        details: {},
        path: {},
        outputs: false,
      },
    };
  },
  methods: {
    fetchLocks(): void {
      const url = `http://localhost:8080/api/locks`;
      axios
        .get(url)
        .then((response) => {
          this.locks = response.data;
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
    getVersions(): void {
      const url =
        `http://localhost:8080/api/lineage/activity/` + this.url.lineage;
      axios
        .get(url)
        .then((response) => {
          for (let i = 0; i < response.data.length; i++) {
            const version = {
              versionId: response.data[i].version_id,
              date: new Date(response.data[i].last_modified).toLocaleString(),
            };
            this.versions.unshift(version);
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
    isLocked(path: string): boolean {
      if (path in this.locks) {
        return true;
      }
      return false;
    },
    setSelected(mod: any, res: any): void {
      this.selectedMod = mod;
      this.selectedRes = res;
      this.state.outputs = false;
      this.display.outputs = false;
      var hash = res.type + "." + res.name;
      router.push({
        path: `/lineage/${this.url.lineage}`,
        query: { versionid: this.url.versionid, ressource: hash },
      });
    },
    setOutputs(mod: any): void {
      this.selectedMod = mod;
      this.selectedRes = null;
      this.state.outputs = true;
      this.display.outputs = true;
      router.push({
        path: `/lineage/${this.url.lineage}`,
        query: {
          versionid: this.url.versionid,
          ressource: mod.path + "." + "outputs",
        },
      });
    },
    getDetails(versionId: string) {
      if (versionId == undefined) {
        versionId = "";
      }
      const url =
        "http://localhost:8080/api/state/" +
        this.url.lineage +
        "?versionid=" +
        versionId +
        "#" +
        router.currentRoute.value.query.ressource;
      axios
        .get(url)
        .then((response) => {
          this.state.path = response.data["path"];
          this.state.details = response.data;
          this.selectedVersion = this.state.details.version.version_id;
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
    compareVersions(): void {
      if (
        this.compareVersion != undefined
      ) {
        router.replace({
          name: "State",
          params: { lineage: this.url.lineage },
          query: {
            versionid: this.url.versionid,
            compare: this.compareVersion.versionId,
          },
        });
        this.display.details = false;
        this.display.compare = true;

        const url =
          `http://localhost:8080/api/state/compare/` +
          this.url.lineage +
          "?from=" +
          this.selectedVersion +
          "&to=" +
          this.compareVersion;
        axios
          .get(url)
          .then((response) => {
            this.compare = response.data;

            this.compareDiff.only_in_old = Object.keys(
              this.compare.differences.only_in_old
            ).length;
            this.compareDiff.only_in_new = Object.keys(
              this.compare.differences.only_in_new
            ).length;
            this.compareDiff.differences = Object.keys(
              this.compare.differences.resource_diff
            ).length;
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
      } else {
        router.replace({
          name: "State",
          params: { lineage: this.url.lineage },
          query: { versionid: this.url.versionid },
        });
        this.display.details = true;
        this.display.compare = false;
      }
    },
  },
  watch: {
    "$data.state.details.modules": {
      handler: function(nv) {
        if (nv == undefined) {
          // Do not compute resource if no mods are known
          return;
        }

        // Sort the modules
        nv.sort(function(a: any, b: any) {
          return a.path.localeCompare(b.path);
        });

        for (let i = 0; i < nv.length; i++) {
          nv[i].resources.sort(function(a: any, b: any) {
            return a.name.localeCompare(b.name);
          });
        }
        this.selectedMod = nv[0];
        this.display.mod = nv[0];
        this.selectedRes = this.selectedMod.resources[0];

        if (router.currentRoute.value.query.ressource != undefined) {
          // Search for module in selected res
          let targetRes = router.currentRoute.value.query.ressource as string;
          for (let i = 0; i < nv.length; i++) {
            if (targetRes.startsWith(nv[i].path + ".")) {
              this.selectedMod = nv[i];
              this.display.mod = nv[i];
            }
          }

          let resources = this.selectedMod.resources;
          for (let j = 0; j < resources.length; j++) {
            if (targetRes == this.selectedMod.path + "." + resources[j].name) {
              this.selectedRes = resources[j];
              break;
            }
          }
        }
      },
    },
    "$data.selectedVersion": {
      handler: function(nv, ov) {
        router.replace({
          name: "State",
          params: { lineage: this.url.lineage },
          query: { versionid: nv },
        });
        if (nv != ov) {
          this.getDetails(nv);
        }
        this.compareVersions();
      },
    },
    "$data.compareVersion": {
      handler: function() {
        this.compareVersions();
      },
    },
  },
  created() {
    this.updateTitle();
  },
  mounted() {
    this.url.lineage = this.$route.params.lineage;
    this.url.versionid = router.currentRoute.value.query.versionid;
    this.selectedVersion = this.url.versionid;
    this.compareVersion = router.currentRoute.value.query.compare;
    this.fetchLocks();
    this.getVersions();
    this.getDetails(router.currentRoute.value.query.versionid);
    this.display.mod = this.selectedMod;
  },
  updated() {
    hljs.highlightAll();
  }
})
export default class State extends Vue {}
</script>

<style lang="scss">
#nodeslist .list-group-item .resource:hover,
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

#nodeslist .list-group-item .resource.selected {
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
