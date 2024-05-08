import React from "react";

function PrivacyStaticPage() {
  // For debugging purposes only.
  console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
  console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
  console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
  console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div className="container">
        <section className="section">
          <h1>TODO: Privacy page</h1>
        </section>
      </div>
    </>
  );
}

export default PrivacyStaticPage;
