<template>
<div class="navbar mt-auto">
    <div class="container-fluid mx-1">
        <ul class="nav navbar-nav" id="navbar-collapse-menu">
            <li><a href="https://github.com/camptocamp/terraboard/releases" target="_blank">Terraboard v{{ version }}</a></li>
        </ul>
        <ul class="nav navbar-nav navbar-right" id="navbar-collapse-menu">
            <li><a href="https://www.camptocamp.com/" target="_blank">{{ copyright }}</a></li>
        </ul>
    </div>
</div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import axios from "axios";

@Options({
  data() {
    return {
      version: "",
      copyright: "",
    };
  },
  methods: {
    fetchVersion(): void {
      const url = `/api/version`;
      axios.get(url)
        .then((response) => {
          this.version = response.data.version;
          this.copyright = response.data.copyright;
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
  created() {
    this.fetchVersion();
  },
})
export default class Footer extends Vue {
  version!: string
  copyright!: string
}
</script>

<style scoped lang="scss">

</style>
