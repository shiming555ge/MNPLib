import './scss/styles.scss'
import * as bootstrap from 'bootstrap'
import "/node_modules/bootstrap-icons/font/bootstrap-icons.min.css"

import { createI18n } from 'vue-i18n';
import zh_cn from './locales/zh_cn.json';
import en_us from './locales/en_us.json'

const i18n = createI18n({
    locale: "en_us",
    fallbackLocale: "zh_cn",
    messages: {zh_cn, en_us}
})

import router from './router';

import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)
app.use(i18n)
app.use(router)
app.mount("#app")


