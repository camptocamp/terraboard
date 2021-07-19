<template>
  <!-- Resource details view -->
  <div v-if="resource.attributes !== undefined" class="mt-3">
    <h3 class="node-title">
      {{ resource.type }}.{{ resource.name }}{{ resource.index }}
    </h3>
    <div class="panel-group">
      <div class="card">
        <h4 class="card-header">Attributes</h4>
        <table class="table">
          <thead>
            <th>Attribute</th>
            <th>Value</th>
          </thead>
          <tbody>
            <tr v-for="attr in sortedAttributes" v-bind:key="attr">
              <td class="attr-key">{{ attr.key }}</td>
              <td class="attr-val">{{ attr.value }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  props: {
    resource: {},
  },
  computed: {
    sortedAttributes() {
      if (this.resource.attributes !== undefined) {
        return this.resource.attributes.sort((a: any, b: any) => {
          return a.key.localeCompare(b.key);
        });
      }
    },
  },
})
export default class StateDetails extends Vue {}
</script>

<style lang="scss"></style>
