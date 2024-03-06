import React from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCloudDownload} from '@fortawesome/free-solid-svg-icons';
import getCustomAxios from "../../Helpers/customAxios";
// import axios from 'axios';


function DataDisplayRowImage(props) {
    const { label, src, alt } = props;

    console.log("DataDisplayRowImage | src:", src); // For debugging purposes only.

    const onDownloadClick = () => {
        // Get the axios handler which is has access to authenticated token.
        const axios = getCustomAxios(null);

        //
        axios.get(src, {
           responseType: 'blob', // Specify the response type as a blob
        }).then((successResponse) => {
            const responseData = successResponse.data;

            // Create a link element
            const link = document.createElement('a');

            // Create a URL for the blob and set it as the href for the link
            const blobURL = URL.createObjectURL(responseData);
            link.href = blobURL;

            // Set the filename for the download
           const filename = src.split('/').pop(); // Extract filename from the URL
           link.download = filename;

           // Append the link to the body
           document.body.appendChild(link);

           // Trigger a click on the link to start the download
           link.click();

           // Remove the link from the body
           document.body.removeChild(link);

           // Release the URL object
           URL.revokeObjectURL(blobURL);

        }).catch( (exception) => {
            console.log(exception)
        }).then(

        );
    }

    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <figure class="image">
                    <img src={src} alt={alt} />
                </figure>
            </div>
            <button class="button is-small" type="button" onClick={onDownloadClick}><FontAwesomeIcon className="fas" icon={faCloudDownload}/>&nbsp;Download Image</button>
        </div>
    );
}

export default DataDisplayRowImage;
