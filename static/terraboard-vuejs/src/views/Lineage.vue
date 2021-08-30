<template>
    <h1 v-if="this.lastPath != ''">Workspace: {{this.lastPath}}</h1>
    <label id="navigate" v-if="data.length != 0">
    <span
      class="fas fa-caret-left"
      role="button"
      v-on:click="refreshList(pager.prevPage)"
      v-if="pager.prevPage"
    ></span>
    {{ pager.startItems }}-{{ pager.itemsInPage }}/{{ data.length }}
    <span
      class="fas fa-caret-right"
      role="button"
      v-on:click="refreshList(pager.nextPage)"
      v-if="pager.nextPage"
    ></span>
  </label>
  <div id="results">
    <table class="table table-border table-striped">
      <thead>
        <th>Type</th>
        <th>Created At</th>
        <th>Actions</th>
      </thead>
      <tbody>
        <tr v-for="r in this.pager.items" v-bind:key="r">
          <td>
            <span v-if="r.type === 'state'" class="badge rounded-pill bg-primary">State</span>
            <span v-else-if="r.type === 'plan'" class="badge rounded-pill bg-secondary">Plan</span>
          </td>
          <td>{{ formatDate(r.createdAt) }}</td>
          <td><button type="button" class="btn btn-outline-primary" @click="goToRessource(r.route)">Consult</button></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import router from "../router";

class ObjWrapper {
  name: string
  type: string
  createdAt: Date
  obj: any
  route: any

  constructor(name: string, type: string, createdAt: Date, obj: any, route: any) {
    this.name = name;
    this.type = type;
    this.createdAt = createdAt;
    this.obj = obj;
    this.route = route;
  }
}

@Options({
  title: "Lineage",
  emits: ["refresh"],
  data() {
    return {
      data: [] as ObjWrapper[],
      lastPath: '',
      lineage: "",
      pager: {
        items: [],
        pages: 0,
        page: 0,
        prevPage: 0,
        nextPage: 0,
        startItems: 0,
        itemsInPage: 0,
        itemsPerPage: 20,
      },
    };
  },
  methods: {
    formatDate(date: Date): string {
      return date.toUTCString();
    },
    getLastStatePath(): string {
      for (let i = 0; i < this.data.length; i++) {
        const obj = this.data[i];
        if (obj.type == 'state') {
          return obj.name;
        }
      }
      return "";
    },
    goToRessource(route: any) {
      router.push(route);
    },
    fetchData() {
      const url = `/api/lineages/` + this.lineage + `/activity`
      axios.get(url)
        .then((response) => {
          // handle success
          response.data.forEach((obj: any) => {
            let entry = new ObjWrapper(
              obj.path, 
              "state", 
              new Date(obj.last_modified), 
              obj,
              {name: "States", params: {lineage: this.lineage}, query: { versionid: obj.version_id }},
            );
            this.data.push(entry)
          });

          const url = `/api/plans?lineage=`+this.lineage;
          axios.get(url)
            .then((response) => {
              // handle success
              response.data.plans.forEach((obj: any) => {
                let entry = new ObjWrapper(
                  obj.lineage_data.lineage, 
                  "plan", 
                  new Date(obj.CreatedAt), 
                  obj,
                  {name: "Plans", params: {lineage: this.lineage}, query: { planid: obj.ID }},
                );
                this.data.push(entry)
              });
              this.data.sort((a: any, b: any) => b.createdAt.getTime() - a.createdAt.getTime())
              this.lastPath = this.getLastStatePath();
              this.refreshList(1);
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
      this.pager.pages = Math.ceil(
        this.data.length / this.pager.itemsPerPage
      );
      this.pager.page = page;
      this.pager.items = this.data.slice(this.pager.itemsPerPage * (this.pager.page - 1), this.pager.itemsPerPage * this.pager.page);
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
        this.data.length
      );
    },
  },
  created() {
    this.updateTitle();
  },
  mounted() {
    this.lineage = this.$route.params.lineage;
    this.fetchData();
  },
})
export default class Lineage extends Vue {}
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
