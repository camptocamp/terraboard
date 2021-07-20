<template>
  <!-- Compare view -->
  <div
    class="row mt-3"
    v-if="
      compareDiff.differences !== undefined || compare.differences !== undefined
    "
  >
    <div class="panel-group">
      <div class="card mt-3">
        <h4 id="card-title-diff" class="card-header bg-warning">
          <div data-toggle="collapse" data-target="#card-body-diff">
            Differences
            <span
              id="badge-diff"
              class="badge rounded-pill bg-secondary float-end"
              >{{ compareDiff.differences }}</span
            >
          </div>
        </h4>
        <div id="card-body-diff" class="panel-collapse collapse show">
          <div class="card-body">
            <div
              class="list-group resource"
              v-for="(diff, resource) in compare.differences.resource_diff"
              v-bind:key="resource"
            >
              <div class="resource-title">{{ resource }}</div>
              <pre><code class="language-diff">{{diff.unified_diff}}</code></pre>              
            </div>
          </div>
        </div>
      </div>
      <div class="card mt-4">
        <h4 id="card-title-in-old" class="card-header bg-danger">
          <div data-toggle="collapse" data-target="#card-body-in-old">
            Only in serial {{ compare.stats.from.serial }}
            <span
              id="badge-diff"
              class="badge rounded-pill bg-secondary float-end"
              >{{ compareDiff.only_in_old }}</span
            >
          </div>
        </h4>
        <div id="card-body-in-old" class="panel-collapse collapse show">
          <div class="card-body">
            <div
              class="list-group resource"
              v-for="(code, resource) in compare.differences.only_in_old"
              v-bind:key="resource"
            >
              <div class="resource-title">{{ resource }}</div>
              <pre><code class="language-ruby">{{code}}</code></pre>
            </div>
          </div>
        </div>
      </div>
      <div class="card mt-4">
        <h4 id="card-title-in-new" class="card-header bg-success">
          <div data-toggle="collapse" data-target="#card-body-in-new">
            Only in serial {{ compare.stats.to.serial }}
            <span
              id="badge-diff"
              class="badge rounded-pill bg-secondary float-end"
              >{{ compareDiff.only_in_new }}</span
            >
          </div>
        </h4>
        <div id="card-body-in-new" class="panel-collapse collapse show">
          <div class="card-body">
            <div
              class="list-group resource"
              v-for="(code, resource) in compare.differences.only_in_new"
              v-bind:key="resource"
            >
              <div class="resource-title">{{ resource }}</div>
              <pre><code class="language-ruby">{{code}}</code></pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  props: {
    compare: {},
    compareDiff: {},
  },
})
export default class StatesCompare extends Vue {}
</script>

<style lang="scss">
</style>
