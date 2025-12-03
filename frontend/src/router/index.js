import { createRouter, createWebHistory } from "vue-router";
import Home from '../components/Home.vue'
import Browse from "../components/Browse.vue";
import Query from "../components/Query.vue";
import About from '../components/About.vue';
import SuperAdmin from '../components/SuperAdmin.vue';
import NotFound from '../components/NotFound.vue';

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
    },
    {
        path: '/query',
        name: 'query',
        component: Query
    },
    {
        path: '/about',
        name: 'about',
        component: About
    },
    {
        path: '/superadmin',
        name: 'superadmin',
        component: SuperAdmin
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'not-found',
        component: NotFound
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

// 路由守卫示例：全局前置守卫
router.beforeEach((to, from, next) => {
    console.log(`路由跳转: ${from.path} -> ${to.path}`);
    
    next();
})

// 全局后置守卫
router.afterEach((to, from) => {
    window.scrollTo(0, 0);
    
    // 示例：更新页面标题
    const pageTitles = {
        '/': 'MNPLib - Home',
        '/browse': 'MNPLib - Browse Compounds',
        '/query': 'MNPLib - Chemical Query',
        '/about': 'MNPLib - About',
        '/superadmin': 'MNPLib - Super Admin'
    };
    
    const title = pageTitles[to.path] || 'MNPLib - Marine Natural Product Library';
    document.title = title;
})

export default router;
