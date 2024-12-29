<template>
    <div id="qr-reader-container">
        <!-- <p  v-show="loading" style=" color: white; text-align: center;  margin:auto; padding:auto; ">Loading...</p> -->
        <p  v-show="loading" style=" color: white; position: absolute; top:45%; left: 45%;" >Loading...</p>
        <div id="qr-reader"></div>
    </div>
</template>
<script setup>
import { onMounted, watch, defineProps, defineEmits, ref } from 'vue'
import { Html5Qrcode } from 'html5-qrcode'
var scanner = null
var scannerRunning = false
const loading = ref(false)

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
    loading.value = false
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
    loading.value = true
    console.log("Starting QR Scanner")
    const qrCodeSuccessCallback = (decodedText, decodedResult) => {
        emit("qrDecoded", decodedText, decodedResult)
    };
    const qrCodeErrorCallback = (msg, err) => {
        console.err(err, msg);

        /* handle success */
    };
    const config = { fps: 10, qrbox: { width: 200, height: 200 }, aspectRatio: 1 };

    // html5QrCode.start({ facingMode: "environment" }, config, qrCodeSuccessCallback);

    // If you want to prefer front camera
    // await scanner.start({ facingMode: "user" }, config, qrCodeSuccessCallback, qrCodeErrorCallback)
    await scanner.start({ facingMode: "environment" }, config, qrCodeSuccessCallback, qrCodeErrorCallback)
    const observerConfig = { childList: true, subtree: true };
    const observerQRReader = new MutationObserver(callbackQRReader);
    let elReader = document.getElementById("qr-reader")
    observerQRReader.observe(elReader, observerConfig);
}

async function stopScaner() {
    console.log("Stopping QR Scanner")
    if (scanner != null && scannerRunning == true) {
        loading.value = false
        scannerRunning = false
        scanner.pause()
        await scanner.stop()
        scanner.clear()
        console.log("Stopped QR Scanner")
    }
}
</script>


<style>
video {
    height: 400px;
}
#qr-reader {
    min-width: 200px;
    max-width: 500px;
    height: 400px;
    width: 400px;
    /*
    */
    margin: auto;
    padding: auto;
    top:20px;
    bottom: 10%;
    
    
    /*
    display: flex;
    vertical-align: middle;
    position: relative;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    */
}

#qr-reader-container {
    margin: auto;
    background-color: #404040ff;
    padding-top: 110px;
    padding-bottom: 110px;
    height: 100%;
    /*
    display: flex;
    align-items: center;
    justify-content: center;
    padding-bottom: 10%;
    margin-bottom: 0;
    vertical-align: middle;
    */
}
</style>