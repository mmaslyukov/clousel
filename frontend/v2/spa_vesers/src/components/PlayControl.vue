<template>
    <div style="height: inherit; background-color: beige; display: flex; flex-direction: column;">
        <div style=" align-items: center; justify-content: center; vertical-align: middle; display: flex; flex-direction: column; margin:auto; margin-bottom: 0;  padding:0;">
            <span>
                <p class="pc-font" style="font-size: 14pt; text-align: center; margin:0; margin-bottom: 10px;">Play fee:</p>
            </span>
            <span style="display: flex; flex-direction: row; ">
                <p class="pc-font pc-bold" style="font-size: 22pt; text-align: right; margin:0; margin-right: 5px; " >{{ props.fee }}</p>
                <!-- <p class="pc-font" >{{ props.fee }}</p> -->
                <p class="pc-font pc-bold" style="font-size: 22pt; margin:0;">Tickets</p>
            </span>
            <span>
                <p v-show="!props.enable" class="pc-font" style="font-size: 10pt; color: #c14b1dff; margin:0; margin-top: 10px;" >Not enought tickets</p>
            </span>
        </div>
        <div style="display: flex; flex-direction: row;align-items: center; justify-content: center; margin:auto; margin-top: 0; " >
            <button id="pc-play" class="pc-font pc-bold pc-button-play-off" type="button" disabled @click="onPlayClick" >PLAY</button>
        </div>
    </div>
</template>
<script setup>
import { onMounted, watch, defineProps, defineEmits } from 'vue'
const emit = defineEmits(["playClicked"])

const props = defineProps({
    fee: Number,
    enable: Boolean
})
onMounted(() => {
})

watch(props, () => {
    let btn = document.getElementById("pc-play");
    if (btn != null) {
        if (props.enable == true) {
            btn.classList.remove("pc-button-play-off")
            btn.classList.add("pc-button-play-on")
            btn.disabled = false
        } else {
            btn.classList.remove("pc-button-play-on")
            btn.classList.add("pc-button-play-off")
            btn.disabled = true
        }
    }
})

function onPlayClick(){
    emit('playClicked')
    console.log("onPlayClick")
}

</script>


<style>
.pc-bold {
    font-weight: bold;
}
.pc-font {
    font-family: sans-serif;
    color: #404040ff;
}


.pc-button-play-off,
.pc-button-play-on,
.pc-button-play {
    border: none;
    width: 120px;
    height: 50px;
    margin-top: 40px;
    font-weight: bold;
    font-size: 20pt;
    border-radius: 10px;
}
.pc-button-play-on {
    background-color: #404040ff;
    color: #ffeb99ff;
}
.pc-button-play-on:active {
    background-color: #1a1a1aff;
    color: #cca400ff;
}

.pc-button-play-off {
    background-color: #404040ff;
    color: #808080ff;
}

</style>