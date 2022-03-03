<template>
  <div>
     <table>
        <tbody>
           <tr>
              <td>algo</td>
              <td>algo</td>
              <td>algo</td>
           </tr>
        </tbody>
     </table>
  </div>
</template>

<script>
import axios from 'axios'
export default {
   data:() =>{

   },
   mounted(){
      console.log ('Hola mundo desde mounted')
      this.getAnimals()
   },
methods:{
   async getAnimals(){
      
      await axios.get('https://localhost:9041/api/animales')
      .then(data => console.log(data))
      .catch(e=> console.log(e))
   }
}
}
</script>

<style>

</style>