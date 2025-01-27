<template>
  <SignPage v-show="signRequired" @registered="onRegistered" @logged-in="onLoggedIn" />
  <div v-show="!signRequired" style="height: 100%;">
    <HeaderToolbar @po-changed="onPoChanged" @logo-clicked="logoClicked" v-bind="tickets" />
    <QRScannerNew v-show="viewMain == viewNameScanner" v-bind="scannerView" @qr-decoded="onQrDecoded" />
    <HelpScreen v-show="viewMain == viewNameHelp" />
    <PaymentOptions v-show="viewOverlay == viewNameOverlay" v-bind="priceOptions" @buy-clicked="onBuyClicked" />
    <!-- <PlayControl v-bind="play" @play-clicked="onPlayClicked" /> -->
    <PlayControl v-show="viewMain == viewNameMachine" v-bind="play" @play-clicked="onPlayClicked" />
    <FooterToolbar @scanner-enabled="onScannerEnabled" v-bind="scannerFeedback" />
  </div>
</template>

<script>
import HeaderToolbar from './components/HeaderToolbar.vue'
import FooterToolbar from './components/FooterToolbar.vue'
import SignPage from './components/SignPage.vue'
import QRScannerNew from './components/QRScannerNew.vue'
import PaymentOptions from './components/PaymentOptions.vue'
import PlayControl from './components/PlayControl.vue'
import HelpScreen from './components/HelpScreen.vue'

export default {
  name: 'App',
  components: {
    SignPage,
    HeaderToolbar,
    PaymentOptions,
    HelpScreen,
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
// import cookiejar from 'cookiejar'
import { config } from './utils/config.js'

import axios from 'axios'
// import cookiejar from 'cookiejar'

const viewNameHelp = "help"
const viewNameOverlay = "pay"
const viewNameScanner = "scanner"
const viewNameMachine = "machine"
var tocken = undefined
const signRequired = ref(true)
const poShown = ref(false)
const viewMain = ref(viewNameHelp)
const viewOverlay = ref("")
const scannerView = reactive({
  enabled: false,
  counter: 0
})
const tickets = reactive({
  balance: 0,
  closeOptions: 0
})
const play = reactive({
  cost: 0,
  balance: 0,
  status: undefined,
  machId: undefined,
  enable: false,
  unable: false,
  unableMsg: ""

})
const priceOptions = reactive({
  details: [
    // {
    //   tickets: 2,
    //   price: 1,
    //   priceId: "pid1",
    // },
    // {
    //   tickets: 5,
    //   price: 2,
    //   priceId: "pid2",
    // },
    // {
    //   tickets: 10,
    //   price: 3.50,
    //   priceId: "pid3",
    // },
    // {
    //   tickets: 15,
    //   price: 5,
    //   priceId: "pid4",
    // },
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
// var machineId = ""
onMounted(() => {
  let c = Object.fromEntries(document.cookie.split('; ').map(v => v.split(/=(.*)/s).map(decodeURIComponent)))
  tocken = c['Tocken']
  // let machId = window.location.href.split('/').at(-1).split('?').at(0);
  console.log("App is mounted")

  fetchBalance().then(function (res) {
    signRequired.value = !res
    fetchPrices()
    machinePlayView(window.location.href)
  })

})

// async function fetchBalance() {

// }
// async function fetchPrices() {

// }

function updatePlayEnabledMarker() {
  if (tickets.balance < play.cost) {
    play.enable = false
    play.unableMsg = "Not enought tickets"
    // play.unable = true
  } else if (play.status != "online") {
    play.enable = false
    play.unableMsg = `Unexpected machine status: "${play.status}"`
  } else {
    play.enable = true
    play.unableMsg = ""
  }
  console.log(`Enabled:${play.enable}, Msg:"${play.unableMsg}"`)
  // play.enable = (tickets.balance >= play.cost) && (play.status == "online")
}

async function onPlayClicked() {
  if (tickets.balance >= play.cost) {

    let url = config.url().play(tocken, play.machId)
    try {
      const response = await axios.get(url);
      if (response.status != 200) {
        throw new Error("Fail to fetch balance")
      }
      let eventId = response.data.EventId
      play.enable = false
      let tid = setInterval(async () => {
        try {
          let url = config.url().poll(tocken, eventId)
          let response = await axios.get(url)
          if (response.status != 200) {
            throw new Error("Fail to fetch balance")
          }
          if (response.data.Status != "pending") {
            clearInterval(tid)
          }
          fetchBalance()
        } catch (e) {
          console.error(e)
        }
      }, 500);

    } catch (e) {
      if (e.status == 401) {
        signRequired.value = true
      } else {
        console.error(e)
      }
    }

    // tickets.balance -= play.cost
    // play.balance = tickets.balance
    // updatePlayEnabledMarker()
  }
}

async function onBuyClicked(selectedPriceOption) {

  // price:Number,
  //   tickets: Number,
  //   priceId:String,

  const home = window.location.href.split('?').at(0)

  let url = config.url().buyTickets(tocken, home, selectedPriceOption.priceId)
  try {
    const response = await axios.get(url);
    if (response.status != 200) {
      throw new Error("Fail to fetch balance")
    }
    window.location = new DOMParser().parseFromString(response.data, "text/html").querySelector("a").getAttribute("href")

  } catch (e) {
    if (e.status == 401) {
      signRequired.value = true
    } else {
      console.error(e)
    }
  }
}

async function onQrDecoded(txt, result) {
  console.log("QR detected ", txt)
  scannerFeedback.detected = true
  scannerFeedback.counter++;
  scannerView.enabled = false
  scannerView.counter++;
  machinePlayView(txt)

  // play.cost = 1
  // updatePlayEnabledMarker()
}
function onScannerEnabled(value) {
  scannerView.enabled = value
  if (scannerView.enabled) {
    viewMain.value = viewNameScanner
  } else {
    viewMain.value = viewNameHelp
  }
  console.log("onScannerEnabled", scannerView.enabled)
}
function onRegistered() {
  console.log("onRegistered")
  // signRequired.value = false
  onLoggedIn()
}
function onLoggedIn() {
  let c = Object.fromEntries(document.cookie.split('; ').map(v => v.split(/=(.*)/s).map(decodeURIComponent)))
  tocken = c['Tocken']
  if (tocken != undefined) {
    fetchBalance()
    fetchPrices()
    signRequired.value = false
    machinePlayView(window.location.href)


  }
}

async function fetchBalance() {

  let url = config.url().balance(tocken)
  try {
    const response = await axios.get(url);
    if (response.status != 200) {
      throw new Error("Fail to fetch balance")
    }
    tickets.balance = response.data.Balance
    play.balance = response.data.Balance
    updatePlayEnabledMarker()
    return true

  } catch (e) {
    console.error(e)
  }
  return false
}

async function machinePlayView(address) {

  // console.log(address)
  // console.log(address.split('/').at(-1))
  // console.log(address.split('/').at(-1).split('?').at(0))
  let machId = address.split('/').at(-1).split('?').at(0);
  console.log("machine id", machId)
  if (machId == undefined || machId.length == 0) {
    return
  }
  let url = config.url().pick(tocken, machId)
  console.log(`fetch ${machId}`)
  try {
    const response = await axios.get(url);
    if (response.status != 200) {
      throw new Error("Fail to fetch machine")
    }
    play.machId = response.data.MachId
    play.status = response.data.Status
    play.cost = response.data.Cost
    viewMain.value = viewNameMachine
    updatePlayEnabledMarker()

  } catch (e) {
    console.error(e)
    viewMain.value = viewNameHelp
  }
  console.log("veiw", viewMain.value)
}



// console.log(c)
// console.log(c['Tocken'])
// let jar = cookiejar.CookieJar()
// // console.log(`App cookies: ${jar.getCookies()}`)
// jar.getCookie("Tocken")
// let t = new cookiejar.Cookie()
//   t.name = "Tocken"
//   t.value = data.Tocken
//   console.log(data)
//   console.log(data.Tocken)
//   if (data.Tocken.length > 0) {
//       emit('LoggedIn')
//   }
//   document.cookie = t.toString()
// console.log(`Cookie Tocken: ${tocken}`)

async function logoClicked() {
  viewMain.value = viewNameHelp
}

async function fetchPrices() {
  let url = config.url().prices(tocken)
  try {
    const response = await axios.get(url);
    if (response.status != 200) {
      throw new Error("Fail to login")
    }
    console.log(response.data)
    priceOptions.details = []
    for (let e of response.data) {
      console.log(e)
      priceOptions.details.push({
        tickets: e.Tickets,
        price: e.Amount * 0.01,
        priceId: e.PriceId
      })
    }
    // Amount  int    `json:"Amount"`
    // Tickets int    `json:"Tickets"`
    // PriceId string `json:"PriceId"`
    // tickets: 2,
    // price: 1,
    // priceId: "pid1",

  } catch (e) {
    console.error(e)
  }
}

async function onPoChanged(visible) {
  viewOverlay.value = visible ? viewNameOverlay : ""
  if (visible) {
    await fetchPrices()
  }
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
  font-family: Verdana, Geneva, Tahoma, sans-serif;
  color: #404040ff;
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
