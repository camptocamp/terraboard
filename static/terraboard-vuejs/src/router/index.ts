import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'
import State from '../views/State.vue'
import Search from '../views/Search.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/state/:path',
    name: 'State',
    component: State
  },
  {
    path: '/search',
    name: 'Search',
    component: Search
  },
  {
    path: "/:catchAll(.*)",
    redirect: {
      name: 'Home'
    }
  },
  // {
  //   path: '/about',
  //   name: 'About',
  //   component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  // }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
