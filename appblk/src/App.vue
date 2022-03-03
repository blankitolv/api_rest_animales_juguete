<template>
  <div class="div">
    <img alt="Vue logo" src="./assets/logo.png">
    <div class="div">
      <ButtonAPI />
      <GetAPI />
    </div>
  </div>
</template>

<script> 
import ButtonAPI from "./components/ButtonAPI";
import GetAPI from "./components/GetApi"

export default {
  name: 'App',
  components: {
    ButtonAPI,
    GetAPI
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
td {
  display: flex;
  flex-direction: column;
}
</style>
