<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld msg="Welcome to Your Vue.js + TypeScript App"/>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import HelloWorld from '@/components/HelloWorld.vue';
import router from "../router";

@Options({
  title: 'States',
  components: {
    HelloWorld,
  },
  data() {
    return {
      lineage: "",
      versionid: "",
    }
  },
  created() {
    this.lineage = this.$route.params.lineage;
    this.versionid = router.currentRoute.value.query.versionid;
  },

})
export default class State extends Vue {}
</script>
