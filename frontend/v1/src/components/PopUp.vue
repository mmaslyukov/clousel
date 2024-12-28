<script setup>
import { defineProps, watch, onMounted, nextTick } from 'vue'
import { PopupSuccess, PopupError } from './popup.js'
const props = defineProps(['show', 'type', 'msg'])

// function drawSuccessPopup(msg) {
//     let bc = "#e2f9e2ff"
//     let sc = "#2e8e37"
//     let canvas = document.getElementById("popup_id");
//     var parent = document.getElementById("popup_canvas_container_id");
//     canvas.width = parent.offsetWidth;
//     //     canvas.height = parent.offsetHeight;
//     let ctx = canvas.getContext("2d");
//     ctx.beginPath();
//     ctx.fillStyle = bc
//     ctx.strokeStyle = sc
//     ctx.lineWidth = 4
//     ctx.roundRect(canvas.width * 0.05, 10, canvas.width * 0.9, 100, 10)
//     ctx.fill()
//     ctx.stroke();

    
//     ctx.beginPath()
//     ctx.fillStyle = sc
//     let D = canvas.height / 6
//     ctx.arc(canvas.width * 0.15, canvas.height / 2 - D/2, D, 0, Math.PI * 2, true);
//     ctx.fill()
    
//     let fhcm = 24
//     ctx.font = fhcm + "pt Arial";
//     ctx.fillStyle = bc
//     ctx.textAlign = "left";
//     ctx.fillText("✔", canvas.width * 0.125, canvas.height /2 - fhcm*0.1 )

//     let fhh = 20
//     ctx.font = "bold " + fhh + "pt Arial";
//     ctx.fillStyle = sc;
//     ctx.textAlign = "left";
//     ctx.fillText("SUCCESS", canvas.width * 0.25, canvas.height*0.37 )

//     let fh = 16
//     ctx.font = fh + "pt Arial";
//     ctx.fillStyle = sc;
//     ctx.textAlign = "left";
//     ctx.fillText(msg, canvas.width * 0.25, canvas.height*0.55)
// }



function drawPopup(type, msg) 
    {
    let canvas = document.getElementById("popup_id")
    var parent = document.getElementById("popup_canvas_container_id")
    canvas.width = parent.offsetWidth
    canvas.height = parent.offsetHeight
    if (canvas == undefined) {
        console.error("Canvas is undefined")
        return
    }
    if (type == "popup_success") {
        new PopupSuccess(canvas, canvas.width, canvas.height).msg(msg).redraw()
    } else if (type == "popup_error") {
        new PopupError(canvas, canvas.width, canvas.height).msg(msg).redraw()
    }

}

onMounted(() => {
    // console.log(">>>> mounted.PopUp ", props)
})
watch(props, async (newProps) => {
    await nextTick();
    // console.log(">>>> watch.PopUp ", newProps)
    drawPopup(newProps.type, newProps.msg)
    // drawPopup("Payment has been confirmed")

})



</script>


<template>
    <!-- <div id="popup_canvas_container_id" class="popup_canvas_container"> -->
    <div id="popup_canvas_container_id" class="popup_canvas_container" v-show="props.show">
        <canvas id="popup_id" class="popup_canvas"></canvas>
    </div>
</template>

<style>
.popup_canvas {
    margin: auto;
    /* background-color: grey; */
}

.popup_canvas_container {
    text-align: center;
    /* height: 100%; */
    height: 120px;
    left: 0;
    right: 0;
    /* position:absolute; */
    margin: auto;
}

/* .popup_canvas {
    width: 100%;
    height: 100%;
} */

/* .popup_canvas_container {
    padding-top: 40px;
    max-width: fit-content;
    margin: auto;
} */
</style>