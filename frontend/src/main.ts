import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import enrutador from './enrutador'
import './assets/estilos.css'

const app = createApp(App)

app.use(createPinia())
app.use(enrutador)
app.mount('#app')
