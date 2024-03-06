import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from 'humps';
import { DateTime } from "luxon";

import {
    BP8_FITNESS_RANK_POINTS_API_ENDPOINT,
    BP8_FITNESS_LEADERBOARD_API_ENDPOINT,
} from "../Constants/API";


export function getRankPointListAPI(filtersMap = new Map(), onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    // The following code will generate the query parameters for the url based on the map.
    let aURL = BP8_FITNESS_RANK_POINTS_API_ENDPOINT;
    filtersMap.forEach(
        (value, key) => {
            let decamelizedkey = decamelize(key)
            if (aURL.indexOf('?') > -1) {
                aURL += "&" + decamelizedkey + "=" + value;
            } else {
                aURL += "?" + decamelizedkey + "=" + value;
            }
        }
    )

    axios.get(aURL).then((successResponse) => {
        const responseRank = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseRank);

        // // Bugfixes.
        // console.log("getRankPointListAPI | pre-fix | results:", data);
        // if (data.results !== undefined && data.results !== null && data.results.length > 0) {
        //     data.results.forEach(
        //         (item, index) => {
        //             item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
        //             console.log(item, index);
        //         }
        //     )
        // }
        // console.log("getRankPointListAPI | post-fix | results:", data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch((exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

// export function getLeaderboardListAPI(filtersMap=new Map(), onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();

//     // The following code will generate the query parameters for the url based on the map.
//     let aURL = BP8_FITNESS_LEADERBOARD_API_ENDPOINT;
//     filtersMap.forEach(
//         (value, key) => {
//             let decamelizedkey = decamelize(key)
//             if (aURL.indexOf('?') > -1) {
//                 aURL += "&"+decamelizedkey+"="+value;
//             } else {
//                 aURL += "?"+decamelizedkey+"="+value;
//             }
//         }
//     )

//     axios.get(aURL).then((successResponse) => {
//         const responseRank = successResponse.data;

//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseRank);

//         // // Bugfixes.
//         // console.log("getRankPointListAPI | pre-fix | results:", data);
//         // if (data.results !== undefined && data.results !== null && data.results.length > 0) {
//         //     data.results.forEach(
//         //         (item, index) => {
//         //             item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
//         //             console.log(item, index);
//         //         }
//         //     )
//         // }
//         // console.log("getRankPointListAPI | post-fix | results:", data);

//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }

// API.js

export function getLeaderboardListAPI(filtersMap = new Map()) {
    const axios = getCustomAxios();
    let aURL = BP8_FITNESS_LEADERBOARD_API_ENDPOINT;

    // Generate the query parameters for the URL based on the map.
    filtersMap.forEach((value, key) => {
        let decamelizedKey = decamelize(key);
        aURL += aURL.includes('?') ? `&${decamelizedKey}=${value}` : `?${decamelizedKey}=${value}`;
    });

    // Return the Axios promise chain
    return axios.get(aURL).then((successResponse) => {
        const responseRank = successResponse.data;
        return camelizeKeys(responseRank);
    }).catch((error) => {
        // In case of error, rethrow the error to be caught by the calling function
        throw camelizeKeys(error);
    });
}
