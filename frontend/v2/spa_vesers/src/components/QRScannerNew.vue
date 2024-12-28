<template>
    <div id="qr-reader-container">
        <p style="color: white; position: absolute; left: 45%; top:45%">Loading...</p>
        <!-- <p style="color: white; position: absolute;  top:55%" id="qr-result"></p> -->
        <div id="qr-reader"></div>
    </div>
</template>
<script setup>
import { onMounted, watch, defineProps, defineEmits } from 'vue'
import { Html5Qrcode } from 'html5-qrcode'
var scanner = null
var scannerRunning = false

const props = defineProps({
    enabled: Boolean,
    counter: Number
})
const emit = defineEmits(["qrDecoded"])

onMounted(() => {
    console.log("QRScanner mounted")

    scanner = new Html5Qrcode("qr-reader");
})
watch(props, async () => {
    /* ... */
    if (props.enabled == true) {
        startScaner()
    } else {
        stopScaner()
    }
})


function callbackQRReader(mutationsList) {
    mutationsList.forEach((mutation) => {
        if (mutation.type === 'childList') {
            let overlayElement = document.getElementById("qr-shaded-region");
            if (overlayElement != null) {
                overlayElement.style = "position: absolute; border-width: 95px; border-style: solid; border-color: #404040ff; box-sizing: border-box; inset: -1%;"
                scannerRunning = true
                return
            }
        }
    });
}
async function startScaner() {
    console.log("Starting QR Scanner")
    const qrCodeSuccessCallback = (decodedText, decodedResult) => {
        emit("qrDecoded", decodedText, decodedResult)
    };
    const qrCodeErrorCallback = (msg, err) => {
        console.err(err, msg);

        /* handle success */
    };
    const config = { fps: 10, qrbox: { width: 200, height: 220 }, aspectRatio: 1 };

    // html5QrCode.start({ facingMode: "environment" }, config, qrCodeSuccessCallback);

    // If you want to prefer front camera
    await scanner.start({ facingMode: "user" }, config, qrCodeSuccessCallback, qrCodeErrorCallback)
    const observerConfig = { childList: true, subtree: true };
    const observerQRReader = new MutationObserver(callbackQRReader);
    let elReader = document.getElementById("qr-reader")
    observerQRReader.observe(elReader, observerConfig);
}

async function stopScaner() {
    console.log("Stopping QR Scanner")
    if (scanner != null && scannerRunning == true) {
        scannerRunning = false
        scanner.pause()
        await scanner.stop()
        scanner.clear()
        console.log("Stopped QR Scanner")
    }
}
</script>


<style>
#reader {
    margin: auto;
    width: 90%;
}

#qr-reader-container {
    padding-top: 60%;
    margin: auto;
    height: 70%;
    background-color: #404040ff;
    /*
    padding-bottom: 10%;
    margin-bottom: 0;
    */
    align-items: center;
}
</style>