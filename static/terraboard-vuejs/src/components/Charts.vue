<template>
<div class="row justify-content-around">
    <div class="overview-chart col-6 col-md-3 col-xxl-4 text-center" style="min-width: 100px; max-width: 300px;">
        <canvas id="chart-pie-resource-types" class="chart mb-2" chart-click="searchType"></canvas>
        <h5>Resource types</h5>
    </div>
    <div class="overview-chart col-6 col-md-3 col-xxl-4 text-center" style="min-width: 100px; max-width: 300px;">
        <canvas id="chart-pie-terraform-versions" class="chart mb-2"
            chart-click="searchTfVersion"></canvas>
        <h5>Terraform versions</h5>
    </div>
    <div class="overview-chart col-6 col-md-3 col-xxl-4 text-center" style="min-width: 100px; max-width: 300px;">
        <canvas id="chart-pie-ls" class="chart mb-2"></canvas>
        <h5>States locked</h5>
    </div>
</div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import { Chart, ChartItem, PieController, ArcElement, Tooltip } from 'chart.js'
import axios from "axios"

Chart.register( PieController, ArcElement, Tooltip )

const chartOptions = 
{ 
  responsive: true,
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      display: true,
    },
  } 
}

@Options({
  data() {
    return {
      locks: {},
      statesTotal: 0,
      pieResourceTypes: {
        labels: [[], [], [], [], [], [], ["Total"]],
        data: [0, 0, 0, 0, 0, 0, 0],
        options: chartOptions,
      },
      pieTfVersions: {
        labels: [[], [], [], [], [], [], ["Total"]],
        data: [0, 0, 0, 0, 0, 0, 0],
        options: chartOptions,
      },
      pieLockedStates: {
        labels: ["Locked", "Unlocked"],
        data: [0, 0],
        options: chartOptions,
      },
    };
  },
  methods: {
    isLocked(path: string): boolean {
      if (path in this.locks) {
          return true;
      }
      return false;
    },
    fetchResourceTypes(): void {
      const url = `http://172.18.0.5:8080/api/resource/types/count`;
      axios.get(url)
        .then((response) => {
          response.data.forEach((value: any, i: number) => {
            if(i < 6) {
                this.pieResourceTypes.labels[i] = value.name;
                this.pieResourceTypes.data[i]   = parseInt(value.count, 10);
            } else {
                this.pieResourceTypes.labels[6].push(value.name+": "+value.count);
                this.pieResourceTypes.data[6] += parseInt(value.count, 10);
            }
          });

          const ctx = document.getElementById('chart-pie-resource-types') as ChartItem;
          const resourcesChart = new Chart(ctx, {
              type: 'pie',
              data: {
                  labels: this.pieResourceTypes.labels,
                  datasets: [{
                      label: 'States Resources Type',
                      data: this.pieResourceTypes.data,
                      backgroundColor: [
                        '#4dc9f6',
                        '#f67019',
                        '#f53794',
                        '#537bc4',
                        '#acc236',
                        '#166a8f',
                        '#00a950',
                      ],
                      hoverOffset: 4
                  }]
              },
              options: this.pieResourceTypes.options
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
    fetchVersions(): void {
      const url = `http://172.18.0.5:8080/api/states/tfversion/count?orderBy=version`;
      axios.get(url)
        .then((response) => {
          response.data.forEach((value: any, i: number) => {
            if(i < 6) {
                this.pieTfVersions.labels[i] = [value.name];
                this.pieTfVersions.data[i]   = parseInt(value.count, 10);
            } else {
                this.pieTfVersions.labels[6].push(value.name+": "+value.count);
                this.pieTfVersions.data[6] += parseInt(value.count, 10);
            }
          });

          const ctx = document.getElementById('chart-pie-terraform-versions') as ChartItem;
          const versionsChart = new Chart(ctx, {
              type: 'pie',
              data: {
                  labels: this.pieTfVersions.labels,
                  datasets: [{
                      label: 'States Versions',
                      data: this.pieTfVersions.data,
                      backgroundColor: [
                        '#4dc9f6',
                        '#f67019',
                        '#f53794',
                        '#537bc4',
                        '#acc236',
                        '#166a8f',
                        '#00a950',
                      ],
                      hoverOffset: 4
                  }]
              },
              options: this.pieTfVersions.options
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
    fetchLocks(): void {
      const url = `http://172.18.0.5:8080/api/locks`;
      axios.get(url)
        .then((response) => {
          this.locks = response.data;

          const ctx = document.getElementById('chart-pie-ls') as ChartItem;
          const locksChart = new Chart(ctx, {
              type: 'pie',
              data: {
                  labels: this.pieLockedStates.labels,
                  datasets: [{
                      label: 'States Locks Status',
                      data: this.pieLockedStates.data,
                      backgroundColor: [
                        '#f67019',
                        '#4dc9f6',
                      ],
                      hoverOffset: 4
                  }]
              },
              options: this.pieLockedStates.options
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
  },
  watch: {
    '$data.locks': {
      handler: function(nv) {
        this.pieLockedStates.data[0] = Object.keys(nv).length;
        this.pieLockedStates.data[1] -= Object.keys(nv).length;
      },
      deep: true
    },
    '$data.statesTotal': {
      handler: function(nv) {
        this.pieLockedStates.data[1] = nv - this.pieLockedStates.data[0];
      },
    }
  },
  created() {
    this.fetchResourceTypes();
    this.fetchVersions();

    const url = `http://172.18.0.5:8080/api/states/stats?page=1`;
      axios.get(url)
        .then((response) => {
          this.statesTotal = response.data.total;
          this.fetchLocks();
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
})
export default class Charts extends Vue {}
</script>

<style scoped lang="scss">

</style>

À titre perso, si on se base sur le 01/01/2000 à 00h j'avais approximativement 15 minutes 
