<script setup>
import { ref, onMounted } from "vue"
// import axios from 'axios'

// const token = ref(null)
const stripe = ref(null)
const elements = ref(null)

onMounted(() => {
    // axios.post('INITIATE_PAYMENT_API', {
    //     amount: 150,
    //     currency: 'USD'
    // }).then(response => {
    //     token.value = response.token // Use to identify the payment
    //     // stripe.value = Stripe(STRIPE_PUBLISHABLE_KEY);
    //     const options = {
    //         clientSecret: response.clientSecret,
    //     }

    //     elements.value = stripe.value.elements(options);
    //     const paymentElement = elements.value.create('payment');
    //     paymentElement.mount('#payment-element');
    // }).catch(error => {
    //     console.log(error)
    //     // throw error
    // })
})

async function handleSubmit(e) {
    e.preventDefault();

    const { error } = await stripe.value.confirmPayment({
        elements: elements.value,
        redirect: "if_required"
    });

    if (error === undefined) {
        // axios.post("PAYMENT_SUCCESS_API", {
        //     token: token.value,
        // })
    } else {
        // axios.post("PAYMENT_FAILURE_API", {
        //     token: token.value,
        //     code: error.code,
        //     description: error.message,
        // })
    }
}
</script>

<template>
    <form id="payment-form">
        <div id="payment-element">
            <!-- Stripe will create form elements here -->
        </div>
        <button type="submit" @click="handleSubmit">Pay via Stripe</button>
    </form>
</template>

<style></style>