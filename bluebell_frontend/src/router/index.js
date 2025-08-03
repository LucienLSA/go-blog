import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Content from '../views/Content.vue'
import Publish from '../views/Publish.vue'
import Login from '../views/Login.vue'
import SignUp from '../views/SignUp.vue'
import UserInfo from '../views/UserInfo.vue'
import ValidEmail from '../views/ValidEmail.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/post/:id',
    name: 'Content',
    component: Content
  },
  {
    path: '/publish',
    name: 'Publish',
    component: Publish,
    meta: { requireAuth: true }
  },
  {
    path: '/login',
    name: "Login",
    component: Login
  },
  {
    path: '/signup',
    name: "SignUp",
    component: SignUp
  },
  {
    path: '/userinfo',
    name: 'UserInfo',
    component: UserInfo,
    meta: { requireAuth: true }
  },
  {
    path: '/user/validEmail/:token',
    name: 'ValidEmail',
    component: ValidEmail
  }
]

const router = createRouter({
  history: createWebHashHistory(process.env.BASE_URL),
  routes
})

export default router
