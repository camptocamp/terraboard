<template>
  <div class="navbar navbar-light bg-light mb-4 navbar-expand-lg">
    <div class="container-fluid mx-1">
      <!-- .btn-navbar is used as the toggle for collapsed navbar content -->
      <a class="navbar-brand mr-2" href="#"><img src="../assets/logo.png"/></a>
      <button
        type="button"
        class="navbar-toggler collapsed"
        data-bs-toggle="collapse"
        data-bs-target="#navbar-collapse-menu"
        aria-controls="navbar-collapse-menu"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class=" navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbar-collapse-menu">
        <ul class="nav navbar-nav flex-grow-1 align-items-center">
          <li class="nav-item mx-2">
            <router-link to="/" class="nav-link"
              ><span class="fas fa-th-list" aria-hidden="true"></span>
              Overview
            </router-link>
          </li>
          <li class="nav-item mx-2">
            <router-link to="/search" class="nav-link"
              ><span class="fas fa-search" aria-hidden="true"></span> Search
            </router-link>
          </li>
        </ul>
        <ul class="nav navbar-nav align-items-center ml-auto">
          <li id="states-select" class="nav-item">
            <Multiselect
              id="states-quick-access"
              v-model="states_select.value"
              v-bind="states_select"
              placeholder= "Enter a state file path..."
              @change="goToState"
              @select="clearSelect"
              ref="quickAccess"
            >
            </Multiselect>
          </li>
          <ul
            class="nav navbar-nav navbar-collapse collapse"
            id="navbar-collapse-menu"
          >
            <li class="dropdown nav-item">
              <a
                data-toggle="dropdown"
                class="dropdown-toggle nav-link"
                id="dropdownMenuButton"
                data-bs-toggle="dropdown"
                aria-expanded="false"
                ><img src="" height="25px"
              /></a>

              <ul
                class="dropdown-menu dropdown-menu-end"
                aria-labelledby="dropdownMenuButton1"
              >
                <li><a class="dropdown-item" href="#">Logged in as </a></li>
                <li>
                  <a class="dropdown-item" href="/oauth2/sign_in">Sign out</a>
                </li>
              </ul>
            </li>
          </ul>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Multiselect from '@vueform/multiselect'
import axios from "axios"
import router from "../router";

@Options({
  data() {
    return {
      states_select: {
        options: [],
        value: null,
        searchable: true,
      }
    };
  },
  methods: {
    reset() {
      this.item = {};
    },
    goToState(value: any) {
      if (value != null) {
        router.push({name: "State", params: {lineage: value.lineage_value}, query: { versionid: value.version_id } });
      }
    },
    clearSelect() {
      this.$refs.quickAccess.clear()
    },
    fetchStates() {
      const url = `/api/lineages/stats`
      axios.get(url)
        .then((response) => {
          // handle success
          response.data.states.forEach((obj: any) => {
            let entry = {value: obj, label: obj.path}
            this.states_select.options.push(entry)
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
    }
  },
  mounted() {
    this.fetchStates();
  },
  components: {
    Multiselect,
  },
})
export default class Navbar extends Vue {}
</script>

<style scoped lang="scss">
#states-select {
  @media (min-width: 992px) {
    width: 30vw;
  }
  @media (max-width: 992px) {
    width: 40%;
  }
}
</style>
