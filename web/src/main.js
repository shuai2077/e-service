import Vue from 'vue'
import App from './App.vue'
import 'element-ui/lib/theme-chalk/index.css'
//import {Row, Button} from 'element-ui'
import ElementUI from 'element-ui'
import router from './router'
import store from './store'
 
Vue.config.productionTip = false

//按需引入
// Vue.use(Row)
// Vue.use(Button)

//全局引入
Vue.use(ElementUI)

new Vue({
  router,
  store,
  render: h => h(App),
  created() {
    store.commit('addMenu', router)
  }
}).$mount('#app')