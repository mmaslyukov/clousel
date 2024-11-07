<script setup>
import { onMounted, reactive } from 'vue'
import axios from 'axios'
//// eslint-disable-next-line no-unused-vars
const widgetData = reactive({
    current: 0,
    max: 0,
    status: "",
})
const buttons = reactive({
    showRefill: false,
    showPlay: false,
})
var carouselId = "";
//http://localhost:8081/carousel&CarouselId=550e8400-e29b-41d4-a716-446655440000
async function playCarousel(id) {
    // const url = "http://localhost:8081/carousel/play";
    const url = "http://localhost:8081/carousel/play?CarouselId=" + id;
    const response = await axios.get(url);
    if (response.status != 200) {
        return false;
    }
    return true
}

async function refillCarousel(id) {
    const json = JSON.stringify({
        CarouselId: id,
        Rounds: 3,
    });
    const url = "http://localhost:8081/carousel/refill";
    const response = await axios.post(url, json);
    if (response.status != 200) {
        return false;
    }
    return true
}

async function fetchCarousel(id) {
    const url = "http://localhost:8081/carousel?CarouselId=" + id;
    console.log(url);
    var data = {};
    try {
        const response = await axios.get(url);
        if (response.status != 200) {
            return false;
        }
        console.log(response.data);
        data = response.data;
    } catch (e) {
        data.Status = "unknown"
        data.Rounds = 0;
    }
    if (data.Rounds > widgetData.max) {
        widgetData.max = data.Rounds;
    }
    widgetData.current = data.Rounds;

    widgetData.status = data.Status;
    if (widgetData.status != "online") {
        buttons.showRefill = false;
        buttons.showPlay = false;
    } else {
        buttons.showRefill = widgetData.current > 0 ? false : true;
        buttons.showPlay = !buttons.showRefill;
    }
    return true;



}

onMounted(() => {
    console.log("App mounted");
    widgetData.current = 0;
    widgetData.max = 0;
    buttons.showRefill = false;
    buttons.showPlay = false;
    carouselId = window.location.href.split('/').at(-1);
    fetchCarousel(carouselId); // TODO
});

function onClick() {
    if (widgetData.current > 0) {
        widgetData.current--;
    }
    buttons.showRefill = widgetData.current > 0 ? false : true;
    buttons.showPlay = !buttons.showRefill;
}

function onEventPlayClick() {
    // buttons.showRefill = !buttons.showRefill;
    console.log("play clicked ");
    playCarousel(carouselId).then(function (result) {
        if (result) {
            setTimeout(() => {
                fetchCarousel(carouselId)
            }, 1000);
        }

    });
}

function onEventRefillClick() {
    console.log("refill clicked");
    refillCarousel(carouselId).then(function (result) {
        if (result) {
            setTimeout(() => {
                fetchCarousel(carouselId)
            }, 1000);
        }

    });
}

</script>


<script>
import WidgetRounds from "./components/WidgetRounds.vue";
import GeneralButton from './components/GeneralButton.vue';
// import StripeForm from './components/StripeForm.vue';
export default {
    name: "App",
    components: {
        // StripeForm,
        WidgetRounds,
        GeneralButton,
    },
};

</script>

<template>
    <StripeForm />
    <WidgetRounds v-bind="widgetData" @click="onClick" />
    <GeneralButton v-bind="buttons" @play-clicked="onEventPlayClick" @refill-clicked="onEventRefillClick" />
</template>

<style>
/* #app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
} */
</style>
