<template>
  <div id="app">
    <el-row>
      <!-- 路由出口 -->
      <!-- 路由匹配到的组件将渲染在这里 -->
      <router-view></router-view>
</el-row>
  </div>
</template>
<script>
import Cookie from 'js-cookie'
export default({
  name: 'App',
  created() {
  window.addEventListener('beforeunload', this.handlePageClose);
},
  methods: {
    handlePageClose() {
      //清楚cookie中的token
      Cookie.remove('token')
      //清除cookie中的menu
      Cookie.remove('menu')
      //跳转到登录界面
      //this.$router.push('/login')
    }
  }
})
</script>
<style lang="less">
  html,body,h3,p{
    margin: 0;
    padding: 0;
  }
</style>
