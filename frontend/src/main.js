import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import { AddToLibrary } from '../bindings/derpy-launcher072/library/library'

//AddToLibrary(69)

const app = createApp(App)
app.use(router)
app.mount('#app')