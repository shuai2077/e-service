import Vue from 'vue'
import VueRouter from 'vue-router'
// import Record from '../views/Record.vue'
// import User from '../views/User.vue'
import Main from '../views/Main.vue'
import Login from '../views/Login.vue'
import Cookie from 'js-cookie'
Vue.use(VueRouter)
//创建路由组件
//将路由与组件进行映射
//创建router实例

const routes = [
    //主路由
    {
        path: '/',
        component: Main,
        name: 'Main',
        redirect: '/record',
        children:[
            //子路由
            // { path: '/record', name: 'record', component: Record, meta: {requireAuth: true, roles: ['admin']} },  //服务记录管理
            // { path: '/user', name: 'user', component: User, meta: {requireAuth: true, roles: ['admin','user']} },  //服务对象管理
           // { path: '/mall', component: Mall },  //商品管理
        ]
    },
    {
        path: '/login',
        name: 'login',
        component: Login
    }
]

const router = new VueRouter({
    routes // (缩写) 相当于 routes: routes
})

//路由守卫:全局前置导航守卫
router.beforeEach((to, from, next) => {
    // 获取token
    const token = Cookie.get('token')
    if (!token && to.name !== 'login') {
        next({ name: 'login' })
    } else if (token && to.name === 'login') {
        next({ name: 'record' })
    } else {
        next()
    }
})

export default router