import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'
import Lineage from '../views/Lineage.vue'
import State from '../views/State.vue'
import Plan from '../views/Plan.vue'
import Search from '../views/Search.vue'
import PlansExplorer from '../views/PlansExplorer.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/lineage/:lineage',
    name: 'Lineage',
    component: Lineage
  },
  {
    path: '/lineage/:lineage/states',
    name: 'States',
    component: State
  },
  {
    path: '/lineage/:lineage/plans',
    name: 'Plans',
    component: Plan
  },
  {
    path: '/search',
    name: 'Search',
    component: Search
  },
  {
    path: '/plans',
    name: 'PlansExplorer',
    component: PlansExplorer
  },
  {
    path: "/:catchAll(.*)",
    redirect: {
      name: 'Home'
    }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
