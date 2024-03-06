import { BP8_FITNESS_API_BASE_PATH } from "../Constants/API";


/**
 *  Function takes the React environment variables and returns the base URL used
 *  for all API communication with our web-service. Please use this function to set the
 *  ``Axios`` base URL when making API calls to the backend server.
 */
export function getAPIBaseURL() {
    return process.env.REACT_APP_API_PROTOCOL + "://" + process.env.REACT_APP_API_DOMAIN + BP8_FITNESS_API_BASE_PATH;
}

export function getAppBaseURL() {
    return process.env.REACT_APP_WWW_PROTOCOL + "://" + process.env.REACT_APP_WWW_DOMAIN;
}

/**
 * Get the URL parameters
 * source: https://css-tricks.com/snippets/javascript/get-url-variables/
 * @param  {String} url The URL
 * @return {Object}     The URL parameters
 */
export function getParams(url) {
	var params = {};
	var parser = document.createElement('a');
	parser.href = url;
	var query = parser.search.substring(1);
	var vars = query.split('&');
	for (var i = 0; i < vars.length; i++) {
		var pair = vars[i].split('=');
		params[pair[0]] = decodeURIComponent(pair[1]);
	}
	return params;
};
