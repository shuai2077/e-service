<template>
    <el-form ref="form" label-width="70px" :inline="true" class="login-container" :model="form" :rules="rules">
        <h3 class="login_title">系统登录</h3>
        <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username" placeholder="请输入账号"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
            <el-input type="password" v-model="form.password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item>
            <el-button @click="submit" style="margin-left:110px;margin-top:10px" type="primary">登录</el-button>
        </el-form-item>
    </el-form>
</template>
<script>
import Cookie from 'js-cookie'
import { getMenu } from '../api'
export default {
    data() {
        return {
            form: {
            username: '',
            password: ''
            },
            rules: {
            username: [
                { required: true, message: '请输入用户名', trigger: 'blur' }
            ],
            password: [
                { required: true, message: '请输入密码', trigger: 'blur' }
            ]
            },
            adminmenu: [
            {
                path: "/record",
                name: "record",
                label: "服务记录管理",
                icon: "document",
                url: "Record.vue",
            },
            {
                path: "/user",
                name: "user",
                label: "服务对象管理",
                icon: "user",
                url: "User.vue"
            },
            {
                path: "/log",
                name: "log",
                label: "修改记录管理",
                icon: "edit",
                url: "Log.vue"
            }
            ],
            usermenu: [
            {
                path: "/user",
                name: "user",
                label: "服务对象管理",
                icon: "user",
                url: "User.vue"
            }
            ],
        }
    },
    methods:{
        //登录
        submit(){
            this.$refs.form.validate((valid) => {
                if (valid){
                    getMenu(this.form).then(({data}) => {
                        //console.log(data)
                        if(data.data.code === 20000){
                            this.$message({
                                message: '登录成功',
                                type: 'success'
                            });
                            //token信息存入cookie用于不同页面间的通信
                            Cookie.set('token',data.data.token)
                            //跳转到首页
                            if(data.data.flag === 1){
                                //获取菜单的数据，存入store中
                                this.$store.commit('setMenu',this.adminmenu)
                                //动态添加路由
                                this.$store.commit('addMenu',this.$router)
                                this.$router.push('/record')
                            }else{
                                //获取菜单的数据，存入store中
                                this.$store.commit('setMenu',this.usermenu)
                                //动态添加路由
                                this.$store.commit('addMenu',this.$router)
                                this.$router.push('/user')
                            }
                        }else{
                            this.$message.error('不存在该用户')
                        }
                    })
                }
            })
        }
    }
}
</script>
<style lang="less" scoped>
.login-container{
    width: 350px;
    border: 1px solid #eaeaea;
    // 居中
    margin: 180px auto;
    padding: 35px 35px 15px 35px;
    border-radius: 15px;
    background-color: #fff;
    box-shadow: 0 0 25px #cac6c6;
    box-sizing: border-box;
    .el-input {
        width: 198px;
    }
    .login_title {
        color: #505458;
        // 左右居中
        text-align: center;
        margin-bottom: 40px;
    }
}
</style>