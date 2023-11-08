<template>
    <div class="manage">
        <el-dialog
            title="提示"
            :visible.sync="dialogVisible"
            width="30%"
            :before-close="handleClose">
            <!-- 服务记录表单信息 -->
            <el-form ref="form" :rules="rules" :inline="true" :model="form" label-width="80px">
                <el-form-item label="姓名" prop="name">
                    <el-input placeholder="请输入姓名" v-model="form.name"></el-input>
                </el-form-item>
                <el-form-item label="年龄" prop="age">
                    <el-input placeholder="请输入年龄" v-model="form.age"></el-input>
                </el-form-item>
                <el-form-item label="性别" prop="sex">
                    <el-select v-model="form.sex" placeholder="请选择">
                        <el-option label="男" :value="'0'"></el-option>
                        <el-option label="女" :value="'1'"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="联系电话" prop="phone">
                    <el-input placeholder="请输入联系电话" v-model="form.phone"></el-input>
                </el-form-item>
                <el-form-item label="地址" prop="addr">
                    <el-input placeholder="请输入地址" v-model="form.addr"></el-input>
                </el-form-item>
                <el-form-item label="服务内容" prop="content">
                    <el-input placeholder="请输入服务内容" v-model="form.content"></el-input>
                </el-form-item>
                <el-form-item label="服务人员" prop="attendant">
                    <el-input placeholder="请输入服务人员姓名" v-model="form.attendant"></el-input>
                </el-form-item>
                <el-form-item label="完成情况" prop="perform">
                    <el-select v-model="form.perform" placeholder="请选择">
                        <el-option label="未完成" :value="'0'"></el-option>
                        <el-option label="已完成" :value="'1'"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="满意度" prop="satisfaction">
                    <el-select v-model="form.satisfaction" placeholder="请选择">
                        <el-option label="不满意" :value="'0'"></el-option>
                        <el-option label="满意" :value="'1'"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            
            <span slot="footer" class="dialog-footer">
                <el-button @click="cancel">取 消</el-button>
                <el-button type="primary" @click="submit">确 定</el-button>
            </span>
        </el-dialog>
        <div class="manage-header"> 
            <el-button @click="handleAdd" type="primary">
                + 新增
            </el-button>
            <!-- form搜索区域 -->
            <el-form :inline="true" :model="recordForm">
                <el-form-item>
                    <el-input placeholder="请输入姓名或地址" v-model="recordForm.name"></el-input>
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
                <el-table-column
                    prop="phone"
                    label="联系电话">
                </el-table-column>
                <el-table-column
                    prop="addr"
                    label="地址">
                </el-table-column>
                <el-table-column
                    prop="content"
                    label="服务内容">
                </el-table-column>
                <el-table-column
                    prop="attendant"
                    label="服务人员">
                </el-table-column>
                <el-table-column
                    prop="perform"
                    label="完成情况">
                    <template slot-scope="scope">
                        <span>{{ scope.row.perform == "0" ? '未完成': '已完成' }}</span>
                    </template>
                </el-table-column>
                <el-table-column
                    prop="satisfaction"
                    label="满意度">
                    <template slot-scope="scope">
                        <span>{{ scope.row.satisfaction == "0" ? '不满意': '满意' }}</span>
                    </template>
                </el-table-column>
                <el-table-column 
                    fixed="right"
                    label="操作">
                    <template slot-scope="scope">
                        <!--<el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
                        <el-button type="danger" size="mini" @click="handleDelete(scope.row)">删除</el-button>-->
                        <div class="button-container">
                            <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
                            <el-button type="danger" size="mini" @click="handleDelete(scope.row)">删除</el-button>
                        </div>
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
import {getRecord, addRecord, editRecord, delRecord} from '../api'
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
                content: '',
                attendant: '',
                perform: '',
                satisfaction: ''
            },
            allPageData: [],
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
                addr: [
                    { required: true, message: '请输入地址', trigger: 'blur' }
                ],
                content: [
                    { required: true, message: '请输入服务内容', trigger: 'blur' }
                ],
                attendant: [
                    { required: true, message: '请输入服务人员姓名', trigger: 'blur' }
                ],
                perform: [
                    { required: true, message: '请选择完成情况', trigger: 'blur' }
                ],
                satisfaction: [
                    { required: true, message: '请选择满意度', trigger: 'blur' }
                ],
            },
            tableData:[],
            modelType: 0, //0表示新增弹窗，1表示编辑
            total: 0,  //当前的总条数，默认为0
            pageData:{
                page: 1,
                limit: 10
            },
            recordForm:{
                name: ''
            }
        }
    },
    methods:{
        submit(){
            this.$refs.form.validate((valid) => {
                if (valid){
                    //后续对表单数据的处理，0表示新增弹窗，1表示编辑
                    if (this.modelType === 0){
                        addRecord(this.form).then(() => {
                            //重新获取列表的接口
                            this.getList()
                        })
                    }else{
                        editRecord(this.form).then(() => {
                            //重新获取列表的接口
                            this.getList()
                        })
                    }

                    //清空表单数据
                    this.$refs.form.resetFields()
                    //关闭弹窗
                    this.dialogVisible = false
                }
            })
        },
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
            this.$confirm('此操作将永久删除该记录, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                //delRecord({id: row.id}).then(() => {
                delRecord({name: row.name}).then(() => {
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
        //获取后端列表数据
        getList(){
            getRecord({params: {...this.recordForm, ...this.pageData}}).then(({data}) => {
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
            this.pageData.page=val;
            this.getList()
        },
        //列表查询
        onSubmit(){
            this.getList()
        }
    },
    mounted(){
        //调用getList方法来获取后端传过来的结构体数据并展示在表格中，getList方法中，我们调用了后端的API接口getRecord来获取数据，并将返回的数据赋值给tableData，然后将count属性的值赋值给total，用于分页显示
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

//对齐编辑与删除按钮位置
.button-container {
  display: flex;
  justify-content: flex-start;
}
</style>