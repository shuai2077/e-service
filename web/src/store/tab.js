import Cookie from "js-cookie"
//管理菜单
export default{
    state:{
        isCollapse: false,  //用于控制菜单的展开和收起
        tabsList: [
            {
                path: '/',
                name: 'record',
                label: '服务记录管理',
                icon: 'document',
                url: 'Record/Record'
            }
        ], //面包屑的数据:点了哪个路由,首页是一定有的
        menu:[]
    },
    mutations:{
        //修改菜单展开收起的方法
        collapseMenu(state){
            state.isCollapse = !state.isCollapse
        },
        // 更新面包屑的数据
        // selectMenu(state, val) {
        //     // 如果点击的不在面包屑数据中,则添加
        //     if (val.name !== 'record'){
        //         const index = state.tabsList.findIndex(item => item.name === val.name)
        //         if (index === -1) {
        //             state.tabsList.push(val)
        //         }
        //     }
        // },
        //设置menu的数据
        setMenu(state, val){
            state.menu=val
            Cookie.set('menu',JSON.stringify(val))
        },
        addMenu(state,router){
            //判断缓存中是否有数据
            if(!Cookie.get('menu')){
                return
            }
            const menu = JSON.parse(Cookie.get('menu'))
            state.menu = menu
            //组装动态路由的数据
            const menuArray = []
            menu.forEach(item => {
                if(item.children){
                    item.children = item.children.map(item => {
                        item.component = () => import(`../views/${item.url}`)
                        return item
                    })
                    menuArray.push(...item.children)
                }else{
                    item.component = () => import(`../views/${item.url}`)
                    menuArray.push(item)
                }
            })
            //路由的动态添加
            menuArray.forEach(item => {
                router.addRoute('Main',item)
            })
        }
    }
}