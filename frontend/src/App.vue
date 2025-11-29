<script setup>
import HelloWorld from './components/Home.vue'
import Navbar from './components/Navbar.vue';
import { RouterView } from 'vue-router';
</script>

<template>
  <Navbar />
  <RouterView />
</template>
