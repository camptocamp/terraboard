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

@Options({
  data() {
    return {
      version: String,
      copyright: String
    };
  },
  methods: {
    getVersion(): string {
      return "test1"
    },
    getCopyright(): string {
      return "test2"
    }
  },
  created() {
    this.version = this.getVersion();
    this.copyright = this.getCopyright();
  },
})
export default class Footer extends Vue {
  version!: string
  copyright!: string
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

</style>
