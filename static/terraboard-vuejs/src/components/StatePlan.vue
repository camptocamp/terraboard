<template>
  <!-- Plan details view -->
  <div class="mt-3">
    <h3 class="node-title">Plan: {{ this.formatDate(plan.CreatedAt) }}</h3>
    <div class="panel-group">
      <ul class="nav nav-tabs" id="myTab" role="tablist">
        <li class="nav-item" role="presentation">
          <button
            class="nav-link active"
            id="generalInfosPlan-tab"
            data-bs-toggle="tab"
            data-bs-target="#generalInfosPlan"
            type="button"
            role="tab"
            aria-controls="generalInfosPlan"
            aria-selected="true"
          >
            General Informations
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button
            class="nav-link"
            id="parsedPlan-tab"
            data-bs-toggle="tab"
            data-bs-target="#parsedPlan"
            type="button"
            role="tab"
            aria-controls="parsedPlan"
            aria-selected="false"
          >
            Parsed Plan
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button
            class="nav-link"
            id="rawJsonPlan-tab"
            data-bs-toggle="tab"
            data-bs-target="#rawJsonPlan"
            type="button"
            role="tab"
            aria-controls="rawJsonPlan"
            aria-selected="false"
          >
            Raw JSON
          </button>
        </li>
      </ul>
      <div class="tab-content" id="myTabContent">
        <div
          class="tab-pane fade show active p-3"
          id="generalInfosPlan"
          role="tabpanel"
          aria-labelledby="generalInfosPlan-tab"
        >
          <table class="table">
            <tbody>
              <tr>
                <td>Lineage:</td>
                <td>{{ plan.lineage_data.lineage }}</td>
              </tr>
              <tr>
                <td>TF Version:</td>
                <td>{{ plan.terraform_version }}</td>
              </tr>
              <tr>
                <td>Git Commit:</td>
                <td>{{ plan.git_commit }}</td>
              </tr>
              <tr>
                <td>Git Remote:</td>
                <td>{{ plan.git_remote }}</td>
              </tr>
              <tr>
                <td>CI URL:</td>
                <td>{{ plan.ci_url }}</td>
              </tr>
              <tr>
                <td>Source:</td>
                <td>{{ plan.source }}</td>
              </tr>
              <tr>
                <td>Created at:</td>
                <td>{{ this.formatDate(plan.CreatedAt) }}</td>
              </tr>
              <tr>
                <td>Changes:</td>
                <td>
                  <ul>
                    <li>Resources:</li>
                    <ul>
                      <li>Added: {{ this.changes.resources.added }}</li>
                      <li>Changed: {{ this.changes.resources.changed }}</li>
                      <li>Deleted: {{ this.changes.resources.deleted }}</li>
                    </ul>
                    <li>Outputs:</li>
                    <ul>
                      <li>Added: {{ this.changes.outputs.added }}</li>
                      <li>Changed: {{ this.changes.outputs.changed }}</li>
                      <li>Deleted: {{ this.changes.outputs.deleted }}</li>
                    </ul>
                  </ul>
                </td>
              </tr>
              <tr>
                <td class="align-middle">Status:</td>
                <td>
                  <i v-if="isStatusValid" class="fas fa-check-circle fa-2x text-success"></i>
                  <i v-if="!isStatusValid" class="fas fa-exclamation-circle fa-2x text-danger"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div
          class="tab-pane fade p-3"
          id="parsedPlan"
          role="tabpanel"
          aria-labelledby="parsedPlan-tab"
        >
          <table class="table">
            <tbody>
              <tr>
                <td>Format version:</td>
                <td>{{ plan.parsed_plan.format_version }}</td>
              </tr>
              <tr>
                <td>Terraform version:</td>
                <td>{{ plan.parsed_plan.terraform_version }}</td>
              </tr>
              <tr>
                <td>Output changes:</td>
                <td>
                  <ul>
                    <li class="my-2" v-for="output in plan.parsed_plan.output_changes" :key="output">
                      Name: {{ output.name }} <br>
                      Changes:
                      <a class="link-primary" type="button" data-bs-toggle="collapse" :data-bs-target="'#collapse-'+output.name" aria-expanded="false" :aria-controls="'collapse-'+output.name">
                        Show changes
                      </a>
                      <ul class="collapse" :id="'collapse-'+output.name">
                        <li>Actions: {{ output.change.actions }}</li>
                        <li>After: {{ output.change.after }}</li>
                        <li>After sensitive: {{ output.change.after_sensitive }}</li>
                        <li>After unknown: {{ output.change.after_unknown }}</li>
                        <li>Before: {{ output.change.before }}</li>
                        <li>Before sensitive: {{ output.change.before_sensitive }}</li>
                      </ul>
                    </li>
                  </ul>
                </td>
              </tr>
              <tr>
                <td>Resource changes:</td>
                <td>
                  <ul>
                    <li class="my-2" v-for="resource in plan.parsed_plan.resource_changes" :key="resource">
                      Name: {{ resource.name }} <br>
                      Address: {{ resource.address }} <br>
                      Type: {{ resource.type }} <br>
                      Provider: {{ resource.provider_name }} <br>
                      Mode: {{ resource.mode }} <br>
                      Changes:
                      <a class="link-primary" type="button" data-bs-toggle="collapse" :data-bs-target="'#collapse-'+resource.name" aria-expanded="false" :aria-controls="'collapse-'+resource.name">
                        Show changes
                      </a>
                      <ul class="collapse" :id="'collapse-'+resource.name">
                        <li>Actions: {{ resource.change.actions }}</li>
                        <li>After: 
                          <ul>
                            <li v-for="(value, attr) in JSON.parse(resource.change.after)" :key="attr">{{attr}}: {{value}}</li>
                          </ul>
                        </li>
                        <li>After sensitive:
                          <ul>
                            <li v-for="(value, attr) in JSON.parse(resource.change.after_sensitive)" :key="attr">{{attr}}: {{value}}</li>
                          </ul>
                        </li>
                        <li>After unknown:
                          <ul>
                            <li v-for="(value, attr) in JSON.parse(resource.change.after_unknown)" :key="attr">{{attr}}: {{value}}</li>
                          </ul>
                        </li>
                        <li>Before:
                          <ul>
                            <li v-for="(value, attr) in JSON.parse(resource.change.before)" :key="attr">{{attr}}: {{value}}</li>
                          </ul>
                        </li>
                        <li>Before sensitive:
                          <ul>
                            <li v-for="(value, attr) in JSON.parse(resource.change.before_sensitive)" :key="attr">{{attr}}: {{value}}</li>
                          </ul>
                        </li>
                      </ul>
                    </li>
                  </ul>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div
          class="tab-pane fade p-2"
          id="rawJsonPlan"
          role="tabpanel"
          aria-labelledby="rawJsonPlan-tab"
        >
          <pre><code class="language-json">{{JSON.stringify(this.plan.plan_json, null, 2)}}</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  props: {
    plan: {},
  },
  data() {
    return {
      changes: {
        resources: {
          added: 0,
          changed: 0,
          deleted: 0,
        },
        outputs: {
          added: 0,
          changed: 0,
          deleted: 0,
        },
      },
    };
  },
  methods: {
    formatDate(date: string): string {
      return new Date(date).toLocaleString();
    },
    checkPlannedChanges() {
      this.plan.parsed_plan.output_changes.forEach((change: any) => {
        let actions = change.change.actions;
        if (actions.includes("create")) {
          this.changes.outputs.added++;
        } else if (actions.includes("update")) {
          this.changes.outputs.changed++;
        } else if (actions.includes("delete")) {
          this.changes.outputs.deleted++;
        }
      });
      this.plan.parsed_plan.resource_changes.forEach((change: any) => {
        let actions = change.change.actions;
        if (actions.includes("create")) {
          this.changes.resources.added++;
        } else if (actions.includes("update")) {
          this.changes.resources.changed++;
        } else if (actions.includes("delete")) {
          this.changes.resources.deleted++;
        }
      });
    },
  },
  computed: {
    isStatusValid: function(): boolean {
      return (
        this.changes.resources.added +
        this.changes.resources.changed +
        this.changes.resources.deleted +
        this.changes.outputs.added +
        this.changes.outputs.changed +
        this.changes.outputs.deleted
      ) == 0;
    },
    sortedAttributes() {
      if (this.resource.attributes !== undefined) {
        return this.resource.attributes.sort((a: any, b: any) => {
          return a.key.localeCompare(b.key);
        });
      }
    },
  },
  created() {
    console.log(this.plan);
    this.checkPlannedChanges();
  },
})
export default class StatePlan extends Vue {}
</script>

<style lang="scss"></style>
