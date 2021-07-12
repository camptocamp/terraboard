<template>
<div id="results" class="row">
    <label id="navigate"> <span class="fas fa-caret-left"
            v-on:click="fetchStats(prevPage)"></span> {{startItems}}-{{itemsInPage}}/{{results.total}}
        <span class="fas fa-caret-right" v-on:click="fetchStats(nextPage)"></span>
    </label>
    <table class="table table-border table-striped">
        <thead>
            <th>
                Path
            </th>
            <th>
                Lineage
            </th>
            <th>
                TF Version
            </th>
            <th>
                Serial
            </th>
            <th>
                Time
            </th>
            <th>
                Resources
            </th>
            <th>
                Activity
            </th>
        </thead>
        <tbody>
            <tr v-for="r in results.states" :key="r">
                <td><span class="glyphicon glyphicon-link" aria-hidden="true"></span> <a href="lineage/{{r.lineage_value}}?versionid={{r.version_id}}">{{r.path}}</a></td>
                <td>{{r.lineage_value}}</td>
                <td>{{r.terraform_version}}</td>
                <td>{{r.serial}}</td>
                <td>{{formatDate(r.last_modified)}}</td>
                <td>{{r.resource_count}}</td>
                <!-- <td>
                    <sparklinechart lineage="{{r.lineage_value}}" data="{{r.activity}}"></sparklinechart>
                </td> -->
            </tr>
        </tbody>
    </table>
</div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import axios from "axios"

@Options({
  data() {
    return {
      results: {},
      pages: 0,
      page: 0,
      prevPage: 0,
      nextPage: 0,
      startItems: 0,
      itemsInPage: 0,
      itemsPerPage: 20,
    }
  },
  methods: {
    formatDate(date: string): string {
        return new Date(date).toLocaleString();
    },
    updatePager(response: any): void {
      this.results = response.data;
      this.pages = Math.ceil(this.results.total / this.itemsPerPage);
      this.page = this.results.page;
      this.prevPage = (this.page <= 1) ? undefined : this.page - 1;
      this.nextPage = (this.page >= this.pages) ? undefined : this.page + 1;
      this.startItems = this.itemsPerPage * (this.page - 1) + 1;
      this.itemsInPage = Math.min(this.itemsPerPage * this.page, this.results.total);
    },
    fetchStats(page: number): void {
      const url = `http://172.22.0.5:8080/api/states/stats?page=`+page;
      axios.get(url)
        .then((response) => {
          console.log(response);
          this.updatePager(response);
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
    }
  },
  created() {
    this.fetchStats(1);
  },
})
export default class StatesList extends Vue {}
</script>

<style scoped lang="scss">

</style>
