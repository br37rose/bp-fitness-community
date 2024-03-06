import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from 'humps';
import { DateTime } from "luxon";

import {
    BP8_FITNESS_REGISTRY_API_ENDPOINT,
} from "../Constants/API";


export function getRegistryAPI(cpsrn, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    axios.get(BP8_FITNESS_REGISTRY_API_ENDPOINT.replace("{id}", cpsrn)).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = camelizeKeys(responseData);

        // Minor bugfix.
        data.showsSignsOfTamperingOrRestoration = parseInt(data.showsSignsOfTamperingOrRestoration);
        data.issueCoverDate = DateTime.fromISO(data.issueCoverDate).toJSDate();

        // For debugging purposeso pnly.
        console.log(data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}
