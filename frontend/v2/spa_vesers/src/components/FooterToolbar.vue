<template>
    <div id="footer">
        <div id="footer-toolbar">
            <div id="footer-scan">
                <canvas id="footer-scan-canvas" width="70px" height="70px"></canvas>
            </div>
            <!-- <div id="footer_scan">
                <p>scan</p>
            </div> -->
        </div>
    </div>
</template>

<script setup>
import { onMounted, defineEmits, defineProps, watch } from 'vue'
import { HSLColor } from './HSLColor.js'
import { CanvasButtonScan } from './CanvasButton.js'
const emit = defineEmits(["scannerEnabled"])
const props = defineProps({
    detected: Boolean,
    counter: Number
})
var buttonScanner = null

var buttonArray = []

watch(props, async () => {
    if (props.detected == true) {
        console.log("QR detected, turning off button")
        buttonScanner.turnOff().redraw()
    }
})
onMounted(() => {
    {
        let c = document.getElementById("footer-scan-canvas");
        buttonScanner = new CanvasButtonScan(c, c.width, c.height, new HSLColor(48, 100, 53), new HSLColor(0, 0, 25)).redraw()
        buttonArray.push(buttonScanner)

    }
    registerClickers()
})

function registerClickers() {
    for (let b of buttonArray) {
        b.canvas().addEventListener('click', () => {
            b.onClickApplyVisual().redraw()
            emit("scannerEnabled", b.isEnbaled())
        });
    }
}
</script>


<style>
#footer {
    position: absolute;
    margin-top: auto;
    height: 100px;
    bottom: 0;
    width: 100%;
    background-color: #ffcc00ff;
    box-shadow: 0px -5px 10px #404040ff;
    /*
    */
}

#footer-toolbar {
    /**
    border-color: red;
    */
    height: 100%;
    display: flex;

}

#footer-scan {
    width: 70px;
    height: 70px;
    margin: auto;
}

#footer-scan-canvas {
    margin: auto;
    padding: auto;
    box-shadow: 0px 0px 10px #404040ff;

}
</style>