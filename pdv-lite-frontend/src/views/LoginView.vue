<script setup>
import { ref } from 'vue'

// Cria duas "caixinhas" reativas para guardar os valores dos inputs.
// Elas começam com uma string vazia.
const email = ref('')
const password = ref('')

// Esta função será executada quando o formulário for enviado.
const handleLogin = () => {
  // Por enquanto, vamos apenas mostrar os valores no console do navegador
  // para confirmar que estamos pegando os dados corretamente.
  // Usamos .value para acessar o valor que está dentro de uma ref.
  console.log('Tentando logar com:', email.value, password.value)
}
</script>

<template>
  <div>
    <h1>PDV LITE</h1>
    <h2>Login</h2>

    <form @submit.prevent="handleLogin">
      <div>
        <label for="email">E-mail</label>
        <input type="email" id="email" v-model="email" />
      </div>
      <div>
        <label for="password">Senha</label>
        <input type="password" id="password" v-model="password" />
      </div>
      <button type="button" @click="handleLogin">Entrar (Modo Teste)</button>
    </form>
  </div>
</template>

<style scoped>
/* Estilos específicos para esta página virão aqui */
</style>
