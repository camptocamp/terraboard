import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import titleMixin from "./mixins/titleMixin";

createApp(App)
  .mixin(titleMixin)
  .use(router)
  .mount("#app");
