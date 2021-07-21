<template>
  <Charts/>
  <hr>
  <StatesList/>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Charts from '@/components/Charts.vue'; // @ is an alias to /src
import StatesList from '@/components/StatesList.vue'; // @ is an alias to /src

@Options({
  title: 'Overview',
  components: {
    Charts,
    StatesList,
  },
  created() {
    this.updateTitle();
  },
})
export default class Home extends Vue {}
</script>
