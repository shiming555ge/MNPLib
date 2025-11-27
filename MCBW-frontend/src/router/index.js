import { createRouter, createWebHistory } from "vue-router";
import Home from '../components/Home.vue'
import Browse from "../components/Browse.vue";
import Query from "../components/Query.vue";

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/browse',
        name: 'browse',
        component: Browse
    }
    ,
    {
        path: '/query',
        name: 'query',
        component: Query
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router;