import http from '../utils/request'

// 请求首页数据,直接把这个对象导出
export const getData = () => {
    // 返回一个promise对象
    return http.get('/record/getData')
}

//用户列表
export const getUser=(params) => {
    //返回用户列表
    return http.get('/user/getUser',params)
}

export const addUser=(data) => {
    return http.post('/user/add',data)
}

export const editUser=(data) => {
    return http.post('/user/edit',data)
}

export const delUser=(data) => {
    return http.post('/user/del',data)
}

//记录列表
export const getRecord=(params) => {
    return http.get('/record/getRecord',params)
}

export const addRecord=(data) => {
    return http.post('/record/add',data)
}

export const editRecord=(data) => {
    return http.post('/record/edit',data)
}

export const delRecord=(data) => {
    return http.post('/record/del',data)
}

//历史列表
export const getHis=(params) => {
    //返回历史列表
    return http.get('/his/getHis',params)
}

// 登录权限
export const getMenu = (data) => {
    return http.post('/permission/getMenu',data)
}