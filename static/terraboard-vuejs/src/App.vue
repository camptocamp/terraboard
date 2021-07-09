<template>
  <Navbar/>
  <!-- <div id="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div> -->
  <router-view/>
  <AppFooter/>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Navbar from '@/components/Navbar.vue'; // @ is an alias to /src
import AppFooter from '@/components/Footer.vue';

@Options({
  components: {
    Navbar,
    AppFooter,
  },
})
export default class App extends Vue {}
</script>

<style lang="scss">
a {
  text-decoration: none !important;
}

@media (max-width:767px) {
  .navbar-toggle {
    float: left;
  }
  .navbar-right {
    position: absolute;
    top: 0;
    right: 20px;
  }
  .breadcrumb {
    display: none;
  }
  h2.node-title {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

@media (min-width: 992px) {
  body, html, #mainrow, #leftcol {
    margin: 0;
    height: 100%;
  }
  .container {
    max-width: 75vw;
  }
  #maincont, #leftcol .panel-group {
    overflow: hidden;
  }
  #maincont {
    margin-bottom: 50px;
  }
  #nodes, #node {
    overflow-y: auto;
  }
  #nodes {
    max-height: calc(100% - 200px);
  }
  #node {
    height: 100%;
  }
  #states-select {
    margin: 10px;
  }
  #states-select a {
    color: black;
  }
}
</style>
