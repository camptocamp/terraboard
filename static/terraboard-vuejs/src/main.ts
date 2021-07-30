import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import titleMixin from "./mixins/titleMixin";

import "@vueform/multiselect/themes/default.css";

createApp(App)
  .mixin(titleMixin)
  .use(router)
  .mount("#app");
