<script setup lang="ts">
import axios from 'axios'
const ip = axios.get("http://106.74.7.3:8080/ip").then(resp=>resp.data).catch(err=>console.log(err))
</script>

<template>
  <h1>{{ ip }}</h1>
</template>

<style scoped>
a {
  color: #42b983;
}

label {
  margin: 0 0.5em;
  font-weight: bold;
}

code {
  background-color: #eee;
  padding: 2px 4px;
  border-radius: 4px;
  color: #304455;
}
</style>
