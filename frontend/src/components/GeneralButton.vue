<script setup>
import { defineProps, defineEmits, onMounted } from 'vue'
// import  {} from 'vue'
defineEmits(["playClicked", "refillClicked"])
const props = defineProps({
    showRefill: Boolean,
    showPlay: Boolean,
})

onMounted(() => {
    // console.log("GeneralButton is Mounted");
    // console.log("showRefill:", props.showRefill.value, " showPlay: ", props.showPlay.value);
})

</script>


<template>
    <div v-show="props.showPlay" class="button_container">
        <button  type="button" class="btn btn__play" id="PlayBtn" @click="$emit('playClicked')">PLAY</button>
    </div>
    <div v-show="props.showRefill" class="button_container">
        <button  type="button" class="btn btn__refill" id="RefillBtn" @click="$emit('refillClicked')">REFILL</button>
    </div>
</template>

<style>
.btn {
    border: none;
    color: white;
    padding: 25px 50px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    margin: 10px 0px;
    cursor: pointer;
    font-family: "Impact";
    font-size: 60px;
    min-width: 20rem;
    /* text-align: center; */
}
.btn__refill {
  background-color: #71db7aff;
}

.btn__play {
  background-color: #71a6dbff;
}

.btn__refill:active {
  background-color: #39ac43ff;
}


.btn__play:active {
    background-color: #2b6baaff;
}
.button_container {
    max-width: fit-content;
    margin: auto;
    /* border-style: solid; */
}
</style>