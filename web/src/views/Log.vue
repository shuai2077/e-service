<template>
    <div class="manage">
        <div class="manage-header" style="display: flex; justify-content: flex-end;"> 
            <!-- form搜索区域 -->
            <el-form :inline="true" :model="userForm">
                <el-form-item>
                    <el-input placeholder="请输入姓名" v-model="userForm.name"></el-input>
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
                    prop="tx_id"
                    label="交易ID">
                </el-table-column>
                <el-table-column
                    prop="hisvalue"
                    label="历史值">
                </el-table-column>
                <el-table-column
                    prop="time"
                    label="修改时间">
                </el-table-column>
                <!-- <el-table-column
                    prop="birth"
                    label="出生日期">
                </el-table-column> -->
                <el-table-column
                    prop="is_delete"
                    label="是否被删除">
                    <template slot-scope="scope">
                        <span>{{ scope.row.is_delete == true ? '已删除': '未删除' }}</span>
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
import {getHis, addUser, editUser, delUser} from '../api'
export default({
    data(){
        return{
            dialogVisible: false,
            form:{
                tx_id: '',
                hisvalue: '',
                time: '',
                is_delete: '',
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
        //获取列表数据
        getList(){
            getHis({params: {...this.userForm, ...this.pageData}}).then(({data}) => {
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