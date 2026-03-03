import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/modeltraining/training/detail/:id',
    name: 'TrainingTaskDetail',
    meta: {
      title: '训练任务详情',
      closeTab: true,
      client: true
    },
    component: () => import('@/view/modeltraining/training/detail.vue')
  },
  {
    path: '/modeltraining/training/test/:id',
    name: 'ModelTest',
    meta: {
      title: '模型测试',
      closeTab: true,
      client: true
    },
    component: () => import('@/view/modeltraining/training/modelTest.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
