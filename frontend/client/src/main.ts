import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

import "./assets/tailwind.css";
import "./assets/main.css";
import "./assets/frontend-theme.css";
import "./styles/vant-pull-refresh.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount("#app");
