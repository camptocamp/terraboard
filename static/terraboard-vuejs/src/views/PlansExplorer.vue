<template>
  <div class="container-fluid mt-1">
    <form class="row justify-content-center" role="form">
      <div class="col-md-6 col-lg-5 col-xl-4 mb-3">
        <Multiselect
          id="tf_version"
          v-model="filters.lineage"
          :options="data.paths.options"
          :searchable="true"
          placeholder="State Path"
          @select="refreshList(1)"
          @clear="clearPath"
        >
        </Multiselect>
      </div>
      <div class="col-3 col-md-2 col-lg-1">
        <button class="clear btn btn-warning" v-on:click="resetSearch()">
          Reset
        </button>
      </div>
    </form>
  </div>
  <hr />
  <label id="navigate" v-if="results.length != 0">
    <span
      class="fas fa-caret-left"
      role="button"
      v-on:click="refreshList(pager.prevPage)"
      v-if="pager.prevPage"
    ></span>
    {{ pager.startItems }}-{{ pager.itemsInPage }}/{{ results.total }}
    <span
      class="fas fa-caret-right"
      role="button"
      v-on:click="refreshList(pager.nextPage)"
      v-if="pager.nextPage"
    ></span>
  </label>
  <div id="results" class="table-responsive">
    <table class="table table-border table-striped">
      <thead>
        <th></th>
        <th>Status</th>
        <th>Lineage</th>
        <th>TF Version</th>
        <th>Git Remote</th>
        <th>Git Commit</th>
        <th>CI Url</th>
        <th>Source</th>
        <th>Created At</th>
      </thead>
      <tbody>
        <tr v-for="r in results.plans" v-bind:key="r">
          <td>
            <router-link
              :to="
                `/lineage/${r.lineage_data.lineage}/plans?planid=${r.ID}`
              "
            >
              <span class="fas fa-link" aria-hidden="true"></span>
            </router-link>
          </td>
          <td>
            <div v-if="r.exit_code == 0"><i class="fas fa-check-circle text-success me-1"></i></div>
            <div v-else-if="r.exit_code == 1"><i class="fas fa-times-circle text-danger me-1"></i></div>
            <div v-else-if="r.exit_code == 2"><i class="fas fa-exclamation-circle text-warning me-1"></i></div>
          </td>
          <td>{{ r.lineage_data.lineage }}</td>
          <td>{{ r.terraform_version }}</td>
          <td>{{ r.git_remote }}</td>
          <td>{{ r.git_commit }}</td>
          <td>{{ r.ci_url }}</td>
          <td>{{ r.source }}</td>
          <td>{{ this.formatDate(r.CreatedAt) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Multiselect from "@vueform/multiselect";
import axios from "axios";
import router from "../router";

@Options({
  title: "Plans",
  components: {
    Multiselect,
  },
  emits: ["refresh"],
  data() {
    return {
      data: {
        paths: {
          options: [],
        },
      },
      filters: {
        lineage: null,
      },
      pager: {
        pages: 0,
        page: 0,
        prevPage: 0,
        nextPage: 0,
        startItems: 0,
        itemsInPage: 0,
        itemsPerPage: 20,
      },
      results: {},
    };
  },
  methods: {
    clearPath() {
      this.filters.lineage = null;
      this.refreshList();
    },
    resetSearch() {
      this.clearPath();
    },
    formatDate(date: string): string {
      return new Date(date).toUTCString();
    },
    fetchStates() {
      const url = `/api/lineages/stats`
      axios.get(url)
        .then((response) => {
          // handle success
          response.data.states.forEach((obj: any) => {
            let entry = {value: obj.lineage_value, label: obj.path}
            this.data.paths.options.push(entry)
          });
          this.refreshList();
        })
        .catch(function (err) {
          if (err.response) {
            console.log("Server Error:", err)
          } else if (err.request) {
            console.log("Network Error:", err)
          } else {
            console.log("Client Error:", err)
          }
        })
        .then(function () {
          // always executed
        });
    },
    refreshList(page: number) {
      let params: any = {};
      if (this.filters.lineage != null) {
        params.lineage = this.filters.lineage;
      }
      if (page != undefined) {
        params.page = page;
      } else {
        params.page = 1;
      }

      params.limit = this.pager.itemsPerPage;

      let query = Object.keys(params)
        .map(function(k) {
          return encodeURIComponent(k) + "=" + encodeURIComponent(params[k]);
        })
        .join("&");

      router.push({ name: "PlansExplorer", query: params });
      const url = `/api/plans/summary?` + query;
      axios
        .get(url)
        .then((response) => {
          this.results = response.data;
          this.pager.pages = Math.ceil(
            this.results.total / this.pager.itemsPerPage
          );
          this.pager.page = this.results.page;
          this.pager.prevPage =
            this.pager.page <= 1 
            ? undefined 
            : this.pager.page - 1;
          this.pager.nextPage =
            this.pager.page >= this.pager.pages
              ? undefined
              : this.pager.page + 1;
          this.pager.startItems =
            this.pager.itemsPerPage * (this.pager.page - 1) + 1;
          this.pager.itemsInPage = Math.min(
            this.pager.itemsPerPage * this.pager.page,
            this.results.total
          );
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
  },
  created() {
    this.updateTitle();
  },
  mounted() {
    this.fetchStates();
  },
})
export default class PlansExplorer extends Vue {}
</script>

<style lang="scss">
.truncate {
    overflow-x: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    -o-text-overflow: ellipsis;
    -ms-text-overflow: ellipsis;
}
</style>
