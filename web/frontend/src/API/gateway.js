import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys } from 'humps';
import {
    BP8_FITNESS_LOGIN_API_ENDPOINT,
    BP8_FITNESS_VERSION_ENDPOINT,
    BP8_FITNESS_REGISTER_MEMBER_API_ENDPOINT,
    BP8_FITNESS_EMAIL_VERIFICATION_API_ENDPOINT,
    BP8_FITNESS_LOGOUT_API_ENDPOINT,
    BP8_FITNESS_FORGOT_PASSWORD_API_ENDPOINT,
    BP8_FITNESS_PASSWORD_RESET_API_ENDPOINT,
    BP8_FITNESS_2FA_GENERATE_OTP_API_ENDPOINT,
    BP8_FITNESS_2FA_GENERATE_OTP_AND_QR_CODE_API_ENDPOINT,
    BP8_FITNESS_2FA_VERIFY_OTP_API_ENDPOINT,
    BP8_FITNESS_2FA_VALIDATE_OTP_API_ENDPOINT,
    BP8_FITNESS_2FA_DISABLED_OTP_API_ENDPOINT
} from "../Constants/API";

import {
    setAccessTokenInLocalStorage,
    setRefreshTokenInLocalStorage
} from '../Helpers/jwtUtility';

export function postLoginAPI(data, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    // To Snake-case for API from camel-case in React.
    let decamelizedData = decamelizeKeys(data);

    axios.post(BP8_FITNESS_LOGIN_API_ENDPOINT, decamelizedData).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        let profile = camelizeKeys(responseData);

        console.log("postLoginAPI | prefix:", profile);

        // BUGFIX: Have both styles.
        profile.user.organizationID = profile.user.organizationId;

        console.log("postLoginAPI | postfix:", profile);

        // SAVE OUR CREDENTIALS IN PERSISTENT STORAGE. THIS IS AN IMPORTANT
        // STEP BECAUSE OUR TOKEN UTILITY HELPER NEEDS THIS.
        setAccessTokenInLocalStorage(profile.accessToken);
        setRefreshTokenInLocalStorage(profile.refreshToken);

        // Return the callback data.
        onSuccessCallback(profile);
    }).catch( (exception) => {
        let responseData = null;
        if (exception.response !== undefined && exception.response !== null) {
            if (exception.response.data !== undefined && exception.response.data !== null) {
                responseData = exception.response.data;
            } else {
                responseData = exception.response;
            }
        } else {
            responseData = exception;
        }
        let errors = camelizeKeys(responseData);

        // Check for incorrect password and enter our own custom error.
        let errorsStr = JSON.stringify(errors);
        if (errorsStr.includes("Incorrect email or password")) { // NOTE: This is the exact error from backend on incorrect email/pass.
            errors = {
                "auth": "Incorrect email or password",
            };
        }

        onErrorCallback(errors);
    }).then(onDoneCallback);
}


export function postRegisterAPI(decamelizedData, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    axios.post(BP8_FITNESS_REGISTER_MEMBER_API_ENDPOINT, decamelizedData).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        let profile = camelizeKeys(responseData);

        // SAVE OUR CREDENTIALS IN PERSISTENT STORAGE. THIS IS AN IMPORTANT
        // STEP BECAUSE OUR TOKEN UTILITY HELPER NEEDS THIS.
        if (profile.accessToken) {
            setAccessTokenInLocalStorage(profile.accessToken);
        }
        if (profile.refreshToken) {
            setRefreshTokenInLocalStorage(profile.refreshToken);
        }

        // Return the callback data.
        onSuccessCallback(profile);
    }).catch( (exception) => {
        let responseData = null;
        if (exception.response !== undefined && exception.response !== null) {
            if (exception.response.data !== undefined && exception.response.data !== null) {
                responseData = exception.response.data;
            } else {
                responseData = exception.response;
            }
        } else {
            responseData = exception;
        }
        let errors = camelizeKeys(responseData);

        // Check for incorrect password and enter our own custom error.
        let errorsStr = JSON.stringify(errors);
        if (errorsStr.includes("Incorrect email or password")) { // NOTE: This is the exact error from backend on incorrect email/pass.
            errors = {
                "auth": "Incorrect email or password",
            };
        }

        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function getVersionAPI(onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    axios.get(BP8_FITNESS_VERSION_ENDPOINT).then((successResponse) => {
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

export function postEmailVerificationAPI(verificationCode, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    let data = {
        code: verificationCode,
    };
    axios.post(BP8_FITNESS_EMAIL_VERIFICATION_API_ENDPOINT, data).then((successResponse) => {
        onSuccessCallback(null);
    }).catch( (exception) => {
        let responseData = null;
        if (exception.response !== undefined && exception.response !== null) {
            if (exception.response.data !== undefined && exception.response.data !== null) {
                responseData = exception.response.data;
            } else {
                responseData = exception.response;
            }
        } else {
            responseData = exception;
        }
        let errors = camelizeKeys(responseData);

        // Check for incorrect password and enter our own custom error.
        let errorsStr = JSON.stringify(errors);
        if (errorsStr.includes("Incorrect email or password")) { // NOTE: This is the exact error from backend on incorrect email/pass.
            errors = {
                "auth": "Incorrect email or password",
            };
        }

        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postLogoutAPI(onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();
    let data = {};
    axios.post(BP8_FITNESS_LOGOUT_API_ENDPOINT, data).then((successResponse) => {
        onSuccessCallback(null);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postForgotPasswordAPI(email, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    axios.post(BP8_FITNESS_FORGOT_PASSWORD_API_ENDPOINT, {email: email}).then((successResponse) => {
        console.log("postForgotPasswordAPI: Success")
        onSuccessCallback(); // Return the callback data.
    }).catch( (exception) => {
        let responseData = null;
        if (exception.response !== undefined && exception.response !== null) {
            if (exception.response.data !== undefined && exception.response.data !== null) {
                responseData = exception.response.data;
            } else {
                responseData = exception.response;
            }
        } else {
            responseData = exception;
        }
        let errors = camelizeKeys(responseData);

        // Check for incorrect password and enter our own custom error.
        let errorsStr = JSON.stringify(errors);
        if (errorsStr.includes("Incorrect email or password")) { // NOTE: This is the exact error from backend on incorrect email/pass.
            errors = {
                "auth": "Incorrect email or password",
            };
        }

        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postPasswordResetAPI(verificationCode, password, passwordRepeat, onSuccessCallback, onErrorCallback, onDoneCallback) {
    const axios = getCustomAxios();

    axios.post(BP8_FITNESS_PASSWORD_RESET_API_ENDPOINT, {verification_code: verificationCode, password: password, password_repeated: passwordRepeat}).then((successResponse) => {
        console.log("postForgotPasswordAPI: Success")
        onSuccessCallback(); // Return the callback data.
    }).catch( (exception) => {
        let responseData = null;
        if (exception.response !== undefined && exception.response !== null) {
            if (exception.response.data !== undefined && exception.response.data !== null) {
                responseData = exception.response.data;
            } else {
                responseData = exception.response;
            }
        } else {
            responseData = exception;
        }
        let errors = camelizeKeys(responseData);

        // Check for incorrect password and enter our own custom error.
        let errorsStr = JSON.stringify(errors);
        if (errorsStr.includes("Incorrect email or password")) { // NOTE: This is the exact error from backend on incorrect email/pass.
            errors = {
                "auth": "Incorrect email or password",
            };
        }

        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postGenerateOTP(onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
    const axios = getCustomAxios(onUnauthorizedCallback);

    let aURL = BP8_FITNESS_2FA_GENERATE_OTP_API_ENDPOINT;

    axios.post(aURL).then((successResponse) => {
        const responseData = successResponse.data;

        // Snake-case from API to camel-case for React.
        const data = {
            base32: responseData.base32,
            optAuthURL: responseData.otpauth_url,
        };

        // console.log("getTagListAPI | post-fix | results:", data);

        // Return the callback data.
        onSuccessCallback(data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postGenerateOTPAndQRCodeImage(onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
    const axios = getCustomAxios(onUnauthorizedCallback);

    let aURL = BP8_FITNESS_2FA_GENERATE_OTP_AND_QR_CODE_API_ENDPOINT;

    axios.post(aURL, { responseType: 'blob',}).then((successResponse) => {
        const binaryData = successResponse.data;

        // Create a Blob from the binary data
        const blob = new Blob([binaryData], { type: 'image/png' });

        // Create a Blob URL from the Blob object
        const blobUrl = URL.createObjectURL(blob);

        console.log("blobUrl", blobUrl);

        // Call the success callback with the Blob URL
        onSuccessCallback(blobUrl);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postVertifyOTP(payload, onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
    const axios = getCustomAxios(onUnauthorizedCallback);
    let aURL = BP8_FITNESS_2FA_VERIFY_OTP_API_ENDPOINT;
    axios.post(aURL, payload).then((successResponse) => {
        // Return the callback data.
        onSuccessCallback(successResponse.data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postValidateOTP(payload, onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
    const axios = getCustomAxios(onUnauthorizedCallback);
    let aURL = BP8_FITNESS_2FA_VALIDATE_OTP_API_ENDPOINT;
    axios.post(aURL, payload).then((successResponse) => {
        // Return the callback data.
        onSuccessCallback(successResponse.data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}

export function postDisableOTP(onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
    const axios = getCustomAxios(onUnauthorizedCallback);

    let aURL = BP8_FITNESS_2FA_DISABLED_OTP_API_ENDPOINT;

    axios.post(aURL).then((successResponse) => {
        onSuccessCallback(successResponse.data);
    }).catch( (exception) => {
        let errors = camelizeKeys(exception);
        onErrorCallback(errors);
    }).then(onDoneCallback);
}
