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
            <tr v-for="(r, index) in results.states" :key="r">
                <td><span class="glyphicon glyphicon-link" aria-hidden="true"></span> <a href="lineage/{{r.lineage_value}}?versionid={{r.version_id}}">{{r.path}}</a></td>
                <td>{{r.lineage_value}}</td>
                <td>{{r.terraform_version}}</td>
                <td>{{r.serial}}</td>
                <td>{{formatDate(r.last_modified)}}</td>
                <td>{{r.resource_count}}</td>
                <td class="text-center">
                    <canvas v-bind:id="'spark-'+index" width="200" height="80" style="max-width: 200px; max-height: 80px;">
                      {{getActivity(index, r.lineage_value, 'spark-'+index)}}
                    </canvas>
                </td>
            </tr>
        </tbody>
    </table>
</div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import { Chart, ChartItem, CategoryScale, PointElement,
LineController, LineElement, LinearScale, Tooltip } from 'chart.js'
import axios from "axios"

Chart.register( CategoryScale, LineElement, LineController, LinearScale, PointElement, Tooltip )

@Options({
  data() {
    return {
      versionMap: {},
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
    getActivity(idx: number, lineage: string, elementId: string): void {
      const url = `http://172.18.0.5:8080/api/state/activity/`+lineage;
      axios.get(url)
        .then((response) => {
          let states = response.data;
          this.versionMap[lineage] = {};
          let activityData = [];
          for (let i = 0; i < states.length; i++) {
              var date = this.formatDate(states[i].last_modified);
              activityData.push(date+";"+states[i].resource_count);
              this.versionMap[lineage][date] = states[i].version_id;
          }

          var activity = activityData.join(",");
          this.results.states[idx].activity = activity;

          let labels: string[] = [];
          let data: string[] = [];
          activityData.forEach((value: string, i: number) => {
            let split = value.split(';');
            labels[i] = split[0];
            data[i] = split[1];
          });
          
          this.createSparkChart(elementId, labels, data)
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
    formatDate(date: string): string {
        return new Date(date).toLocaleString();
    },
    createSparkChart(id: string, labels: string[], data: string[]): void {
      const ctx = document.getElementById(id) as ChartItem;
      const sparkchart = new Chart(ctx, {
        type: 'line',
        data: {
          labels: labels,
          datasets: [
            {
              data: data
            }
          ]
        },
        options: {
          responsive: true,
          elements: {
            line: {
              borderColor: '#4dc9f6',
              borderWidth: 1
            },
            point: {
              radius: 1
            }
          },
          scales: {
            yAxes:
              {
                display: true,
                ticks: {
                  stepSize: 1
                }
              },
            xAxes:
              {
                display: false
              }
          },
          plugins: {
            legend: {
              display: false
            },
            tooltip: {
              enabled: true
            },
          }
        }
      });
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
      const url = `http://172.18.0.5:8080/api/states/stats?page=`+page;
      axios.get(url)
        .then((response) => {
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
