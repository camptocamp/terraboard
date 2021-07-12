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
  title: 'Home',
  components: {
    Charts,
    StatesList,
  },
})
export default class Home extends Vue {}
</script>
