import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from './service/api'

const app = createApp(App)

// 全局属性
app.config.globalProperties.$axios = axios

// 路由守卫
router.beforeEach((to, from, next) => {
  console.log(to);
  console.log(from);
  if (to.meta.requireAuth) { // 判断该路由是否需要登录权限
    if (localStorage.getItem("loginResult")) { //判断本地是否存在access_token
      next();
    } else {
      if (to.path === '/login') {
        next();
      } else {
        next({
          path: '/login'
        })
      }
    }
  }
  else {
    next();
  }
  /*如果本地 存在 token 则 不允许直接跳转到 登录页面*/
  if (to.fullPath == "/login") {
    if (localStorage.getItem("loginResult")) {
      next({
        path: from.fullPath
      });
    } else {
      next();
    }
  }
})

app.use(router)
app.use(store)
app.mount('#app')
