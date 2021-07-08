/* Mixin used to dynamically manage page title for each vue views */

function getTitle(vm: any) {
  const title = vm.$options;
  if (title) {
    return typeof title === "function" ? title.call(vm) : title;
  }
}

export default {
  mounted() {
    const title = getTitle(this);
    if (title) {
      document.title = "Terraboard - " + title.title;
    }
  },
};
