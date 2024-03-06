import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from 'humps';
import { DateTime } from "luxon";

import {
    BP8_FITNESS_CREATE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT,
    BP8_FITNESS_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT,
    BP8_FITNESS_CANCEL_SUBSCRIPTION_API_ENDPOINT,
    BP8_FITNESS_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT,
    // BP8_FITNESS_PAYMENT_PROCESSOR_SEND_SUBSCRIPTION_REQUEST_EMAIL_API_ENDPOINT,
    // BP8_FITNESS_PAYMENT_PROCESSOR_GRANT_FREE_CREDITS_API_ENDPOINT
} from "../Constants/API";

export function postCreateStripeSubscriptionCheckoutSessionAPI(priceID, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    const postData = {
        "price_id": priceID,
    };

    axios.post(BP8_FITNESS_CREATE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT, postData).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function getCompleteStripeSubscriptionCheckoutSessionAPI(sessionID, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    axios.get(BP8_FITNESS_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT.replace("{sessionID}", sessionID)).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // For debugging purposeso pnly.
        console.log("completeStripeSubscriptionCheckoutSession: Response Data: ", data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postSubscriptionCancelAPI(memberID="", onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    const aURL = BP8_FITNESS_CANCEL_SUBSCRIPTION_API_ENDPOINT + "?member_id="+memberID;

    axios.post(aURL).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function getPaymentProcessorStripeInvoiceListAPI(userID, cursor, pageSize, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    let aURL = BP8_FITNESS_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT.replace("{userID}",userID).replace("{cursor}",cursor).replace("{pageSize}",pageSize);

    axios.get(aURL).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // Bugfixes.
        // console.log("getMemberListAPI | pre-fix | results:", data);
        // if (data.results !== undefined && data.results !== null && data.results.length > 0) {
        //     data.results.forEach(
        //         (item, index) => {
        //             item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
        //             console.log(item, index);
        //         }
        //     )
        // }
        // console.log("getMemberListAPI | post-fix | results:", data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

// export function postPaymentProcessorSendSubscriptionRequestEmailAPI(userID, offerId, onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();
//     const data = {
//         member_id: userID,
//         offer_id: offerId,
//     };
//     axios.post(BP8_FITNESS_PAYMENT_PROCESSOR_SEND_SUBSCRIPTION_REQUEST_EMAIL_API_ENDPOINT, data).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
//
// export function postPaymentProcessorGrantFreeCreditAPI(decamelizedData, onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();
//     axios.post(BP8_FITNESS_PAYMENT_PROCESSOR_GRANT_FREE_CREDITS_API_ENDPOINT, decamelizedData).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
