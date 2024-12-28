<script setup>
import { onMounted, reactive } from 'vue'
import axios from 'axios'

const URLCarouselService = "http://localhost:8081"
const URLAccountantService = "http://localhost:4321"

// const URLCarouselService = "http://109.110.19.178:8081"
// const URLAccountantService = "http://109.110.19.178:4321"

const widgetData = reactive({
    current: 0,
    max: 0,
    status: "",
    show: true,
})
const buyOptData = reactive({
    menu: [],
    show: false,
})
const popup = reactive({
    show: false,
    type: "",
    msg: "",
})
const buttons = reactive({
    showRefill: false,
    showPlay: false,
})
var carouselId = "";

async function playCarousel(id) {
    var details = {
        'CarouselId': id,
    };
    const formBody = Object.entries(details).map(([key, value]) => encodeURIComponent(key) + '=' + encodeURIComponent(value)).join('&')


    const url = URLCarouselService+"/carousel/play";
    const response = await axios.post(url, formBody);
    if (response.status != 200) {
        return false;
    }
    return true
}

async function refillCarousel(id, tickets) {
    var details = {
        'CarouselId': id,
        'Tickets': tickets,
    };
    const formBody = Object.entries(details).map(([key, value]) => encodeURIComponent(key) + '=' + encodeURIComponent(value)).join('&')

    const url = URLCarouselService+"/carousel/refill";
    const response = await axios.post(url, formBody);
    if (response.status != 200) {
        return false;
    }
    return true
}
async function checkout(cid, pid) {
    const home = window.location.href.split('?').at(0)
    const url = URLAccountantService+"/client/checkout?CarouselId=" + cid + "&PriceId=" + pid + "&HomeUrl=" + encodeURI(home);

    
    const response = await axios.get(url);
    if (response.status != 200) {
        return false;
    }
    window.location = new DOMParser().parseFromString(response.data, "text/html").querySelector("a").getAttribute("href")

        // let data = response.data;


        // console.log("url: ", data.Url)
        // if (response.status == 303) {
        // console.log("rurl ", response.url)
        // window.location = response.url
}


async function fetchCarousel(id) {
    const url = URLCarouselService + "/carousel?CarouselId=" + id;
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
        data.Status = ""
        data.Tickets = 0;
    }

    if (data.Status.length == 0) {
        data.Status = "unknown"
    }

    if (data.Tickets > widgetData.max) {
        widgetData.max = data.Tickets;
    }
    widgetData.current = data.Tickets;

    widgetData.status = data.Status;
    console.log("data.Status  ", data.Status)
    console.log("widgetData.status  ", widgetData.status)
    return true;

}

async function fetchPriceOptions(id) {
    const url = URLAccountantService + "/client/prices?CarouselId=" + id;
    console.log(url);
    var data = [];

    try {
        const response = await axios.get(url);
        if (response.status != 200) {
            return false;
        }
        console.log(response.data);
        data = response.data;
    } catch (e) {
        data = []
        console.error(e)
    }
    buyOptData.menu = []
    for (let price of data) {
        console.log(price)
        buyOptData.menu.push({ price: price.Amount / 100, tickets: price.Tickets, priceId: price.PriceId })
    }

    // if (buyOptData.menu.length > 0) {
    //         buyOptData.show = true
    // }

    // if (data.length > 0) {


    //     buyOptData.menu = [{price:2, tickets:1, priceId:"pid_#1"}, {price:4, tickets:3,priceId:"pid_#2"}, {price:5, tickets:5, priceId:"pid_#3"}]

    // }


}

async function updateComponentsView(widget, buttons, buyopt, popup) {
    if (popup.show) {
        widget.show = false
        buyopt.show = false
    } else if (widget.status == "online") {
        widget.show = true
        if (widget.current > 0) {
            buttons.showPlay = true
            buyopt.show = false
        } else {
            buttons.showPlay = false
            if (buyopt.menu != undefined && buyopt.menu.length > 0) {
                buyopt.show = true
            }
        }
    }
}

function getQueryVariables(url) {
    
    let args = []
    if (url.indexOf('?') > 0) {
        url.split('?').at(-1).split('&').forEach(pair => {
            let tup = pair.split('=')
            if (tup != undefined && tup.length > 0) {
                args.push({
                    name: tup[0],
                    // value: decodeURIComponent(tup[1]),
                    value: decodeURIComponent(tup[1].replace(/\+/g, '%20'))//   (tup[1]),
                })
            }
        });
    } 
    return args
}

onMounted(() => {
    // let astr = '<a href="checkout.stripe.com/c/pay/cs_test_a1jkFjpOsbyoTt7ZPhCtgc1dL0XhEefXlxXpbgdKswT9krt9UaVzMs82Xb#fidkdWxOYHwnPyd1blpxYHZxWjA0VWRvXWpXcGd1VmlCVm59RHQwR2BVcFBET2d8NGdNTkFLZ19%2FNE9rZzFBRFZhMV03NE1dc05GNlNobH12PEtiZmxyUmp8TmFKcUhvN21cdUNhPEFHQUNwNTV8aHdxaF1IPCcpJ2N3amhWYHdzYHcnP3F3cGApJ2lkfGpwcVF8dWAnPyd2bGtiaWBabHFgaCcpJ2BrZGdpYFVpZGZgbWppYWB3dic%2FcXdwYHgl">OK</a>'
    // var a = document.createElement('div');
    // a.innerHTML = astr
    // console.log("a >> ", a)
    // console.log("href >> ", a.firstChild.getAttribute('href'))

    // const parser = new DOMParser()
    // const doc = parser.parseFromString(astr, "text/html")
    // console.log("doc >> ", doc)

    // console.log("a1 >> ", doc.querySelector("a").getAttribute("href"))

    // const a = doc.getElementsByTagName("a")
    // console.log("a >> ", a)
    // console.log("a >> ", a.getAttribute("href"))


    console.log("App mounted");
    widgetData.current = 0;
    widgetData.max = 0;
    buttons.showRefill = false;
    buttons.showPlay = false;


    let args = getQueryVariables(window.location.href)
    carouselId = window.location.href.split('/').at(-1).split('?').at(0);
    args.forEach(arg => {
        if (arg.name == "type") {
            popup.type = arg.value
            popup.show = true
        } else if (arg.name == "msg") {
            popup.msg = arg.value
            popup.show = true
        }
    });

    if (popup.show) {
        setTimeout(() => {
            popup.show = false
            fullRefresh()
        }, 2000);
    }


    console.log(args)
    fullRefresh()

});

async function fullRefresh() {
    await fetchPriceOptions(carouselId);
    await fetchCarousel(carouselId);
    await updateComponentsView(widgetData, buttons, buyOptData, popup)
}


function onClick() {
    if (widgetData.current > 0) {
        widgetData.current--;
    }
    updateComponentsView(widgetData, buttons, buyOptData, popup)


}

function onEventPlayClick() {
    console.log("play clicked ");
    playCarousel(carouselId).then(function (result) {
        if (result) {
            setTimeout(() => {
                fullRefresh()
                // fetchCarousel(carouselId)
            }, 500);
        }

    });
}

function onEventRefillClick() {
    console.log("refill clicked");
    refillCarousel(carouselId, 3).then(function (result) {
        if (result) {
            setTimeout(() => {
                fetchCarousel(carouselId)
            }, 1000);
        }

    });
}
function onEventBuyClick(button) {
    console.log("Buy clicked %d, %d, %s", button.getPrice(), button.getTickets(), button.getPriceId());
    console.log("refill clicked");

    checkout(carouselId, button.getPriceId())
}
// function created() {
//         this.$setTitle("Login Page - My Website");
//         this.$setMeta([
//             {
//                 name: "description",
//                 content: "This is the login page description"
//             },
//             {
//                 name: "keywords",
//                 content: "Account, Login, Other, Keywords"
//             }
//         ])
//     }

</script>


<script>
import BuyTicketButton from "./components/BuyTicketButton.vue";
import WidgetActiveTickets from "./components/WidgetActiveTickets.vue";
import GeneralButton from './components/GeneralButton.vue';
import PopUp from './components/PopUp.vue';
export default {
    name: "App",
    components: {
        PopUp,
        WidgetActiveTickets,
        BuyTicketButton,
        GeneralButton,
    },
};

</script>

<template>
    <!-- <teleport to="head"> -->
        <!-- <meta http-equiv="Content-Security-Policy" content="upgrade-insecure-requests"> -->
    <!-- </teleport> -->
    <PopUp v-bind="popup" />
    <WidgetActiveTickets v-bind="widgetData" @click="onClick" />
    <GeneralButton v-bind="buttons" @play-clicked="onEventPlayClick" @refill-clicked="onEventRefillClick" />
    <BuyTicketButton v-bind="buyOptData" @buy-clicked="onEventBuyClick" />
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
