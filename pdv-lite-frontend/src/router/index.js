// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    // --- Adicionaremos outras rotas aqui depois ---
    {
      path: '/',
      redirect: '/login' // Redireciona a rota raiz para o login
    }
  ]
})

export default router
