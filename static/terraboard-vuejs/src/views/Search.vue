<template>
  <div class="container-fluid mt-1">
    <form class="row" role="form">
      <div class="col-md-4 col-lg-3 col-xl-2 mb-3">
        <Multiselect
          id="tf_version"
          v-model="search.tf_version"
          :options="data.tf_versions"
          :searchable="true"
          placeholder="Terraform Version"
          style="max-width: 300px;"
          @select="doSearch(1)"
          @clear="clearTfVersion"
        >
        </Multiselect>
      </div>
      <div class="col-md-4 col-lg-3 col-xl-2 mb-3">
        <Multiselect
          id="resource_type"
          v-model="search.resType"
          :options="data.resTypes"
          :searchable="true"
          placeholder="Resource Type"
          style="max-width: 300px;"
          @select="doSearch(1)"
          @clear="clearResourceType"
        >
        </Multiselect>
      </div>
      <div class="col-md-4 col-lg-3 col-xl-2 mb-3">
        <Multiselect
          id="resource_id"
          v-model="search.resID"
          :options="data.resIDs"
          :searchable="true"
          placeholder="Resource ID"
          style="max-width: 300px;"
          @select="doSearch(1)"
          @clear="clearResourceID"
        >
        </Multiselect>
      </div>
      <div class="col-md-4 col-lg-3 col-xl-2 mb-3">
        <Multiselect
          id="attribute_key"
          v-model="search.attrKey"
          :options="data.attrKeys"
          :searchable="true"
          placeholder="Attribute Key"
          style="max-width: 300px;"
          @select="doSearch(1)"
          @clear="clearAttributeKey"
        >
        </Multiselect>
      </div>
      <div class="col-md-4 col-lg-3 col-xl-2 mb-3">
        <input
          id="attribute_value"
          type="text"
          class="form-control"
          placeholder="Attribute Value"
          v-model="search.attrVal"
          v-on:change="doSearch(1)"
        />
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
      v-on:click="doSearch(pager.prevPage)"
      v-if="pager.prevPage"
    ></span>
    {{ pager.startItems }}-{{ pager.itemsInPage }}/{{ results.total }}
    <span
      class="fas fa-caret-right"
      v-on:click="doSearch(pager.nextPage)"
      v-if="pager.nextPage"
    ></span>
  </label>
  <div id="results">
    <table class="table table-border table-striped">
      <thead>
        <th></th>
        <th>Path</th>
        <th>TF Version</th>
        <th>Serial</th>
        <th>Module Path</th>
        <th>Resource ID</th>
        <th>Key</th>
        <th>Value</th>
      </thead>
      <tbody>
        <tr v-for="r in results.results" v-bind:key="r">
          <td>
            <router-link
              :to="
                `/lineage/${r.lineage_value}?versionid=${r.version_id}#${r.module_path}.${r.resource_name}`
              "
            >
              <span class="fas fa-link" aria-hidden="true"></span>
            </router-link>
          </td>
          <td>{{ r.path }}</td>
          <td>{{ r.tf_version }}</td>
          <td>{{ r.serial }}</td>
          <td>{{ r.module_path }}</td>
          <td>
            {{ r.resource_type }}.{{ r.resource_name }}{{ r.resource_index }}
          </td>
          <td>{{ r.attribute_key }}</td>
          <td class="attr-val">{{ r.attribute_value }}</td>
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
  title: "Search",
  components: {
    Multiselect,
  },
  data() {
    return {
      data: {
        tf_versions: [],
        resTypes: [],
        resIDs: [],
        attrKeys: [],
        attrVals: [],
      },
      search: {
        tf_version: null,
        resType: null,
        resID: null,
        attrKey: null,
        attrVal: null,
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
    clearTfVersion() {
      this.search.tf_version = null;
      this.doSearch();
    },
    clearResourceType() {
      this.search.resType = null;
      this.doSearch();
    },
    clearResourceID() {
      this.search.resID = null;
      this.doSearch();
    },
    clearAttributeKey() {
      this.search.attrKey = null;
      this.doSearch();
    },
    resetSearch() {
      this.search.tf_version = null;
      this.search.resType = null;
      this.search.resID = null;
      this.search.attrKey = null;
      this.search.attrVal = null;
      this.doSearch();
    },
    fetchTfVersions() {
      const url = `http://localhost:8080/api/tf_versions`;
      axios
        .get(url)
        .then((response) => {
          this.data.tf_versions = response.data;
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
    fetchResourceTypes() {
      const url = `http://localhost:8080/api/resource/types`;
      axios
        .get(url)
        .then((response) => {
          this.data.resTypes = response.data;
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
    fetchResourceIDs() {
      const url = `http://localhost:8080/api/resource/names`;
      axios
        .get(url)
        .then((response) => {
          this.data.resIDs = response.data;
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
    fetchAttributeKeys() {
      const url = `http://localhost:8080/api/attribute/keys`;
      axios
        .get(url)
        .then((response) => {
          this.data.attrKeys = response.data;
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
    doSearch(page?: number) {
      let params: any = {};
      if (this.search.tf_version != null) {
        params.tf_version = this.search.tf_version;
      }
      if (this.search.resType != null) {
        params.type = this.search.resType;
      }
      if (this.search.resID != null) {
        params.name = this.search.resID;
      }
      if (this.search.attrKey != null) {
        params.key = this.search.attrKey;
      }
      if (this.search.attrVal != null && this.search.attrVal != "") {
        params.value = this.search.attrVal;
      }
      if (page != undefined) {
        params.page = page;
      } else {
        params.page = 1;
      }

      let query = Object.keys(params)
        .map(function(k) {
          return encodeURIComponent(k) + "=" + encodeURIComponent(params[k]);
        })
        .join("&");

      router.push({ name: "Search", query: params });
      const url = `http://localhost:8080/api/search/attribute?` + query;
      axios
        .get(url)
        .then((response) => {
          this.results = response.data;
          this.pager.pages = Math.ceil(
            this.results.total / this.pager.itemsPerPage
          );
          this.pager.page = this.results.page;
          this.pager.prevPage =
            page != undefined && page <= 1 ? undefined : this.pager.page - 1;
          this.pager.nextPage =
            page != undefined && page >= this.pager.pages
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
    this.fetchTfVersions();
    this.fetchResourceTypes();
    this.fetchResourceIDs();
    this.fetchAttributeKeys();

    if(router.currentRoute.value.query.tf_version != null) {
      this.search.tf_version = router.currentRoute.value.query.tf_version;
    }
    if(router.currentRoute.value.query.type != null) {
      this.search.resType = router.currentRoute.value.query.type;
    }

    this.doSearch();
  },
})
export default class Search extends Vue {}
</script>
