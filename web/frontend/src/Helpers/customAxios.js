import axios from 'axios';

import { getAPIBaseURL } from './urlUtility';
import {
    getAccessTokenFromLocalStorage,
    getRefreshTokenFromLocalStorage,
    setAccessTokenInLocalStorage,
    setRefreshTokenInLocalStorage
} from './jwtUtility';
import { BP8_FITNESS_REFRESH_TOKEN_API_ENDPOINT } from "../Constants/API";

/**
 *  Function returns a custom `Axios` instance tailered to the `cps backend`
 *  API web-service for authenticated users.
 *
 *  Features:
 *  (1) Inform API to expect request encoded with `JSON` format.
 *  (2) Inform API we expect responses to be in `JSON` format.
 *  (3) Attach authorization bearer token.
 *  (4) Integrate automatic refresh token when token expires.
 */
export default function getCustomAxios() {
    // IMPORTANT: THIS IS THE ONLY WAY WE CAN GET THE ACCESS TOKEN.
    const accessToken = getAccessTokenFromLocalStorage();

    // Create a new Axios instance using our oAuth 2.0 bearer token
    // and various other headers.
    const customAxios = axios.create({
        baseURL: getAPIBaseURL(),
        headers: {
            'Authorization': "JWT " + accessToken,
            'Content-Type': 'application/json;',
            'Accept': 'application/json',
        },
    });

    // Attach our Axios "refesh token" interceptor.
    customAxios.interceptors.response.use( // Special thanks to: https://stackoverflow.com/a/69309837
        response => {
            return response;
        },
        async error => {

            // unauthen error 401
            let originalConfig = error.config;
            if (error.response.status === 401) {
                console.log("Axios Interceptor: Detected 401. Handling now...");

                // get token from storage
                const token = getRefreshTokenFromLocalStorage();

                console.log("Axios Interceptor: Detected 401. token:", token);
                if (token) {
                    // Make an API call to the remote service to refresh token
                    // and wait until the service returns a result and then
                    // handle the response here...
                    const respRefresh = await handleRefresh(token);

                    // On successful token refreshal, run the following code.
                    if (respRefresh && respRefresh.status === 200) {

                        // Extract the new values from the response.
                        const newAccessToken = respRefresh.data.access_token;
                        const newRefreshToken = respRefresh.data.refresh_token;

                        // Save.
                        setAccessTokenInLocalStorage(newAccessToken);
                        setRefreshTokenInLocalStorage(newRefreshToken);

                        // Reset our axios authorization header to use our
                        // new token but keep the original configuration intact.
                        originalConfig = {
                            ...originalConfig,
                            headers: {
                                ...originalConfig.headers,
                                Authorization: `JWT ${newAccessToken}`,
                            },
                        };

                        // continue send currently request
                        return customAxios(originalConfig);
                    }
                }
            }

            // Catch-all return case, this will result in the error going through.
            return Promise.reject(error.response.data);
        },
    );

    // Return our custom Axios instance for our application.
    return customAxios;
}

const axiosServiceRefresh = axios.create({
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
    },
});
const handleRefresh = async token => {
    // Set the data we will be sending to the 'refresh token' API endpoint
    // in the backend.
    const param = {
        value: token,
    };

    // set header for axios refresh
    axiosServiceRefresh.defaults.headers.common.Authorization = `Bearer ${token}`;
    return new Promise((resolve, reject) => {
        axiosServiceRefresh
        .post(BP8_FITNESS_REFRESH_TOKEN_API_ENDPOINT, param)
        .then(response => {
            resolve(response);
        })
        .catch(error => {});
    });
};

// Auto-generated comment for change 12
