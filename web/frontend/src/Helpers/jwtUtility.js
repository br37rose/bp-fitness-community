
/**
 *------------------------------------------------------------------------------
 * The purpose of this utility is to handle all API token related functionality
 * and provide an interface for the system to use.
 *------------------------------------------------------------------------------
 */


/**
 *  Saves our access token to persistent storage.
 */
export function setAccessTokenInLocalStorage(accessToken) {
    if (accessToken !== undefined && accessToken !== null) {
        localStorage.setItem("BP8_FITNESS_TOKEN_UTILITY_ACCESS_TOKEN_DATA", accessToken);
    } else {
        console.error("Setting undefined access token");
    }
}


/**
 *  Saves our refresh token to our persistent storage.
 */
export function setRefreshTokenInLocalStorage(accessToken) {
    if (accessToken !== undefined && accessToken !== null)  {
        localStorage.setItem("BP8_FITNESS_TOKEN_UTILITY_REFRESH_TOKEN_DATA", accessToken);
    } else {
        console.error("Setting undefined resfresh token");
    }
}

/**
 *  Gets our access token from persistent storage.
 */
export function getAccessTokenFromLocalStorage() {
    return localStorage.getItem("BP8_FITNESS_TOKEN_UTILITY_ACCESS_TOKEN_DATA");
}


/*
 *  Gets our refresh token from persistent storage.
 */
export function getRefreshTokenFromLocalStorage() {
    return localStorage.getItem("BP8_FITNESS_TOKEN_UTILITY_REFRESH_TOKEN_DATA");
}


/*
 *  Clears all the tokens on the user's browsers persistent storage.
 */
export function clearAllAccessAndRefreshTokensFromLocalStorage() {
    localStorage.removeItem("BP8_FITNESS_TOKEN_UTILITY_ACCESS_TOKEN_DATA");
    localStorage.removeItem("BP8_FITNESS_TOKEN_UTILITY_REFRESH_TOKEN_DATA");
}
