import { camelizeKeys, decamelize } from "humps";
import {
  BP8_FITNESS_T_PROGRAM_API_ENDPOINT,
  BP8_FITNESS_T_PROGRAM_PHASE_PATCH_API_ENDPOINT,
} from "../Constants/API";
import getCustomAxios from "../Helpers/customAxios";
import { DateTime } from "luxon";

export function postTrainingProgCreateAPI(
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();

  axios
    .post(BP8_FITNESS_T_PROGRAM_API_ENDPOINT, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getTrainingProgListApi(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();

  // The following code will generate the query parameters for the url based on the map.
  let aURL = BP8_FITNESS_T_PROGRAM_API_ENDPOINT;
  filtersMap.forEach((value, key) => {
    let decamelizedkey = decamelize(key);
    if (aURL.indexOf("?") > -1) {
      aURL += "&" + decamelizedkey + "=" + value;
    } else {
      aURL += "?" + decamelizedkey + "=" + value;
    }
  });

  axios
    .get(aURL)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      if (
        data.results !== undefined &&
        data.results !== null &&
        data.results.length > 0
      ) {
        data.results.forEach((item, index) => {
          item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(
            DateTime.DATETIME_MED
          );
        });
      }

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getTrainingProgDetailAPI(
  TrainingProgId,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();
  axios
    .get(BP8_FITNESS_T_PROGRAM_API_ENDPOINT + "/" + TrainingProgId)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // For debugging purposeso pnly.

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}
export function putTrainingProgUpdateAPI(
  TrainingProgId,
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();

  axios
    .put(
      BP8_FITNESS_T_PROGRAM_API_ENDPOINT + "/" + TrainingProgId,
      decamelizedData
    )
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function deleteTrainingProgAPI(
  TrainingProgId,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();
  axios
    .delete(BP8_FITNESS_T_PROGRAM_API_ENDPOINT + "/" + TrainingProgId)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function patchTrainingProgAPI(
  prgmid,
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback
) {
  const axios = getCustomAxios();
  let url = BP8_FITNESS_T_PROGRAM_PHASE_PATCH_API_ENDPOINT.replace(
    "{pid}",
    prgmid
  );
  axios
    .patch(url, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}
