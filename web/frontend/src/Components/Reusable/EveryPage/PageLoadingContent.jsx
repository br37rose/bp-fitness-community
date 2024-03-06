import React from "react";

function PageLoadingContent(props) {
    const { displayMessage } = props;
    return (
        <div class="columns is-centered" style={{ paddingTop: "20px" }}>
            <div class="column has-text-centered is-1">
                <div class="loader-wrapper is-centered">
                    <br />
                    <br />
                    <div class="loader is-loading is-centered" style={{ height: "80px", width: "80px" }}></div>
                </div>
                <br />
                <div className="">{displayMessage}</div>
                <br />
                <br />
            </div>
        </div>
    );
}

export default PageLoadingContent;
