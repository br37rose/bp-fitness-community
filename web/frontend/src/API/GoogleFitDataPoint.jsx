import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from 'humps';
import { DateTime } from "luxon";

import {
    BP8_FITNESS_GOOGLE_FIT_DATA_POINTS_API_ENDPOINT
} from "../Constants/API";


export function getGoogleFitDataPointListAPI(filtersMap = new Map(), onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    // The following code will generate the query parameters for the url based on the map.
    let aURL = BP8_FITNESS_GOOGLE_FIT_DATA_POINTS_API_ENDPOINT;
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
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // Bugfixes.
        console.log("getDataPointListAPI | pre-fix | results:", data);
        if (data.results !== undefined && data.results !== null && data.results.length > 0) {
            data.results.forEach(
                (item, index) => {
                    item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
                    item.startAt = DateTime.fromISO(item.startAt).toLocaleString(DateTime.DATETIME_MED);
                    item.endAt = DateTime.fromISO(item.endAt).toLocaleString(DateTime.DATETIME_MED);
                    // console.log(item, index);
                }
            )
        }
        console.log("getDataPointListAPI | post-fix | results:", data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch((exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}
