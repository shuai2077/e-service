<template>
    <div class="manage">
        <div class="manage-header" style="display: flex; justify-content: flex-end;"> 
            <!-- form搜索区域 -->
            <el-form :inline="true" :model="userForm">
                <el-form-item>
                    <el-input placeholder="请输入姓名或地址" v-model="userForm.name"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="onSubmit">查询</el-button>
                </el-form-item>
            </el-form>
        </div>
       <div class="common-table">
            <el-table
                stripe
                height="90%"
                :data="tableData"
                style="width: 100%"
                :cell-style="{ textAlign: 'center' }"
                :header-cell-style="{ textAlign: 'center' }">
                <el-table-column
                    prop="name"
                    label="姓名">
                </el-table-column>
                <el-table-column
                    prop="sex"
                    label="性别">
                    <template slot-scope="scope">
                        <span>{{ scope.row.sex == "0" ? '男': '女' }}</span>
                    </template>
                </el-table-column>
                <el-table-column
                    prop="age"
                    label="年龄">
                </el-table-column>
                <!-- <el-table-column
                    prop="birth"
                    label="出生日期">
                </el-table-column> -->
                <el-table-column
                    prop="phone"
                    label="联系电话">
                </el-table-column>
                <el-table-column
                    prop="addr"
                    label="地址">
                </el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button type="danger" size="mini" @click="handleDelete(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pager">
                <el-pagination
                    @current-change="handleChange"
                    :current-page="page"
                    :page-limit="limit"
                    :layout="prev, pager, next"
                    :total="total">
                </el-pagination>
            </div>
       </div>
    </div>
</template>
<script>
import {getUser, addUser, editUser, delUser} from '../api'
export default({
    data(){
        return{
            dialogVisible: false,
            form:{
                name: '',
                age: '',
                sex: '',
                phone: '',
                addr: '',
            },
            rules: {
                name: [
                    { required: true, message: '请输入名称', trigger: 'blur' }
                ],
                age: [
                    { required: true, message: '请输入年龄', trigger: 'blur' }
                ],
                sex: [
                    { required: true, message: '请选择性别', trigger: 'blur' }
                ],
                phone: [
                    { required: true, message: '请输入联系电话', trigger: 'blur' }
                ],
                // birth: [
                //     { required: true, message: '请选择出生日期', trigger: 'blur' }
                // ],
                addr: [
                    { required: true, message: '请输入地址', trigger: 'blur' }
                ],
            },
            tableData:[],
            modelType: 0, //0表示新增弹窗，1表示编辑
            total: 0,  //当前的总条数，默认为0
            pageData:{
                page: 1,
                limit: 10
            },
            userForm:{
                name: ''
            }
        }
    },
    methods:{
        handleClose(){
            this.$refs.form.resetFields()
            this.dialogVisible = false
        },
        cancel(){
            this.handleClose()
        },
        handleEdit(row){
            this.modelType=1
            this.dialogVisible=true
            //注意需要对当前行数据进行深拷贝，否则会出错
            this.form=JSON.parse(JSON.stringify(row))
        },
        handleDelete(row){
            this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                delUser({name: row.name}).then(() => {
                    this.$message({
                        type: 'success',
                        message: '删除成功!'
                    });
                    //重新获取列表的接口
                    this.getList()
                })
                
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });          
            });
        },
        handleAdd(){
            this.modelType=0
            this.dialogVisible=true
        },
        //获取列表数据
        getList(){
            getUser({params: {...this.userForm, ...this.pageData}}).then(({data}) => {
                  //console.log(data)
                  if (data.list == null){
                    this.tableData=data.list
                    this.total = data.count || 0
                  }else {
                    this.tableData = data.list.slice(
                        (this.pageData.page - 1) * this.pageData.limit,
                        this.pageData.page * this.pageData.limit
                    );
                    this.total = data.count || 0
                  }
            })
        },
        //选择页码的回调函数
        handleChange(val){
            this.pageData.page=val
            this.getList()
        },
        //列表查询
        onSubmit(){
            this.getList()
        }
    },
    mounted(){
        //页面加载时首先执行该函数
        this.getList()
    }
})
</script>
<style lang="less" scoped>
.manage{
    height: 90%;
    .manage-header{
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    .common-table{
        position:relative;
        height: calc(100% - 62px);
        .pager{
            position: absolute;
            bottom: 0;
            right: 20px;
        }
    }
}
</style>