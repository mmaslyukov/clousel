<template>
  <SignPage v-show="signRequired" @register-clicked="onRegisterClicked" @login-clicked="onLoginClicked" />
  <div v-show="!signRequired" style="height: 100%;">
    <HeaderToolbar @po-changed="onPoChanged" v-bind="tickets" />
    <QRScannerNew v-show="scannerView.enabled" v-bind="scannerView" @qr-decoded="onQrDecoded" />
    <PaymentOptions v-show="poShown" v-bind="priceOptions" @buy-clicked="onBuyClicked" />
    <PlayControl  v-bind="play" @play-clicked="onPlayClicked" />
    <!-- <PlayControl v-show="scannerFeedback.detected" v-bind="play" @play-clicked="onPlayClicked" /> -->
    <FooterToolbar @scanner-enabled="onScannedEnabled" v-bind="scannerFeedback" />
  </div>
</template>

<script>
import HeaderToolbar from './components/HeaderToolbar.vue'
import FooterToolbar from './components/FooterToolbar.vue'
import SignPage from './components/SignPage.vue'
import QRScannerNew from './components/QRScannerNew.vue'
import PaymentOptions from './components/PaymentOptions.vue'
import PlayControl from './components/PlayControl.vue'

export default {
  name: 'App',
  components: {
    SignPage,
    HeaderToolbar,
    PaymentOptions,
    PlayControl,
    QRScannerNew,
    FooterToolbar
  },
  // mounted() {
  //   const script = document.createElement("script");
  //   script.src = "https://unpkg.com/html5-qrcode";
  //   document.body.appendChild(script);
  // },
}
</script>

<script setup>
import { onMounted, reactive, ref } from 'vue'

const signRequired = ref(true)
const poShown = ref(false)
const scannerView = reactive({
  enabled: false,
  counter: 0
})
const tickets = reactive({
  balance: 0
})
const play = reactive({
  fee: 0,
  enable: false
})
const priceOptions = reactive({
  details: [
    {
      tickets: 2,
      price: 1,
      prodId: "pid1",
    },
    {
      tickets: 5,
      price: 2,
      prodId: "pid2",
    },
    {
      tickets: 10,
      price: 3.50,
      prodId: "pid3",
    },
    {
      tickets: 15,
      price: 5,
      prodId: "pid4",
    },
  ]
})
const scannerFeedback = reactive({
  detected: false,
  counter: 0
})
// const props = defineProps({
//     disable: Boolean
// })
// const props = defineProps({
//     deactivateScanner: Boolean
// })

// var scannerShown = ref(false)
onMounted(() => {
  console.log("App is mounted")
  // scanner.scannerEnabled = true
})

function updatePlayEnabledMarker() {
  play.enable = tickets.balance >= play.fee

}

function onPlayClicked() {
  if (tickets.balance >= play.fee) {
    tickets.balance -= play.fee
    updatePlayEnabledMarker()
  }
}

function onBuyClicked(selectedPriceOption) {
  tickets.balance = selectedPriceOption.tickets
  updatePlayEnabledMarker()

}

function onQrDecoded(txt, result) {
  console.log("QR detected ", txt)
  scannerFeedback.detected = true
  scannerFeedback.counter++;
  scannerView.enabled = false
  scannerView.counter++;

  play.fee = 1

  updatePlayEnabledMarker()
}
function onScannedEnabled(value) {
  scannerView.enabled = value
  console.log("onScannedEnabled", scannerView.enabled)
}
function onRegisterClicked() {
  console.log("onRegisterClicked")
  signRequired.value = false
}
function onLoginClicked() {
  console.log("onLoginClicked")
  signRequired.value = false
}
function onPoChanged(visible) {
  poShown.value = visible
  console.log("visble " + visible)

}
</script>


<style>
#app {
  height: 100%;
  display: flex;
  flex-direction: column;
}

html,
body {
  padding: 0 !important;
  margin: 0 !important;
  height: 100%;
  overflow: clip;
  background-color: beige;
}

/*
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
  */
</style>
